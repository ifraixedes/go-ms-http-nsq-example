//Package redis is a implementation of the Driver Location Service using Redis
// database
package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"go.fraixed.es/errors"
)

// NewService satisifies the Service interface implementing the buisniess logic
// using the Redis database specified in options.
//
// It returns an error if the Redis connection cannot be established.
//
// NOTE: that any method of the returned service supports the context
// cancellations because the used Redis library doesn't contemplate them.
func NewService(o Options) (drvloc.Service, error) {
	var (
		ro = redis.Options(o)
		c  = redis.NewClient(&ro)
		sc = c.Ping()
	)

	if err := sc.Err(); err != nil {
		return nil, errors.New(ErrInvalidRedisConfig, errors.MD{K: "config", V: o})
	}

	return service{
		cli: c,
	}, nil
}

type service struct {
	cli redis.Cmdable
}

func (s service) SetLocation(ctx context.Context, id uint64, l drvloc.Location) error {
	var (
		key = getDriverLocationSortedSetKey(id)
		z   = redis.Z{
			Score:  float64(l.At.Unix()),
			Member: &l,
		}
		cli = s.getClient(ctx)
	)

	var ic = cli.ZAddNX(key, z)
	if err := ic.Err(); err != nil {
		return errors.Wrap(err, drvloc.ErrUnexpectedErrorStore)
	}

	return nil
}

func (s service) LocationsForLastMinutes(ctx context.Context, id uint64, minutes uint16) ([]drvloc.Location, error) {
	// TODO: WIP
	return nil, nil
}

func (s service) getClient(ctx context.Context) redis.Cmdable {
	switch c := s.cli.(type) {
	case *redis.Client:
		return c.WithContext(ctx)
	case *redis.ClusterClient:
		return c.WithContext(ctx)
	case *redis.Ring:
		return c.WithContext(ctx)
	}

	panic("unsupported redis.Cmdable")
}

func getDriverLocationSortedSetKey(id uint64) string {
	return fmt.Sprintf("drvloc_sorted_set_%d", id)
}
