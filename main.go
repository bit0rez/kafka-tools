package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bit0rez/kafka-tools/commands"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	_ = godotenv.Load()

	app := cli.App{
		Name: "Kafka tools",
		Commands: []*cli.Command{
			&commands.Produce,
			&commands.Consume,
			&commands.Topic,
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
