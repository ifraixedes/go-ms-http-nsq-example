//Package redis is a implementation of the Driver Location Service using Redis
// database
package redis

import (
	"context"

	"github.com/go-redis/redis"
	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"go.fraixed.es/errors"
)

// NewService satisifies the Service interface implementing the buisniess logic
// using the Redis database specified in options.
// It returns an error if the Redis connection cannot be established.
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
	// TODO: WIP
	return nil
}

func (s service) LocationsForLastMinutes(ctx context.Context, id uint64, minutes uint16) ([]drvloc.Location, error) {
	// TODO: WIP
	return nil, nil
}
