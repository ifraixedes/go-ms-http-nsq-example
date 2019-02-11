package internal

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/dimfeld/httptreemux"
	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	nsq "github.com/nsqio/go-nsq"
)

// DownFunc is a function that when it's called, it releases all the resources
// and stop the HTTP server.
type DownFunc func(context.Context) error

// NSQSettings are the required settings for NSQ which the server will use.
type NSQSettings struct {
	Topic        string
	Channel      string
	Cfg          *nsq.Config
	LookupdAddrs []string
}

// HTTPConfig is the HTTP server configuration that the server will use.
type HTTPConfig struct {
	Addr string
}

// ServerUp spins up the NSQ consumer and HTTP server for serving the Driver
// Location Service functionalities.
//
// DownFunc can be called more than once, however only the first call do the
// shutting down operations, the following calls will return a nil error if the
// first call has not finished otherwise the returned error if there was one.
// NOTE that calling the DownFunc concurrently more than once before the first
// call has finished may produce a data race.
//
// Error is returned if the server cannot be up for whatever reason.
func ServerUp(svc drvloc.Service, ns NSQSettings, hc HTTPConfig, l *log.Logger) (DownFunc, error) {
	if svc == nil {
		return nil, errors.New("Driver Location Service instance cannot be nil")
	}

	var consumer, err = createNSQConsumer(svc, ns.Topic, ns.Channel, ns.Cfg, ns.LookupdAddrs)
	if err != nil {
		return nil, err
	}

	var svr = createHTTPServer(svc, hc)

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

		consumer.Stop()
		downErr = svr.Shutdown(ctx)

		<-consumer.StopChan
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
	case err = <-srvErr:
		return nil, err
	}

	return df, nil
}

func createNSQConsumer(
	svc drvloc.Service, topic string, channel string, cfg *nsq.Config, lookupdAddrs []string,
) (*nsq.Consumer, error) {
	if len(lookupdAddrs) == 0 {
		return nil, errors.New("At least one lookupd addres is required for NSQ consumer")
	}

	var consumer, err = nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return nil, err
	}

	var h = setLocationHandler{
		svc: svc,
	}

	var cc = runtime.GOMAXPROCS(0)
	if cc > 3 {
		consumer.AddConcurrentHandlers(h, cc-2)
	} else {
		consumer.AddHandler(h)
	}

	if err = consumer.ConnectToNSQLookupds(lookupdAddrs); err != nil {
		return nil, err
	}

	return consumer, nil
}

func createHTTPServer(svc drvloc.Service, c HTTPConfig) *http.Server {
	var (
		h = locationsForLastMinsHanlder{
			svc: svc,
		}
		router = httptreemux.NewContextMux()
	)

	router.GET("/drivers/:id/locations", h.ServeHTTP)

	return &http.Server{
		Addr:    c.Addr,
		Handler: router,
	}
}
