//Package drvloc is an implementation of the Zombie Driver Service which uses
// the Driver Location Service to get the location of the drivers.
package drvloc

import (
	"context"

	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	zmbdrv "github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver"
	"github.com/umahmood/haversine"
	"go.fraixed.es/errors"
)

// NewService creates an instance of the Service type using the passed
// configured Driver Location Service client and the r.
//
// The following error codes can be returned:
//
// * ErrRequiredDrvLocSvc
//
// * ErrInvalidRules
func NewService(dlsvc drvloc.Service, r zmbdrv.Rules) (*Service, error) {
	if dlsvc == nil {
		return nil, errors.New(ErrRequiredDrvLocSvc)
	}

	if !r.IsValid() {
		return nil, errors.New(ErrInvalidRules, errors.MD{K: "arg:r", V: r})
	}

	return &Service{
		dlsvc: dlsvc,
		rules: r,
	}, nil
}

// Service implements the zmbdrv.Service interface using a drvloc.Service
//client.
type Service struct {
	dlsvc drvloc.Service
	rules zmbdrv.Rules
}

// IsZombie verifies if a driver with id is a zombie or not based on the
// locations that drvloc.Service returns.
func (s *Service) IsZombie(ctx context.Context, id uint64) (bool, error) {
	var locs, err = s.dlsvc.LocationsForLastMinutes(ctx, id, s.rules.LastMinutes)
	if err != nil {
		switch {
		case errors.Is(err, drvloc.ErrNotFoundDriver):
			return false, errors.New(zmbdrv.ErrNotFoundDriver, errors.MD{K: "id", V: id})
		case err == ctx.Err():
			return false, errors.New(zmbdrv.ErrAbortedCtx)
		default:
			return false, errors.Wrap(err, zmbdrv.ErrUnexpectedErrorSystem)
		}
	}

	if len(locs) < 2 {
		return true, nil
	}

	var (
		pdis = locs[0]
		tdis uint64
	)
	for _, l := range locs[1:] {
		var (
			p = haversine.Coord{Lat: pdis.Lat, Lon: pdis.Lng}
			n = haversine.Coord{Lat: l.Lat, Lon: l.Lng}
		)

		var _, km = haversine.Distance(p, n)
		tdis += uint64(km * 1000)

		if tdis >= s.rules.MinDistance {
			return false, nil
		}

		pdis = l
	}

	return true, nil
}
