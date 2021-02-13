package timeout

import (
	"context"
	"errors"
)

func ErrIfTimeout(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.New("operation timed out")
	default:
		return nil
	}
}
