package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/nsqhttp"
	zmbdrv "github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver"
	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver/drvloc"
	"github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver/http/internal"
)

func main() {
	var in = parseInput()
	var dlsvc, err = nsqhttp.NewClientHTTP(in.DrvlocHTTPAddr)
	if err != nil {
		exit(err)
	}

	svc, err := drvloc.NewService(dlsvc, in.Rules)
	if err != nil {
		exit(err)
	}

	down, err := internal.ServerUp(svc, in.HTTP)
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
	HTTP           internal.HTTPConfig
	DrvlocHTTPAddr string
	Rules          zmbdrv.Rules
}

func parseInput() *input {
	var (
		httpAddr = flag.String(
			"httpAddr", "127.0.0.1:9001", "the address where the HTTP server will listen (host:port)",
		)
		drvlocHTTPAddr = flag.String(
			"drvlocHTTPAddr", "127.0.0.1:9000", "the address where the Driver Location Service is running over HTTP (host:port)",
		)
		rMinDis = flag.Uint64(
			"rulesMinDistance", 500, "the minium distance that a driver must drive for not being considered a zombie",
		)
		rLastMins = flag.Uint(
			"rulesLastMin", 5, "the number of minutes which are considered to check the distance that a driver has driven",
		)
	)

	flag.Parse()

	if *rLastMins > math.MaxUint16 {
		exit(fmt.Errorf("rulesLastMin cannot be greater than %d", math.MaxUint16))
	}

	return &input{
		HTTP: internal.HTTPConfig{
			Addr: *httpAddr,
		},
		DrvlocHTTPAddr: *drvlocHTTPAddr,
		Rules: zmbdrv.Rules{
			MinDistance: *rMinDis,
			LastMinutes: uint16(*rLastMins),
		},
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "There has been an error.\n%s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
