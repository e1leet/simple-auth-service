package postgresql

import (
	"strings"
	"time"
)

func FormatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err := fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}

		return
	}

	return
}
