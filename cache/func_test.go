package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/carterjs/webtools/cache"
)

func TestFunc(t *testing.T) {
	var currentTime = time.Now()
	cache.Now = func() (t time.Time) {
		return currentTime
	}

	var count int
	var fail bool
	get := cache.Func(time.Hour, func() (int, error) {
		count++
		if fail {
			return count, errors.New("some error")
		}

		return count, nil
	})

	t.Run("error when failing on initial call", func(t *testing.T) {
		fail = true
		_, err := get()
		if err == nil {
			t.Fatalf("expected error when failing on initial call")
		}

		fail = false
	})

	t.Run("cached calls", func(t *testing.T) {
		// function call should be cached no matter how many times we call
		if v, _ := get(); v != 2 {
			t.Fatalf("expected >1 calls to be cached, but function was invoked %d times", count)
		}
	})

	t.Run("calls after expiration", func(t *testing.T) {
		currentTime = currentTime.Add(time.Hour + time.Second)

		// after caching period, the function should get invoked again just once
		if v, _ := get(); v != 3 {
			t.Fatalf("expected call after expiration to increase count")
		}
	})

	t.Run("return stale data when getter fails after initial success", func(t *testing.T) {
		fail = true

		count++

		currentTime = currentTime.Add(time.Hour + time.Second)

		v, err := get()
		if err != nil {
			t.Fatalf("unexpected error when returning stale data: %v", err)
		}

		if v != 3 {
			t.Fatal("expected true, found false")
		}
	})
}
