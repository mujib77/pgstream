package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mujib77/pgstream/config"
	"github.com/mujib77/pgstream/internal/connector"
	"github.com/mujib77/pgstream/internal/decoder"
	"github.com/mujib77/pgstream/internal/handler"
)

func main() {
	cfg := config.DefaultConfig()

	fmt.Println("connecting to postgres...")
	conn, err := connector.New(cfg.DatabaseURL, cfg.SlotName, cfg.PublicationName)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("connected!")

	fmt.Println("creating replication slot...")
	err = conn.CreateSlot()
	if err != nil {
		fmt.Println("slot may already exist, continuing...")
	}

	fmt.Println("starting replication...")
	err = conn.Start()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	dec := decoder.New(conn.GetConn())
	han := handler.New()

	fmt.Println("listening for changes...")
	for {
		event, err := dec.NextEvent(context.Background())
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		err = han.Handle(event)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}