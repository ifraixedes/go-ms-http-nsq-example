package redis_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/internal/testassert"
	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/redis"
	"github.com/mmcloughlin/spherand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	type params struct {
		ctx context.Context
		id  uint64
		l   drvloc.Location
	}

	type tcase struct {
		desc   string
		args   params
		assert func(*testing.T, tcase, error)
	}

	// Values to use for the different test cases
	var dID = uint64(rand.Int63n(99) + 1)
	var loc1 drvloc.Location
	{
		lat, lng := spherand.Geographical()
		loc1 = drvloc.Location{
			Lat: lat,
			Lng: lng,
			At:  time.Unix(time.Now().Unix()-rand.Int63n(5)-5, 0).Round(0),
		}
	}
	var loc2 drvloc.Location
	{
		lat, lng := spherand.Geographical()
		loc2 = drvloc.Location{
			Lat: lat,
			Lng: lng,
			At:  time.Now().Round(0),
		}
	}

	var tcases = []tcase{
		{
			desc: "ok: new driver and location",
			args: params{
				ctx: context.Background(),
				id:  dID,
				l:   loc1,
			},
			assert: func(t *testing.T, _ tcase, err error) {
				assert.NoError(t, err)
			},
		},
		{
			desc: "ok: existing driver and new location",
			args: params{
				ctx: context.Background(),
				id:  dID,
				l:   loc2,
			},
			assert: func(t *testing.T, _ tcase, err error) {
				assert.NoError(t, err)
			},
		},
		{
			desc: "ok: existing driver and same location",
			args: params{
				ctx: context.Background(),
				id:  dID,
				l:   loc2,
			},
			assert: func(t *testing.T, _ tcase, err error) {
				assert.NoError(t, err)
			},
		},
	}

	var svc, err = redis.NewService(redis.Options{Addr: testRedisAddr})
	require.NoError(t, err)

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			var err = svc.SetLocation(tc.args.ctx, tc.args.id, tc.args.l)
			tc.assert(t, tc, err)
		})
	}
}

func TestService_Locations(t *testing.T) {
	t.Skipf("WIP: Implement")
}
