package ratelimit

import (
	"testing"
	"testing/synctest"
	"time"
)

func Test_SlidingWindow(t *testing.T) {
	synctest.Run(func() {
		limit := 52
		duration := time.Second
		key := "my-key"

		sw := newSlidingWindow(limit, duration)

		want := true
		for range limit {
			got, _ := sw.allow(key)
			if got != want {
				t.Errorf("allow() = %v, want %v", got, want)
			}
		}

		want = false
		for range limit {
			got, remaining := sw.allow(key)
			if got != want {
				t.Errorf("allow() = %v, want %v", got, want)
			}
			if remaining > duration || remaining <= 0 {
				t.Errorf("remaining = %v, out of expected range", remaining)
			}
		}

		time.Sleep(2 * time.Second)

		want = true
		for range limit {
			got, _ := sw.allow(key)
			if got != want {
				t.Errorf("allow() = %v, want %v", got, want)
			}
		}
	})
}

func Test_SlidingWindow_Concurrency(t *testing.T) {
	synctest.Run(func() {
		limit := 203
		duration := time.Minute
		key := "my-key"

		sw := newSlidingWindow(limit, duration)

		want := true
		for range limit {
			go func() {
				got, _ := sw.allow(key)
				if got != want {
					t.Errorf("allow() = %v, want %v", got, want)
				}
			}()
		}

		synctest.Wait()

		want = false
		for range limit {
			go func() {
				got, remaining := sw.allow(key)
				if got != want {
					t.Errorf("allow() = %v, want %v", got, want)
				}
				if remaining > duration || remaining <= 0 {
					t.Errorf("remaining = %v, out of expected range", remaining)
				}
			}()
		}

		synctest.Wait()

		time.Sleep(2 * time.Minute)

		want = true
		for range limit {
			go func() {
				got, _ := sw.allow(key)
				if got != want {
					t.Errorf("allow() = %v, want %v", got, want)
				}
			}()
		}

		synctest.Wait()
	})
}

func Test_SlidingWindow_MultipleKeys(t *testing.T) {
	synctest.Run(func() {
		limit := 52
		duration := time.Second
		keys := []string{"key1", "key2", "key3"}

		sw := newSlidingWindow(limit, duration)

		for range limit {
			for _, key := range keys {
				ok, _ := sw.allow(key)
				if !ok {
					t.Errorf("allow(%s) = false, want true", key)
				}
			}
		}

		for _, key := range keys {
			ok, _ := sw.allow(key)
			if ok {
				t.Errorf("allow(%s) = true, want false", key)
			}
		}

		time.Sleep(2 * time.Second)

		for _, key := range keys {
			ok, _ := sw.allow(key)
			if !ok {
				t.Errorf("allow(%s) = false, want true", key)
			}
		}
	})
}

func Test_SlidingWindow_TimeRemaining(t *testing.T) {
	synctest.Run(func() {
		limit := 2
		duration := time.Minute
		key := "my-key"

		sw := newSlidingWindow(limit, duration)

		for range limit {
			ok, _ := sw.allow(key)
			if !ok {
				t.Error("allow() = false, want true")
			}
		}

		time.Sleep(32 * time.Second)

		ok, remain := sw.allow(key)
		if ok {
			t.Error("allow() = true, want false")
		}

		if remain != 28*time.Second {
			t.Errorf("remaining = %v, want 28s", remain)
		}
	})
}
