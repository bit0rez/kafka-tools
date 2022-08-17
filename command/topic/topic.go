package topic

import (
	"fmt"
	"sort"

	"github.com/Shopify/sarama"
	"github.com/bit0rez/kafka-tools/command"
	"github.com/bit0rez/kafka-tools/flags"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	flagPartitions = cli.UintFlag{
		Name:    "partitions",
		Usage:   "Count of topic partitions",
		Aliases: []string{"p"},
		Value:   1,
	}

	flagReplicationFactor = cli.UintFlag{
		Name:    "replication_factor",
		Usage:   "Topic replication factor",
		Aliases: []string{"r"},
		Value:   1,
	}

	Topic = cli.Command{
		Name:  "topic",
		Usage: "Manage topics",
		Subcommands: []*cli.Command{
			{
				Name:   "create",
				Action: command.HandleCliCtx(topicCreate),
				Flags: []cli.Flag{
					&flags.BootstrapFlag,
					&flags.TopicFlag,
					&flagPartitions,
					&flagReplicationFactor,
				},
			},
			{
				Name:   "list",
				Action: command.HandleCliCtx(topicList),
				Flags: []cli.Flag{
					&flags.BootstrapFlag,
				},
			},
			{
				Name:   "delete",
				Action: command.HandleCliCtx(topicDelete),
				Flags: []cli.Flag{
					&flags.BootstrapFlag,
					&flags.TopicFlag,
				},
			},
		},
	}
)

type topic struct {
	Name              string
	NumPartitions     int32
	ReplicationFactor int16
}

type topics []topic

func (t topics) Print() {
	maxLen := 0
	for i := range t {
		l := len(t[i].Name)
		if l > maxLen {
			maxLen = l
		}
	}

	for i := range t {
		fill := maxLen - len(t[i].Name)
		fmt.Print("Topic: ", t[i].Name, "  ")
		for j := 0; j < fill; j++ {
			fmt.Print(" ")
		}
		fmt.Print("Partitions: ", t[i].NumPartitions, "\t")
		fmt.Print("Replication factor: ", t[i].ReplicationFactor, "\n")
	}
}

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

func topicList(ctx *cli.Context) error {
	cfg := sarama.NewConfig()
	adm, err := sarama.NewClusterAdmin(ctx.StringSlice(flags.BootstrapFlag.Name), cfg)
	if err != nil {
		return errors.Wrap(err, "create cluster admin")
	}

	l, err := adm.ListTopics()
	if err != nil {
		return errors.Wrap(err, "list topics")
	}

	t := make(topics, len(l))
	var i int
	for k, details := range l {
		t[i].Name = k
		t[i].NumPartitions = details.NumPartitions
		t[i].ReplicationFactor = details.ReplicationFactor
		i++
	}

	sort.Slice(t, func(i, j int) bool {
		return t[i].Name < t[j].Name
	})

	t.Print()

	return nil
}
