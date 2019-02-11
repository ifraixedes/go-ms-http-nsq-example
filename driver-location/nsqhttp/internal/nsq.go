package internal

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"github.com/nsqio/go-nsq"
)

var (
	errNSQMsgEmpty       = errors.New("the message received in the NSQL consumer is nil or its body is empty")
	errNSQMsgBodyInvalid = errors.New("the message received in the NSQL consumer has an invalid body content")
)

// SetLocationMsgBody is the NSQ Message body type for sending a driver location.
type SetLocationMsgBody struct {
	ID  uint64          `json:"id"`
	Loc drvloc.Location `json:"location"`
}

type setLocationHandler struct {
	svc drvloc.Service
}

func (h setLocationHandler) HandleMessage(msg *nsq.Message) error {
	if msg == nil || len(msg.Body) == 0 {
		return errNSQMsgEmpty
	}

	var slmb SetLocationMsgBody
	if err := json.Unmarshal(msg.Body, &slmb); err != nil {
		return errNSQMsgBodyInvalid
	}

	if err := h.svc.SetLocation(context.Background(), slmb.ID, slmb.Loc); err != nil {
		return err
	}

	return nil
}

// nsqlLoggerMiddleware wraps a NSQ Handler for logging errors which aren't
// returned.
// The errors are not returned because re-queueing the message won't solve the
// problem.
type nsqlLoggerMiddleware struct {
	l *log.Logger
	h nsq.HandlerFunc
}

func (nlm nsqlLoggerMiddleware) HandleMessage(msg *nsq.Message) error {
	var err error
	switch err = nlm.h(msg); err {
	case errNSQMsgEmpty, errNSQMsgBodyInvalid:
		nlm.l.Printf("%v", err)
		return nil
	}

	return err
}
