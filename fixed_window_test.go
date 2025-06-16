package ratelimit

import (
	"testing"
	"testing/synctest"
	"time"
)

func Test_FixedWindow(t *testing.T) {
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

func Test_FixedWindow_Concurrency(t *testing.T) {
	synctest.Run(func() {
		limit := 200
		duration := time.Minute
		key := "my-key"

		fw := newFixedWindow(limit, duration)

		want := true
		for range limit {
			go func() {
				got, _ := fw.allow(key)
				if got != want {
					t.Errorf("allow() = %v, want %v", got, want)
				}
			}()
		}

		synctest.Wait()

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

		synctest.Wait()
		time.Sleep(2 * time.Minute)

		want = true
		for range limit {
			go func() {
				got, _ := fw.allow(key)
				if got != want {
					t.Errorf("allow() = %v, want %v", got, want)
				}
			}()
		}

		synctest.Wait()

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
