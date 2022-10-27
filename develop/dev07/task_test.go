package dev07

import (
	"fmt"
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	t.Run("Test-0", func(t *testing.T) {
		fmt.Println("geg")
		sig := func(after time.Duration) <-chan interface{} {
			c := make(chan interface{})
			go func() {
				defer close(c)
				time.Sleep(after)
			}()
			return c
		}

		start := time.Now()
		<-or(
			sig(2*time.Hour),
			sig(5*time.Minute),
			sig(6*time.Second),
			sig(1*time.Second),
			sig(54*time.Second),
		)
		fmt.Printf("fone after %v", time.Since(start))
	})
}
