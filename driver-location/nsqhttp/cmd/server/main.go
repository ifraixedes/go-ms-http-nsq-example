package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/nsqhttp/internal"
	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/redis"
	nsq "github.com/nsqio/go-nsq"
)

func main() {
	var in = parseInput()
	var svc, err = redis.NewService(in.Redis)
	if err != nil {
		exit(err)
	}

	down, err := internal.ServerUp(svc, in.NSQ, in.HTTP, log.New(os.Stderr, "", log.Ldate))
	if err != nil {
		exit(err)
	}

	var termChan = make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	err = down(context.Background())
	if err != nil {
		exit(err)
	}
}

type input struct {
	NSQ   internal.NSQSettings
	HTTP  internal.HTTPConfig
	Redis redis.Options
}

func parseInput() *input {
	var (
		redisAddr       = flag.String("redisAddr", "127.0.0.1:6379", "Redis address (host:port)")
		nsqLookupdAddrs = flag.String("nsqLookupdAddrs", "127.0.0.1:4161", "comma separated list of NSQ Lookupd addresses")
		nsqTopic        = flag.String("nsqTopic", "locations", "the NSQL topic to be consumed")
		nsqChan         = flag.String("nsqChan", "driver-location-svc", "the NSQL channel to be consumed")
		httpAddr        = flag.String("httpAddr", "127.0.0.1:9000", "the address where the HTTP server will listen (host:port)")
	)

	flag.Parse()

	return &input{
		NSQ: internal.NSQSettings{
			LookupdAddrs: strings.Split(*nsqLookupdAddrs, ","),
			Topic:        *nsqTopic,
			Channel:      *nsqChan,
			Cfg:          nsq.NewConfig(),
		},
		HTTP: internal.HTTPConfig{
			Addr: *httpAddr,
		},
		Redis: redis.Options{Addr: *redisAddr},
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "There has been an error.\n%s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
