package main

import (
	"context"
	_ "embed"
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

//go:embed completion/bash
var bashCompletion string

//go:embed completion/zsh
var zshCompletion string

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
			{
				Name:   "completion",
				Usage:  "autocompletion for Bash and Zsh",
				Action: completion,
			},
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func completion(ctx *cli.Context) error {
	shell, ok := os.LookupEnv("SHELL")
	if !ok {
		return nil
	}

	switch shell {
	case "/bin/bash":
		fmt.Print(bashCompletion)
		return nil
	case "/bin/zsh":
		fmt.Print(zshCompletion)
		return nil
	default:
		return nil
	}
}
