package gateway

import (
	"errors"
	"net/http"

	"github.com/dimfeld/httptreemux"
	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location/nsqhttp"
	zmbdrv "github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver/http"
	nsq "github.com/nsqio/go-nsq"
)

// NewGateway creates a new HTTP Server which acts as a gateway with the
// specified c.
//
// The gateway uses the nsqhttp client of the Driver Location Service and the
// http client of the Zombie Drier Service for operating with them.
//
// The returned server is only configured with the appropriated gateway
// handlers, not any other configuration is set.
func NewGateway(c Config) (*http.Server, error) {
	if c == (Config{}) {
		return nil, errors.New("Invalid configuration")
	}

	var dlsvc, err = drvloc.NewClientNSQ(c.NSQdAddr, nsq.NewConfig(), c.setDriverLocation.NSQ.Topic)
	if err != nil {
		return nil, err
	}

	zdsvc, err := zmbdrv.NewClient(c.getDriver.HTTP.Host)
	if err != nil {
		return nil, err
	}

	var eps = &endpoints{
		dlsvc: dlsvc,
		zdsvc: zdsvc,
	}

	return &http.Server{
		Handler: routeEndpoints(c, eps),
	}, nil
}

func routeEndpoints(c Config, eps *endpoints) http.Handler {
	var router = httptreemux.NewContextMux()

	switch c.setDriverLocation.Method {
	case http.MethodPost:
		router.POST(c.setDriverLocation.Path, eps.setDriverLocation)
	case http.MethodPut:
		router.PUT(c.setDriverLocation.Path, eps.setDriverLocation)
	case http.MethodPatch:
		router.PATCH(c.setDriverLocation.Path, eps.setDriverLocation)
	default:
		panic("there is a bug in some of the Config contructor " +
			"functions, setDriverLocation has an unsupported HTTP Method",
		)
	}

	switch c.getDriver.Method {
	case http.MethodGet:
		router.GET(c.getDriver.Path, eps.getDriver)
	case http.MethodPost:
		router.POST(c.getDriver.Path, eps.getDriver)
	case http.MethodPut:
		router.PUT(c.getDriver.Path, eps.getDriver)
	case http.MethodPatch:
		router.PATCH(c.getDriver.Path, eps.getDriver)
	default:
		panic("there is a bug in some of the Config contructor " +
			"functions, getDriver has an unsupported HTTP Method",
		)
	}

	return router
}
