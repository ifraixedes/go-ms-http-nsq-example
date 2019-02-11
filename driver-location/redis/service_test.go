package redis_test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/redis"
	"github.com/ifraixedes/go-ms-http-nsq-example/internal/testassert"
	"github.com/mmcloughlin/spherand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.fraixed.es/errors"
)

func TestNewService(t *testing.T) {
	type tcase struct {
		desc   string
		o      redis.Options
		assert func(*testing.T, tcase, *redis.Service, error)
	}

	var tcases = []tcase{
		{
			desc: "ok",
			o: redis.Options{
				Addr: testRedisAddr,
			},
			assert: func(t *testing.T, _ tcase, s *redis.Service, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, s)
			},
		},
		{
			desc: "error: invalid connection",
			o: redis.Options{
				Addr: "google.com:6379",
			},
			assert: func(t *testing.T, tc tcase, _ *redis.Service, err error) {
				testassert.ErrorWithCode(t, err, redis.ErrInvalidRedisConfig, errors.MD{K: "config", V: tc.o})
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			var s, err = redis.NewService(tc.o)
			if s != nil {
				defer require.NoError(t, s.Close(context.Background()))
			}

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

func TestService_LocationsForLastMinutes(t *testing.T) {
	var svc *redis.Service
	{
		var err error
		svc, err = redis.NewService(redis.Options{Addr: testRedisAddr})
		require.NoError(t, err)
	}

	t.Run("error: NotFoundDriver", func(t *testing.T) {
		t.Parallel()

		var (
			dID     = uint64(rand.Int63n(100) + 100)
			minutes = uint16(rand.Intn(5) + 1)
		)

		var _, err = svc.LocationsForLastMinutes(context.Background(), dID, minutes)
		testassert.ErrorWithCode(t, err, drvloc.ErrNotFoundDriver, errors.MD{K: "id", V: dID})
	})

	t.Run("ok location", func(t *testing.T) {
		t.Parallel()

		// Expected values
		var (
			ctx   = context.Background()
			dID   = uint64(rand.Int63n(100) + 200)
			locs  []drvloc.Location
			nlocs = rand.Intn(60) + 30
			// due to the fact that this base time is calculate in advance than the base
			// time used by the service function, +1 second is added for avoiding
			// discrepancies in the expected result
			baseTime = time.Unix(time.Now().Unix()+1-int64(nlocs*5), 0)
		)

		for i := 0; i < nlocs; i++ {
			var (
				lat, lng = spherand.Geographical()
				loc      = drvloc.Location{
					Lat: lat,
					Lng: lng,
					At:  time.Unix(baseTime.Unix()+int64(i*5), 0).Round(0),
				}
			)
			locs = append(locs, loc)

			var err = svc.SetLocation(ctx, dID, loc)
			require.NoError(t, err)
		}

		// Insert some repeated locations for ensuring that these one won't be present
		for i := rand.Intn(10) + (nlocs - 10); i < nlocs; i++ {
			var err = svc.SetLocation(ctx, dID, locs[i])
			require.NoError(t, err)
		}

		t.Run("some", func(t *testing.T) {
			t.Parallel()

			// Take 1 or 2 minutes randomly and calculate the index of the first
			// location that it should be returned
			var (
				mins  = uint16(rand.Intn(1) + 1)
				flIdx = nlocs - (60 / 5 * int(mins))
			)

			var ls, err = svc.LocationsForLastMinutes(ctx, dID, mins)
			require.NoError(t, err)

			assert.Equal(t, locs[flIdx:], ls)
		})

		t.Run("all", func(t *testing.T) {
			t.Parallel()

			var mins = uint16(nlocs*5/60) + 1

			var ls, err = svc.LocationsForLastMinutes(ctx, dID, mins)
			require.NoError(t, err)

			assert.Equal(t, locs, ls)
		})
	})

	t.Run("ok: no locations", func(t *testing.T) {
		t.Parallel()

		var (
			ctx = context.Background()
			dID = uint64(rand.Int63n(100) + 300)
		)

		var loc drvloc.Location
		{
			lat, lng := spherand.Geographical()
			loc = drvloc.Location{
				Lat: lat,
				Lng: lng,
				At:  time.Unix(time.Now().Unix()-120, 0), // it's 2 minutes ago
			}
		}

		var err = svc.SetLocation(ctx, dID, loc)
		require.NoError(t, err)

		ls, err := svc.LocationsForLastMinutes(ctx, dID, 1)
		require.NoError(t, err)

		assert.Empty(t, ls)
	})
}

func TestService_Close(t *testing.T) {
	t.Run("error: Close after close", func(t *testing.T) {
		var svc, err = redis.NewService(redis.Options{Addr: testRedisAddr})
		require.NoError(t, err)

		var ctx = context.Background()
		err = svc.Close(ctx)
		require.NoError(t, err)

		err = svc.Close(ctx)
		testassert.ErrorWithCode(t, err, redis.ErrClosedService)
	})

	t.Run("error: SetLocation after close", func(t *testing.T) {
		var svc, err = redis.NewService(redis.Options{Addr: testRedisAddr})
		require.NoError(t, err)

		var ctx = context.Background()
		err = svc.Close(ctx)
		require.NoError(t, err)

		err = svc.SetLocation(ctx, 1000, drvloc.Location{
			Lat: 41.34617,
			Lng: -118.50037,
			At:  time.Now(),
		})
		testassert.ErrorWithCode(t, err, redis.ErrClosedService)
	})

	t.Run("error: LocationsForLastSeconds after close", func(t *testing.T) {
		var svc, err = redis.NewService(redis.Options{Addr: testRedisAddr})
		require.NoError(t, err)

		var ctx = context.Background()
		err = svc.Close(ctx)
		require.NoError(t, err)

		_, err = svc.LocationsForLastMinutes(ctx, 1000, 1)
		testassert.ErrorWithCode(t, err, redis.ErrClosedService)
	})

	t.Run("ok: Close call concurrently", func(t *testing.T) {
		var svc *redis.Service
		{
			var err error
			svc, err = redis.NewService(redis.Options{Addr: testRedisAddr})
			require.NoError(t, err)
		}

		var (
			ctx = context.Background()
			wg  sync.WaitGroup
		)

		wg.Add(1)
		var closeCall = func() {
			var err = svc.Close(ctx)
			assert.NoError(t, err)
			wg.Done()
		}

		wg.Add(1)
		var svcCall = func() {
			var err = svc.SetLocation(ctx, 1000, drvloc.Location{
				Lat: 41.34617,
				Lng: -118.50037,
				At:  time.Now(),
			})

			if err != nil {
				fmt.Printf("\n\n%+v\n\n", err)
				testassert.ErrorWithCode(t, err, redis.ErrClosedService)
			}

			wg.Done()
		}

		if rand.Int()%2 == 0 {
			go svcCall()
			go closeCall()
		} else {
			go closeCall()
			go svcCall()
		}

		wg.Wait()
	})
}
