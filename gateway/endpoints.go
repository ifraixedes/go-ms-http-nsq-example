package gateway

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	zmbdrv "github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver"
	"go.fraixed.es/errors"
)

type endpoints struct {
	dlsvc drvloc.Service
	zdsvc zmbdrv.Service
}

func (e *endpoints) setDriverLocation(w http.ResponseWriter, r *http.Request) {
	var (
		params = httptreemux.ContextParams(r.Context())
		idp    = params["id"]
	)

	id, err := strconv.ParseUint(idp, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var sdlocb setDriverLocationReqBody
	err = json.Unmarshal(b, &sdlocb)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var ctx = r.Context()
	err = e.dlsvc.SetLocation(ctx, id, drvloc.Location{
		Lat: sdlocb.Lat,
		Lng: sdlocb.Lng,
		At:  time.Now(),
	})
	if err != nil {
		if err == ctx.Err() {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *endpoints) getDriver(w http.ResponseWriter, r *http.Request) {
	var (
		params = httptreemux.ContextParams(r.Context())
		idp    = params["id"]
	)

	id, err := strconv.ParseUint(idp, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var ctx = r.Context()
	isz, err := e.zdsvc.IsZombie(ctx, id)
	if err != nil {
		if errors.Is(err, zmbdrv.ErrNotFoundDriver) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err == ctx.Err() {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(getDriverResBody{ID: id, Zombie: isz})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

type setDriverLocationReqBody struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

type getDriverResBody struct {
	ID     uint64 `json:"id"`
	Zombie bool   `json:"zombie"`
}
