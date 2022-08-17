package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bit0rez/kafka-tools/command/consume"
	"github.com/bit0rez/kafka-tools/command/produce"
	"github.com/bit0rez/kafka-tools/command/topic"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	_ = godotenv.Load()

	app := cli.App{
		Authors: []*cli.Author{
			{
				Name:  "Merkulov Ilia aka bit0rez",
				Email: "ilia@merkulov.online",
			},
		},
		Name:                 "Kafka tools",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			&produce.Produce,
			&consume.Consume,
			&topic.Topic,
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
