package produce

import (
	"io"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/bit0rez/kafka-tools/flags"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var Produce = cli.Command{
	Name:        "produce",
	Usage:       "Simple message producer from STDIN",
	Description: "Simple message producer from STDIN",
	UsageText:   `$ echo '{"id": 1, "name": "user"}' | kafka-tools produce -t test -H key=value`,
	Action:      produce,
	Flags: []cli.Flag{
		&flags.BootstrapFlag,
		&flags.TopicFlag,
		&flags.HeadersFlag,
	},
}

func produce(ctx *cli.Context) error {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "read message from STDIN")
	}

	cfg := sarama.NewConfig()
	cfg.ClientID = "kafka-tools"
	cfg.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(
		ctx.StringSlice(flags.BootstrapFlag.Name),
		cfg,
	)
	if err != nil {
		return errors.Wrap(err, "create sync producer")
	}

	headersArgs := ctx.StringSlice(flags.HeadersFlag.Name)

	var headers []sarama.RecordHeader
	if len(headersArgs) > 0 {
		headers = make([]sarama.RecordHeader, 0, len(headersArgs))
		for i := range headersArgs {
			kv := strings.Split(headersArgs[i], "=")
			if len(kv) != 2 {
				return errors.Wrap(errors.New("wrong header value"), headersArgs[i])
			}

			headers = append(
				headers,
				sarama.RecordHeader{
					Key:   sarama.ByteEncoder(kv[0]),
					Value: sarama.ByteEncoder(kv[1]),
				},
			)
		}
	}

	msg := sarama.ProducerMessage{
		Topic:   ctx.String(flags.TopicFlag.Name),
		Headers: headers,
		Value:   sarama.ByteEncoder(b),
	}

	errCh := make(chan error, 1)
	go func() {
		_, _, err := p.SendMessage(&msg)
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err = <-errCh:
		if err != nil {
			return errors.Wrap(err, "send message")
		}
	}

	return p.Close()
}
