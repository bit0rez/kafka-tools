package commands

import (
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/bit0rez/kafka-tools/flags"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const (
	OUTPUT_STAT = "stat"
	OUTPUT_RAW  = "raw"
	OUTPUT_HEX  = "hex"
)

var (
	flagGroupID = cli.StringFlag{
		Name:    "group",
		Aliases: []string{"g"},
		EnvVars: []string{"CONSUME_GROUP"},
		Value:   "kafka-tools",
	}

	flagFormat = cli.StringFlag{
		Name:    "format",
		Usage:   "one of: stat, raw, hex",
		Aliases: []string{"f"},
		Value:   OUTPUT_STAT,
	}

	flagOffset = cli.Int64Flag{
		Name:    "offset",
		Aliases: []string{"o"},
		Usage:   "For newest: -1. For oldest: -2.",
		Value:   sarama.OffsetNewest,
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
			&flagFormat,
			&flagOffset,
		},
	}
)

func consume(ctx *cli.Context) error {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = ctx.Int64(flagOffset.Name)

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
	for {
		select {
		case m, ok := <-dataCh:
			if !ok {
				return nil
			}
			out(ctx.String(flagFormat.Name), m)
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

func out(format string, msg *sarama.ConsumerMessage) {
	switch format {
	case OUTPUT_STAT:
		fmt.Printf("Topic: %s\tPartition: %d\tOffset: %d\n", msg.Topic, msg.Partition, msg.Offset)
	case OUTPUT_HEX:
		fmt.Printf("%x\n", msg.Value)
	case OUTPUT_RAW:
		os.Stdout.Write(msg.Value)
	default:
	}
}
