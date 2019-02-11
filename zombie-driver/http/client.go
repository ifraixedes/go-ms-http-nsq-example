package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	zmbdrv "github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver"
	"github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver/http/internal"
	"go.fraixed.es/errors"
)

// NewClient creates a new client of the Zoombie Driver Service to operate with
// the service HTTP.
func NewClient(httpAddr string) (zmbdrv.Service, error) {
	// TODO: IMPROVEMENT
	// validate httpAddr and verify it, making a noop request to the server for
	// ensuring that the passed address is OK

	return &client{
		httpBase: url.URL{
			Scheme: "http",
			Host:   httpAddr,
		},
		httpClient: &http.Client{},
	}, nil
}

type client struct {
	httpBase   url.URL
	httpClient *http.Client
}

func (c *client) IsZombie(ctx context.Context, id uint64) (bool, error) {
	var u = c.httpBase
	u.Path = fmt.Sprintf("drivers/%d", id)

	var req, err = http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return false, errors.Wrap(err, zmbdrv.ErrUnexpectedErrorSystem)
	}

	req = req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	if err != nil {
		if err == ctx.Err() {
			return false, errors.Wrap(err, zmbdrv.ErrAbortedCtx)
		}
		return false, errors.Wrap(err, zmbdrv.ErrUnexpectedErrorSystem)
	}

	switch {
	case res.StatusCode == http.StatusOK:
	case res.StatusCode == http.StatusNotFound:
		return false, errors.New(zmbdrv.ErrNotFoundDriver)
	default:
		return false, errors.Wrap(err, zmbdrv.ErrUnexpectedErrorSystem)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, errors.Wrap(err, zmbdrv.ErrUnexpectedErrorSystem)
	}

	var isz internal.IsZombieResBody
	err = json.Unmarshal(b, &isz)
	if err != nil {
		return false, errors.Wrap(err, zmbdrv.ErrUnexpectedErrorSystem)
	}

	return isz.Zombie, nil
}
