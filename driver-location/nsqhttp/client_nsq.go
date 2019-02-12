package nsqhttp

import (
	"context"
	"errors"

	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"github.com/nsqio/go-nsq"
)

// NewClientNSQ creates a new client of the Driver Location Service which only
// offers the functionalities exposed over NSQ, and it panics if any of the
// functionalities exposed over HTTP is called.
//
// NOTE only use this client if you DON'T NEED any Driver Location Service
// function exposed over HTTP.
func NewClientNSQ(nsqAddr string, cfg *nsq.Config, topic string) (drvloc.Service, error) {
	if !nsq.IsValidTopicName(topic) {
		return nil, errors.New("Invalid topic")
	}

	var setLoc, err = nsq.NewProducer(nsqAddr, cfg)
	if err != nil {
		return nil, err
	}

	if err = setLoc.Ping(); err != nil {
		return nil, errors.New("Impossible to connect to the indicated NSQ address")
	}

	return &clientNSQ{
		client: &client{
			setLoc: setLoc,
			topic:  topic,
		},
	}, nil
}

type clientNSQ struct {
	*client
}

func (c *clientNSQ) LocationsForLastMinutes(ctx context.Context, id uint64, minutes uint16) ([]drvloc.Location, error) {
	panic("Driver Location Service function exposed over HTTP MUST NOT be executed with this client")
}
