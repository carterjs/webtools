package cache

import (
	"log"
	"time"
)

var Now = time.Now

// Func caches the result of a function for the given duration
func Func[T any](ttl time.Duration, get func() (T, error)) func() (T, error) {
	var value T
	var lastUpdate time.Time

	return func() (T, error) {
		if Now().Sub(lastUpdate) > ttl {
			newValue, err := get()
			if err != nil {
				if lastUpdate.IsZero() {
					return value, err
				} else {
					log.Printf("serving stale data after error: %v", err)
					return value, nil
				}
			}

			value = newValue
			lastUpdate = Now()
		}

		return value, nil
	}
}
