package redis_test

import (
	"testing"

	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/internal/testassert"
	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/redis"
	"github.com/stretchr/testify/assert"
	"go.fraixed.es/errors"
)

func TestNewService(t *testing.T) {
	type tcase struct {
		desc   string
		o      redis.Options
		assert func(*testing.T, tcase, drvloc.Service, error)
	}

	var tcases = []tcase{
		{
			desc: "ok",
			o: redis.Options{
				Addr: testRedisAddr,
			},
			assert: func(t *testing.T, _ tcase, s drvloc.Service, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, s)
			},
		},
		{
			desc: "error: invalid connection",
			o: redis.Options{
				Addr: "google.com:6379",
			},
			assert: func(t *testing.T, tc tcase, _ drvloc.Service, err error) {
				testassert.ErrorWithCode(t, err, redis.ErrInvalidRedisConfig, errors.MD{K: "config", V: tc.o})
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			var s, err = redis.NewService(tc.o)
			tc.assert(t, tc, s, err)
		})
	}
}

func TestService_SetLocation(t *testing.T) {
	t.Skipf("WIP: Implement")
}

func TestService_Locations(t *testing.T) {
	t.Skipf("WIP: Implement")
}
