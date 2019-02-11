package nsqhttp

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/nsqhttp/internal"
	"github.com/nsqio/go-nsq"
	"go.fraixed.es/errors"
)

// NewClient creates a new client of the Driver Location Service to operate
// with the service through NSQ and HTTP.
func NewClient(nsqAddr string, cfg *nsq.Config, topic string, httpAddr string) (drvloc.Service, error) {
	// TODO: IMPROVEMENT
	// validate httpAddr and verify it, making a noop request to the server for
	// ensuring that the passed address is OK

	if !nsq.IsValidTopicName(topic) {
		return nil, stderrors.New("Invalid topic")
	}

	var setLoc, err = nsq.NewProducer(nsqAddr, cfg)
	if err != nil {
		return nil, err
	}

	if err = setLoc.Ping(); err != nil {
		return nil, stderrors.New("Impossible to connect to the indicated NSQ address")
	}

	return &client{
		setLoc: setLoc,
		topic:  topic,
		httpBase: url.URL{
			Scheme: "http",
			Host:   httpAddr,
		},
		httpClient: &http.Client{},
	}, nil
}

type client struct {
	setLoc     *nsq.Producer
	topic      string
	httpBase   url.URL
	httpClient *http.Client
}

func (c *client) SetLocation(_ context.Context, id uint64, l drvloc.Location) error {
	var slmb = internal.SetLocationMsgBody{
		ID:  id,
		Loc: l,
	}

	var b, err = json.Marshal(slmb)
	if err != nil {
		return errors.Wrap(err, drvloc.ErrUnexpectedErrorSerializer)
	}

	if err = c.setLoc.Publish(c.topic, b); err != nil {
		return errors.Wrap(err, drvloc.ErrUnexpectedErrorTransport)
	}

	return nil
}

func (c *client) LocationsForLastMinutes(ctx context.Context, id uint64, minutes uint16) ([]drvloc.Location, error) {
	var qv = url.Values{}
	qv.Set("minutes", strconv.Itoa(int(minutes)))

	var u = c.httpBase
	u.Path = fmt.Sprintf("drivers/%d/locations", id)
	u.RawQuery = qv.Encode()

	var req, err = http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, drvloc.ErrUnexpectedErrorSystem)
	}

	req = req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	if err != nil {
		if err == ctx.Err() {
			return nil, errors.Wrap(err, drvloc.ErrAbortedCtx)
		}
		return nil, errors.Wrap(err, drvloc.ErrUnexpectedErrorSystem)
	}

	switch {
	case res.StatusCode == http.StatusOK:
	case res.StatusCode == http.StatusNotFound:
		return nil, errors.New(drvloc.ErrNotFoundDriver)
	default:
		return nil, errors.Wrap(err, drvloc.ErrUnexpectedErrorSystem)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, drvloc.ErrUnexpectedErrorSystem)
	}

	var locs []drvloc.Location
	err = json.Unmarshal(b, &locs)
	if err != nil {
		return nil, errors.Wrap(err, drvloc.ErrUnexpectedErrorSystem)
	}

	return locs, nil
}
