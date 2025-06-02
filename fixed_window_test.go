package ratelimit

import (
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func Test_fixedWindow_allow(t *testing.T) {
	synctest.Run(func() {
		limit := 2
		duration := time.Second
		key := "my-key"

		fw := newFixedWindow(limit, duration)

		want := true
		for range limit {
			got, _ := fw.allow(key)
			if got != want {
				t.Errorf("allow() = %v, want %v", got, want)
			}
		}

		want = false
		for range limit {
			got, remaining := fw.allow(key)
			if got != want {
				t.Errorf("allow() = %v, want %v", got, want)
			}
			if remaining != duration {
				t.Errorf("allow() = %v, want %v", remaining, duration)
			}
		}

		time.Sleep(3 * time.Second)

		want = true
		for range limit {
			got, _ := fw.allow(key)
			if got != want {
				t.Errorf("allow() = %v, want %v", got, want)
			}
		}

		want = false
		for range limit {
			got, remaining := fw.allow(key)
			if got != want {
				t.Errorf("allow() = %v, want %v", got, want)
			}
			if remaining != duration {
				t.Errorf("allow() = %v, want %v", remaining, duration)
			}
		}
	})
}

func Test_fixedWindow_allow_concurrency(t *testing.T) {
	synctest.Run(func() {
		limit := 200
		duration := time.Minute
		key := "my-key"

		fw := newFixedWindow(limit, duration)
		wg := &sync.WaitGroup{}

		wg.Add(limit)
		want := true
		for range limit {
			go func() {
				got, _ := fw.allow(key)
				if got != want {
					t.Errorf("allow() = %v, want %v", got, want)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		wg.Add(limit)
		want = false
		for range limit {
			go func() {
				got, remaining := fw.allow(key)
				if got != want {
					t.Errorf("allow() = %v, want %v", got, want)
				}
				if remaining != duration {
					t.Errorf("allow() = %v, want %v", remaining, duration)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		time.Sleep(2 * time.Minute)

		wg.Add(limit)
		want = true
		for range limit {
			go func() {
				got, _ := fw.allow(key)
				if got != want {
					t.Errorf("allow() = %v, want %v", got, want)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		want = false
		for range limit {
			go func() {
				got, remaining := fw.allow(key)
				if got != want {
					t.Errorf("allow() = %v, want %v", got, want)
				}
				if remaining != duration {
					t.Errorf("allow() = %v, want %v", remaining, duration)
				}
			}()
		}
	})
}
