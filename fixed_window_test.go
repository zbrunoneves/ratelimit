package ratelimit

import (
	"testing"
	"testing/synctest"
	"time"
)

func Test_FixedWindow(t *testing.T) {
	synctest.Run(func() {
		limit := 52
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
		limit := 203
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

func Test_FixedWindow_MultipleKeys(t *testing.T) {
	synctest.Run(func() {
		limit := 52
		duration := time.Second
		keys := []string{"key1", "key2", "key3"}

		fw := newFixedWindow(limit, duration)

		for range limit {
			for i := 0; i < len(keys); i++ {
				ok, _ := fw.allow(keys[i])
				if !ok {
					t.Error("allow() = false, want true")
				}
			}
		}

		for i := 0; i < len(keys); i++ {
			ok, _ := fw.allow(keys[i])
			if ok {
				t.Error("allow() = true, want false")
			}
		}

		time.Sleep(2 * time.Second)

		for i := 0; i < len(keys); i++ {
			ok, _ := fw.allow(keys[i])
			if !ok {
				t.Error("allow() = false, want true")
			}
		}
	})
}

func Test_FixedWindow_TimeRemaining(t *testing.T) {
	synctest.Run(func() {
		limit := 2
		duration := time.Hour
		key := "my-key"

		fw := newFixedWindow(limit, duration)

		for range limit {
			ok, _ := fw.allow(key)
			if !ok {
				t.Error("allow() = false, want true")
			}
		}

		time.Sleep(59 * time.Minute)

		want := time.Minute

		ok, remain := fw.allow(key)
		if ok {
			t.Error("allow() = true, want false")
		}

		if remain != want {
			t.Errorf("allow() time remaining = %v, want %v", remain, want)
		}
	})
}
