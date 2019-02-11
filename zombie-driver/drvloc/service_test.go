package drvloc_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	sdrvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	mock "github.com/ifraixedes/go-ms-http-nsq-example/driver-location/mock"
	"github.com/ifraixedes/go-ms-http-nsq-example/internal/testassert"
	zmbdrv "github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver"
	"github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver/drvloc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.fraixed.es/errors"
)

func TestNewService(t *testing.T) {
	// TODO: ENHANCEMENT
	// Implement this test
	t.Skipf("TODO: TO BE IMPLEMENTED")
}

func TestService_IsZombie(t *testing.T) {
	var startT = time.Unix(time.Now().Unix()-5*60, 0)
	var locsUntil2MinsBefore = []sdrvloc.Location{
		{
			Lat: 42.005863,
			Lng: 2.292601,
			At:  startT,
		},
		{
			Lat: 42.00566,
			Lng: 2.292855,
			At:  startT.Add(time.Second * 30).Round(0),
		},
		{
			Lat: 42.005309,
			Lng: 2.292884,
			At:  startT.Add(time.Second * 60).Round(0),
		},
		{
			Lat: 42.004902,
			Lng: 2.292533,
			At:  startT.Add(time.Second * 90).Round(0),
		},
		{
			Lat: 41.991466,
			Lng: 2.280345,
			At:  startT.Add(time.Second * 120).Round(0),
		},
		{ // Stops here
			Lat: 41.961731,
			Lng: 2.266615,
			At:  startT.Add(time.Second * 150).Round(0),
		},
	}
	var locsAfter2MinsBefore = []sdrvloc.Location{
		{
			Lat: 41.961731,
			Lng: 2.266615,
			At:  startT.Add(time.Second * 180).Round(0),
		},
		{
			Lat: 41.961731,
			Lng: 2.266615,
			At:  startT.Add(time.Second * 210).Round(0),
		},
		{
			Lat: 41.961731,
			Lng: 2.266615,
			At:  startT.Add(time.Second * 240).Round(0),
		},
		{
			Lat: 41.961731,
			Lng: 2.266615,
			At:  startT.Add(time.Second * 270).Round(0),
		},
	}

	t.Run("ok: it is not a zombie", func(t *testing.T) {
		t.Parallel()
		var (
			ctx           = context.Background()
			dID           = uint64(rand.Int63n(99) + 1)
			dlsvc, finish = getDLSvcMock(t)
		)
		defer finish()

		var rules = zmbdrv.Rules{
			MinDistance: 100,
			LastMinutes: 5,
		}
		var svc, err = drvloc.NewService(dlsvc, rules)
		require.NoError(t, err)

		var locs []sdrvloc.Location
		locs = append(locs, locsUntil2MinsBefore...)
		locs = append(locs, locsAfter2MinsBefore...)

		dlsvc.EXPECT().LocationsForLastMinutes(ctx, dID, rules.LastMinutes).Return(locs, nil)
		ok, err := svc.IsZombie(ctx, dID)
		require.NoError(t, err)

		assert.False(t, ok)
	})

	t.Run("ok: it is a zombie", func(t *testing.T) {
		var (
			ctx           = context.Background()
			dID           = uint64(rand.Int63n(99) + 1)
			dlsvc, finish = getDLSvcMock(t)
		)
		defer finish()

		var rules = zmbdrv.Rules{
			MinDistance: 10,
			LastMinutes: 2,
		}
		var svc, err = drvloc.NewService(dlsvc, rules)
		require.NoError(t, err)

		dlsvc.EXPECT().LocationsForLastMinutes(ctx, dID, rules.LastMinutes).Return(locsAfter2MinsBefore, nil)
		ok, err := svc.IsZombie(ctx, dID)
		require.NoError(t, err)

		assert.True(t, ok)
	})

	t.Run("error: driver not found", func(t *testing.T) {
		var (
			ctx           = context.Background()
			dID           = uint64(rand.Int63n(99) + 1)
			dlsvc, finish = getDLSvcMock(t)
		)
		defer finish()

		var rules = zmbdrv.Rules{
			MinDistance: 8000,
			LastMinutes: 90,
		}
		var svc, err = drvloc.NewService(dlsvc, rules)
		require.NoError(t, err)

		dlsvc.EXPECT().LocationsForLastMinutes(ctx, dID, rules.LastMinutes).
			Return(nil, errors.New(sdrvloc.ErrNotFoundDriver))

		_, err = svc.IsZombie(ctx, dID)
		testassert.ErrorWithCode(t, err, zmbdrv.ErrNotFoundDriver, errors.MD{K: "id", V: dID})
	})

	t.Run("error: unexpected error", func(t *testing.T) {
		t.Parallel()
		var (
			ctx           = context.Background()
			dID           = uint64(rand.Int63n(99) + 1)
			dlsvc, finish = getDLSvcMock(t)
		)
		defer finish()

		var rules = zmbdrv.Rules{
			MinDistance: 1000,
			LastMinutes: 20,
		}
		var svc, err = drvloc.NewService(dlsvc, rules)
		require.NoError(t, err)

		dlsvc.EXPECT().LocationsForLastMinutes(ctx, dID, rules.LastMinutes).
			Return(nil, errors.New(sdrvloc.ErrInvalidDataFormatStore))

		_, err = svc.IsZombie(ctx, dID)
		testassert.ErrorWithCode(t, err, zmbdrv.ErrUnexpectedErrorSystem)
	})

	t.Run("error: context cancelation", func(t *testing.T) {
		t.Parallel()
		var (
			ctx, cancel   = context.WithCancel(context.Background())
			dID           = uint64(rand.Int63n(99) + 1)
			dlsvc, finish = getDLSvcMock(t)
		)
		defer finish()
		cancel()

		var rules = zmbdrv.Rules{
			MinDistance: 50,
			LastMinutes: 1,
		}
		var svc, err = drvloc.NewService(dlsvc, rules)
		require.NoError(t, err)

		dlsvc.EXPECT().LocationsForLastMinutes(ctx, dID, rules.LastMinutes).
			Return(nil, ctx.Err())

		_, err = svc.IsZombie(ctx, dID)
		testassert.ErrorWithCode(t, err, zmbdrv.ErrAbortedCtx)
	})
}

func getDLSvcMock(t *testing.T) (dlsvc *mock.Service, finish func()) {
	var ctrl = gomock.NewController(t)
	dlsvc = mock.NewService(ctrl)

	return dlsvc, ctrl.Finish
}
