package commands

import (
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/bit0rez/kafka-tools/flags"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	flagGroupID = cli.StringFlag{
		Name:    "group",
		Aliases: []string{"g"},
		EnvVars: []string{"CONSUME_GROUP"},
		Value:   "kafka-tools",
	}

	flagStatOnly = cli.BoolFlag{
		Name:    "stat",
		Usage:   "print only stat",
		Aliases: []string{"s"},
		Value:   false,
	}

	Consume = cli.Command{
		Name:        "consume",
		Usage:       "Simple message consumer",
		Description: "Simple message consumer",
		UsageText:   `$ kafka-tools consume -t test | jq `,
		Action:      consume,
		Flags: []cli.Flag{
			&flags.BootstrapFlag,
			&flags.TopicsFlag,
			&flagGroupID,
			&flagStatOnly,
		},
	}
)

func consume(ctx *cli.Context) error {
	cfg := sarama.NewConfig()

	g, err := sarama.NewConsumerGroup(
		ctx.StringSlice(flags.BootstrapFlag.Name),
		ctx.String(flagGroupID.Name),
		cfg,
	)
	if err != nil {
		return errors.Wrap(err, "create consumer group")
	}

	errCh := make(chan error, 1)
	dataCh := make(chan *sarama.ConsumerMessage, 1)

	// consumer
	go func() {
		defer close(errCh)
		defer close(dataCh)
		h := handler{dataCh: dataCh}

		for {
			select {
			case <-ctx.Done():
				if err := g.Close(); err != nil {
					errCh <- err
				}
				return
			default:
			}

			if err = g.Consume(ctx.Context, ctx.StringSlice(flags.TopicsFlag.Name), h); err != nil {
				errCh <- errors.Wrap(err, "consume")
			}
		}
	}()

	// App handler
	statOnly := ctx.Bool(flagStatOnly.Name)
	for {
		select {
		case m, ok := <-dataCh:
			if !ok {
				return nil
			}
			if statOnly {
				fmt.Printf("Topic: %s\tPartition: %d\tOffset: %d\n", m.Topic, m.Partition, m.Offset)
			} else {
				os.Stdout.Write(m.Value)
			}
		case err, ok := <-g.Errors():
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		case err, ok := <-errCh:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		}
	}
}

type handler struct {
	dataCh chan *sarama.ConsumerMessage
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (handler) Setup(s sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (handler) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (h handler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-s.Context().Done():
			return nil
		case m, ok := <-c.Messages():
			if !ok {
				return nil
			}
			h.dataCh <- m
		}
	}
}
