package shutdown

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCloser_Add(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		closer := &Closer{}
		closer.Add(func(ctx context.Context) error {
			return nil
		})
		assert.Equal(t, 1, len(closer.funcs))
	})
}

func TestCloser_Close(t *testing.T) {
	t.Run("deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		closer := &Closer{}
		closer.Add(func(ctx context.Context) error {
			time.Sleep(time.Second)
			return nil
		})
		err := closer.Close(ctx)
		assert.ErrorContains(t, err, context.DeadlineExceeded.Error())
	})

	t.Run("close func throws error", func(t *testing.T) {
		const errorText = "empty error"
		closer := &Closer{}
		closer.Add(func(ctx context.Context) error {
			return errors.New(errorText)
		})
		err := closer.Close(context.Background())
		assert.ErrorContains(t, err, errorText)
	})

	t.Run("close", func(t *testing.T) {
		closer := &Closer{}
		closer.Add(func(ctx context.Context) error {
			return nil
		})
		err := closer.Close(context.Background())
		assert.NoError(t, err)
	})
}
