package flags

import "github.com/urfave/cli/v2"

var (
	BootstrapFlag = cli.StringSliceFlag{
		Name:        "bootstrap-server",
		Aliases:     []string{"b"},
		EnvVars:     []string{"BOOTSTRAP_SERVER"},
		DefaultText: "localhost:9092",
		Value:       cli.NewStringSlice("localhost:9092"),
	}

	TopicFlag = cli.StringFlag{
		Name:     "topic",
		Aliases:  []string{"t"},
		EnvVars:  []string{"TOPIC"},
		Required: true,
	}

	TopicsFlag = cli.StringSliceFlag{
		Name:     "topic",
		Aliases:  []string{"t"},
		EnvVars:  []string{"TOPIC"},
		Required: true,
	}

	HeadersFlag = cli.StringSliceFlag{
		Name:        "header",
		Aliases:     []string{"H"},
		EnvVars:     []string{"HEADER"},
		DefaultText: "header=value",
	}
)
