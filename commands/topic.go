package commands

import (
	"github.com/Shopify/sarama"
	"github.com/bit0rez/kafka-tools/flags"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	flagPartitions = cli.UintFlag{
		Name:    "partitions",
		Aliases: []string{"p"},
		Value:   1,
	}

	flagReplicationFactor = cli.UintFlag{
		Name:    "replication_factor",
		Aliases: []string{"r"},
		Value:   1,
	}

	Topic = cli.Command{
		Name: "topic",
		Subcommands: []*cli.Command{
			{
				Name:   "create",
				Action: topicCreate,
				Flags: []cli.Flag{
					&flags.BootstrapFlag,
					&flags.TopicFlag,
					&flagPartitions,
					&flagReplicationFactor,
				},
			},
			{
				Name:   "delete",
				Action: topicDelete,
				Flags: []cli.Flag{
					&flags.BootstrapFlag,
					&flags.TopicFlag,
				},
			},
		},
	}
)

func topicCreate(ctx *cli.Context) error {
	cfg := sarama.NewConfig()
	adm, err := sarama.NewClusterAdmin(ctx.StringSlice(flags.BootstrapFlag.Name), cfg)
	if err != nil {
		return errors.Wrap(err, "create cluster admin")
	}

	td := sarama.TopicDetail{
		NumPartitions:     int32(ctx.Uint(flagPartitions.Name)),
		ReplicationFactor: int16(ctx.Uint(flagReplicationFactor.Name)),
	}

	return adm.CreateTopic(ctx.String(flags.TopicFlag.Name), &td, false)
}

func topicDelete(ctx *cli.Context) error {
	cfg := sarama.NewConfig()
	adm, err := sarama.NewClusterAdmin(ctx.StringSlice(flags.BootstrapFlag.Name), cfg)
	if err != nil {
		return errors.Wrap(err, "create cluster admin")
	}

	return adm.DeleteTopic(ctx.String(flags.TopicFlag.Name))
}
