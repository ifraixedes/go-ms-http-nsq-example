//Package redis is a implementation of the Driver Location Service using Redis
// database
package redis

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"go.fraixed.es/errors"
)

// NewService return a Service instance implementing the buisniess logic
// using the Redis database specified in options.
//
// NOTE: that any method of the returned service supports the context
// cancellations because the used Redis library doesn't contemplate them.
//
// It returns an error if the Redis connection cannot be established (error code
// ErrInvalidRedisConfig)
func NewService(o Options) (*Service, error) {
	var (
		ro = redis.Options(o)
		c  = redis.NewClient(&ro)
		sc = c.Ping()
	)

	if err := sc.Err(); err != nil {
		return nil, errors.New(ErrInvalidRedisConfig, errors.MD{K: "config", V: o})
	}

	return &Service{
		cli:    c,
		closed: make(chan struct{}),
	}, nil
}

// Service implements the drvloc.Service interface using Redis as Database.
// An instance of this type must always be initialized with the constructor
// functions and using the zero value will panic.
//
// All the exported methods of an instance of this type can return the
// ErrClosedService after the Close method is called. Once Close is called, the
// instance isn't usable anymore.
type Service struct {
	cli redis.Cmdable

	// the following fields are used by the isClose and Close methods, the rest
	// of function shouldn't never use them
	closed  chan struct{}
	closing sync.Once
}

// SetLocation stores the driver location satisfying the drvloc.Service
// interface.
//
// The following error codes can be returned:
//
// * ErrUnexpectedErrorStore
func (s *Service) SetLocation(ctx context.Context, id uint64, l drvloc.Location) error {
	var cli, err = s.getClient(ctx)
	if err != nil {
		return err
	}

	var (
		key = getDriverLocationSortedSetKey(id)
		z   = redis.Z{
			Score:  float64(l.At.Unix()),
			Member: &l,
		}
	)

	var ic = cli.ZAddNX(key, z)
	if err := s.handleUnexpectedErr(ic); err != nil {
		return err
	}

	return nil
}

// LocationsForLastMinutes retrieve the driver location satisfying the drvloc.Service
// interface.
//
// The following error codes can be returned:
//
// * ErrInvalidDataFormat
//
// * ErrUnexpectedErrorStore
func (s *Service) LocationsForLastMinutes(ctx context.Context, id uint64, minutes uint16) ([]drvloc.Location, error) {
	var cli, err = s.getClient(ctx)
	if err != nil {
		return nil, err
	}

	var (
		key = getDriverLocationSortedSetKey(id)
		ic  = cli.Exists(key)
	)

	if err := s.handleUnexpectedErr(ic); err != nil {
		return nil, err
	}

	if n, err := ic.Result(); err != nil {
		return nil, errors.Wrap(err, drvloc.ErrUnexpectedErrorStore)
	} else if n == 0 {
		return nil, errors.New(drvloc.ErrNotFoundDriver, errors.MD{K: "id", V: id})
	}

	var ssc = cli.ZRangeByScore(key, redis.ZRangeBy{
		Min: strconv.Itoa(int(time.Now().Unix() - (int64(minutes) * 60))),
		Max: "+inf",
	})
	if err := s.handleUnexpectedErr(ssc); err != nil {
		return nil, err
	}

	var locs []drvloc.Location
	if err := ssc.ScanSlice(&locs); err != nil {
		return nil, errors.Wrap(err, drvloc.ErrInvalidDataFormatStore)
	}

	return locs, nil
}

// Close closes release the resources of the service and close the connections
// with the external services (i.e. Redis)
//
// The following error codes can be returned:
//
// * ErrUnexpectedErrorStore
func (s *Service) Close(ctx context.Context) error {
	if s.isClosed() {
		return errors.New(ErrClosedService)
	}

	var err = errors.New(ErrClosedService)
	s.closing.Do(func() {
		close(s.closed)
		if c, ok := s.cli.(closer); ok {
			if err := c.Close(); err != nil {
				err = errors.New(drvloc.ErrUnexpectedErrorStore)
			}
		}

		err = nil
	})

	return err
}

func (s *Service) isClosed() bool {
	select {
	case _, ok := <-s.closed:
		return ok == false
	default:
		return false
	}
}

func (s *Service) getClient(ctx context.Context) (redis.Cmdable, error) {
	if s.isClosed() {
		return nil, errors.New(ErrClosedService)
	}

	switch c := s.cli.(type) {
	case *redis.Client:
		return c.WithContext(ctx), nil
	case *redis.ClusterClient:
		return c.WithContext(ctx), nil
	case *redis.Ring:
		return c.WithContext(ctx), nil
	}

	panic("unsupported redis.Cmdable")
}

func (s *Service) handleUnexpectedErr(ce errorer) error {
	if err := ce.Err(); err != nil {
		if s.isClosed() {
			return errors.Wrap(err, ErrClosedService)
		}
		return errors.Wrap(err, drvloc.ErrUnexpectedErrorStore)
	}

	return nil
}

func getDriverLocationSortedSetKey(id uint64) string {
	return fmt.Sprintf("drvloc_sorted_set_%d", id)
}
