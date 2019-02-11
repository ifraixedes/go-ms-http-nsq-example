package internal

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dimfeld/httptreemux"
	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"go.fraixed.es/errors"
)

type locationsForLastMinsHanlder struct {
	svc drvloc.Service
}

func (h locationsForLastMinsHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		params = httptreemux.ContextParams(r.Context())
		idp    = params["id"]
		minsq  = r.FormValue("minutes")
	)

	if minsq == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var mins, err = strconv.ParseUint(minsq, 10, 16)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	id, err := strconv.ParseUint(idp, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var ctx = r.Context()
	ls, err := h.svc.LocationsForLastMinutes(ctx, id, uint16(mins))
	if err != nil {
		if errors.Is(err, drvloc.ErrNotFoundDriver) {
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

	b, err := json.Marshal(ls)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}
