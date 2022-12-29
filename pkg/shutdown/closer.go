package shutdown

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type CloseFunc func(ctx context.Context) error

type Closer struct {
	mu    sync.Mutex
	funcs []CloseFunc
}

func (c *Closer) Add(f CloseFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		messages = make([]string, 0, len(c.funcs))
		complete = make(chan struct{}, 1)
	)

	go func() {
		for _, f := range c.funcs {
			if err := f(ctx); err != nil {
				messages = append(messages, fmt.Sprintf("err: %v", err))
			}
		}
		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown canceled: %v", ctx.Err())
	}

	if len(messages) > 0 {
		return fmt.Errorf("shutdown finished with error(s): \n%s", strings.Join(messages, "\n"))
	}

	return nil
}
