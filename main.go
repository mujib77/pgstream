package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mujib77/pgstream/config"
	"github.com/mujib77/pgstream/internal/connector"
	"github.com/mujib77/pgstream/internal/decoder"
	"github.com/mujib77/pgstream/internal/handler"
)

func main() {
	cfg := config.DefaultConfig()

	fmt.Println("connecting to postgres...")
	ctx := context.Background()
	conn, err := connector.New(ctx, cfg.DatabaseURL, cfg.SlotName, cfg.PublicationName)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	fmt.Println("connected!")

	fmt.Println("creating replication slot...")
	err = conn.CreateSlot(ctx)
	if err != nil {
		fmt.Println("slot may already exist, continuing...")
	}

	fmt.Println("starting replication...")
	err = conn.Start(ctx)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	dec := decoder.New(conn.GetConn())
	han := handler.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		fmt.Println("\nshutting down...")
		cancel()
	}()

	fmt.Println("listening for changes...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			event, err := dec.NextEvent(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				fmt.Println("error:", err)
				return
			}
			err = han.Handle(event)
			if err != nil {
				fmt.Println("error:", err)
			}
		}
	}
}