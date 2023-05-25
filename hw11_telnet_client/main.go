package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "")
}

func read(ctx context.Context, client TelnetClient) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := client.Receive()
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}
}

func write(ctx context.Context, client TelnetClient) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := client.Send()
			if err != nil {
				fmt.Println(err)
			}

			return
		}
	}
}

func main() {
	flag.Parse()

	if len(os.Args) < 3 {
		fmt.Println("Необходимо передать host и port для подключения")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	hostAndPort := os.Args[len(os.Args)-2:]
	client := NewTelnetClient(net.JoinHostPort(hostAndPort[0], hostAndPort[1]), timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		read(ctx, client)
	}()

	go func() {
		defer cancel()
		defer wg.Done()
		write(ctx, client)
	}()

	wg.Wait()

	err = client.Close()
	if err != nil {
		fmt.Println(err)
	}
}
