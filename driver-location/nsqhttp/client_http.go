package nsqhttp

import (
	"context"
	"net/http"
	"net/url"

	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
)

// NewClientHTTP creates a new client of the Driver Location Service which only
// offers the functionalities exposed over HTTP, and it panics if any of the
// functionalities exposed over NSQ is called.
//
// NOTE only use this client if you DON'T NEED any Driver Location Service
// function exposed over NSQ.
func NewClientHTTP(httpAddr string) (drvloc.Service, error) {
	// TODO: IMPROVEMENT
	// validate httpAddr and verify it, making a noop request to the server for
	// ensuring that the passed address is OK

	return &clientHTTP{
		client: &client{
			httpBase: url.URL{
				Scheme: "http",
				Host:   httpAddr,
			},
			httpClient: &http.Client{},
		},
	}, nil
}

type clientHTTP struct {
	*client
}

func (c *clientHTTP) SetLocation(_ context.Context, id uint64, l drvloc.Location) error {
	panic("Driver Location Service function exposed over NSQ MUST NOT be executed with this client")
}
