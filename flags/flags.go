package flags

import "github.com/urfave/cli/v2"

var (
	BootstrapFlag = cli.StringSliceFlag{
		Name:        "bootstrap-server",
		Usage:       "List of Kafka servers",
		Aliases:     []string{"b"},
		EnvVars:     []string{"BOOTSTRAP_SERVER"},
		DefaultText: "localhost:9092",
		Value:       cli.NewStringSlice("localhost:9092"),
	}

	TopicFlag = cli.StringFlag{
		Name:     "topic",
		Usage:    "Topic name",
		Aliases:  []string{"t"},
		EnvVars:  []string{"TOPIC"},
		Required: true,
	}

	TopicsFlag = cli.StringSliceFlag{
		Name:     "topic",
		Usage:    "List of topics",
		Aliases:  []string{"t"},
		EnvVars:  []string{"TOPIC"},
		Required: true,
	}

	HeadersFlag = cli.StringSliceFlag{
		Name:    "header",
		Usage:   "Message header. Example: header=value",
		Aliases: []string{"H"},
		EnvVars: []string{"HEADER"},
	}
)
