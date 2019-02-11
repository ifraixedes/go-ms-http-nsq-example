package internal

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
	zmbdrv "github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver"
	"go.fraixed.es/errors"
)

// DownFunc is a function that when it's called, it releases all the resources
// and stop the HTTP server.
type DownFunc func(context.Context) error

// HTTPConfig is the HTTP server configuration that the server will use.
type HTTPConfig struct {
	Addr string
}

// ServerUp runs an HTTP server listening in the httpAddr which expose though
// the HTTP transport the functionality of the passed Zombie Driver service.
func ServerUp(svc zmbdrv.Service, c HTTPConfig) (DownFunc, error) {
	if svc == nil {
		return nil, stderrors.New("Zombie Driver Service cannot be nil")
	}

	var (
		h = isZombieHandler{
			svc: svc,
		}
		router = httptreemux.NewContextMux()
	)

	router.GET("/drivers/:id", h.ServeHTTP)
	var svr = &http.Server{
		Addr:    c.Addr,
		Handler: router,
	}

	var (
		downErr    error
		downCalled bool
	)
	var df = func(ctx context.Context) error {
		if !downCalled {
			downCalled = true
		} else {
			return downErr
		}

		downErr = svr.Shutdown(ctx)
		return downErr
	}

	// TODO: IMPROVEMENT
	// This mechanism of waiting 2 seconds for assuming that the HTTP server isn't
	// good, so it must be replace for once that we can know better that the HTTP
	// server started to listen without any issue.
	var (
		srvErr = make(chan error)
		tm     = time.NewTimer(2 * time.Second)
	)
	go func() {
		srvErr <- svr.ListenAndServe()
	}()

	select {
	case <-tm.C:
	case err := <-srvErr:
		return nil, err
	}

	return df, nil
}

// IsZombieResBody is the HTTP response body format of the IsZombien service
// function
type IsZombieResBody struct {
	ID     uint64 `json:"id"`
	Zombie bool   `json:"zombie"`
}

type isZombieHandler struct {
	svc zmbdrv.Service
}

func (h isZombieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	isz, err := h.svc.IsZombie(ctx, id)
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

	b, err := json.Marshal(IsZombieResBody{ID: id, Zombie: isz})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}
