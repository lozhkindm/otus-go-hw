package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout time.Duration
)

func main() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Client timeout")
	flag.Parse()

	if len(os.Args) < 3 {
		log.Fatal("host and port required")
	}
	addr := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT)
	defer signal.Stop(stopCh)

	clientCh := make(chan error)
	go func() {
		clientCh <- client.Send()
	}()
	go func() {
		clientCh <- client.Receive()
	}()

	select {
	case <-stopCh:
	case err := <-clientCh:
		if err != nil {
			log.Fatal(err)
		}
	}
}
