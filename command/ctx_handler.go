package command

import "github.com/urfave/cli/v2"

// HandleCliCtx wraps cli.ActionFunc because of sarama.Client has no context on connect and can stuck.
func HandleCliCtx(a cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		errCh := make(chan error, 1)

		go func() {
			err := a(ctx)
			errCh <- err
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		}
	}
}
