package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ifraixedes/go-ms-http-nsq-example/gateway"
)

func main() {
	var in = parseInput()
	var cb, err = ioutil.ReadAll(in.cfg)
	if err != nil {
		exit(fmt.Errorf("error when reading the configuration file. %+v", err))
	}

	cfg, err := gateway.NewConfigFromYAML(cb)
	if err != nil {
		exit(fmt.Errorf("invalid configuration. %+v", err))
	}

	cfg.NSQdAddr = in.nsqdAddr
	svr, err := gateway.NewGateway(cfg)
	if err != nil {
		exit(fmt.Errorf("error when creating the gateway. %+v", err))
	}

	svr.Addr = in.httpAddr
	_ = svr.ListenAndServe()
}

type input struct {
	cfg      *os.File
	nsqdAddr string
	httpAddr string
}

func parseInput() *input {
	var (
		cfgfp    = flag.String("c", "", "the configuration YAML file path")
		nsqdAddr = flag.String("nsqdAddr", "127.0.0.1:4150", "NSQd network address")
		httpAddr = flag.String("httpAddr", "127.0.0.1:9002", "the address where the HTTP server will listen (host:port)")
	)

	flag.Parse()

	if *cfgfp == "" {
		exit(errors.New("Config file path must be indicated"))
	}

	f, err := os.Open(*cfgfp)
	if err != nil {
		perr, ok := err.(*os.PathError)
		if ok {
			exit(fmt.Errorf("Error while opening the YAML file (%s): %s", perr.Path, perr.Err.Error()))
		} else {
			exit(fmt.Errorf("Error while opening the YAML file (%s): %s", *cfgfp, err.Error()))
		}
	}

	return &input{
		cfg:      f,
		nsqdAddr: *nsqdAddr,
		httpAddr: *httpAddr,
	}

}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "There has been an error.\n%s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
