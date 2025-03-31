package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func init() {
	flag.StringVar(&nickname, "nickname", "notbot", "bot nickname")
	flag.StringVar(&listen, "listen", "127.0.0.1:6789", "listen address")
	flag.Var(&jitsiChannels, "jitsi.channels", "jitsiServer/jitsiRoom mapping; may be specified multiple times")
}

func main() {
	flag.Parse()

	log.Println("starting api server on:", listen)
	apiServer := StartAPIServer()

	log.Println("monitoring channels:", jitsiChannels)
	jitsiDone := JitsiRunWrapper(apiServer)

	log.Println("running...")
	sig := waitForSignal()
	log.Println("received signal", sig, ", shutting down")

	log.Println("disconnecting from jitsi servers")
	for _, ch := range jitsiDone {
		ch <- true
	}

	log.Println("stopping the api server")
	if err := apiServer.h.Close(); err != nil {
		log.Println("error while stopping the http server, ¯\\_(ツ)_/¯", err)
	}

	log.Println("jitsi-monitor has been stopped")
}

func waitForSignal() os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGHUP)
	for {
		sig := <-ch
		signal.Stop(ch)
		return sig
	}
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	jitsiChannels arrayFlags
	nickname      string
	listen        string
)
