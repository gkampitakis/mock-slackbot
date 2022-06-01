package utils

import (
	"log"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestConcurrent(t *testing.T) {
	t.Run("should panic", func(t *testing.T) {
		defer func() {
			e := recover()
			if e != "Deadlock can't process items" {
				t.Errorf("Concurrent panicked with %s\n", e)
			}
		}()
		Concurrent([]string{}, func(item string) {}, -1)

		t.Error("Concurrent did not panic")
	})

	t.Run("should run concurrently", func(t *testing.T) {
		start := time.Now()
		items := make([]int, 20)

		Concurrent(items, func(item int) {
			time.Sleep(time.Second)
		}, 10)

		timePassed := time.Since(start)
		// we do 2 seconds and 100ms to allow some margin of error
		if timePassed > time.Millisecond*2100 {
			t.Errorf("Looks like function didn't run concurrently")
		}

		// this tries to test the max concurrency
		if timePassed < time.Millisecond*1900 {
			t.Errorf("Looks like everything run at the same time")
		}
	})
}

func TestEscapeTags(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{
			input:    "hello <!here> world",
			expected: "hello  world",
			name:     "here tag in between",
		},
		{
			input:    "<!channel> hello world",
			expected: " hello world",
			name:     "channel tag at start",
		},
		{
			input:    "<@4324234UDAsdd> hello world",
			expected: " hello world",
			name:     "username tag-ignore cases",
		},
		{
			input:    "<@!blue> hello world",
			expected: "<@!blue> hello world",
			name:     "should ignore invalid tags",
		},
		{
			input:    "<1231321> hello world",
			expected: "<1231321> hello world",
			name:     "should ignore invalid tags",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := EscapeSlackTags(test.input)

			if test.expected != got {
				t.Errorf("expected %s but got %s\n", test.expected, got)
			}
		})
	}
}

func TestConfiguration(t *testing.T) {
	t.Run("should init config", func(t *testing.T) {
		t.Setenv("SLACK_SIGNING_SECRET", "mock-secret")
		t.Setenv("SLACK_BOT_TOKEN", "mock-token")
		t.Setenv("BOT_MODE", "mock")

		snaps.MatchSnapshot(t, NewConfiguration())
	})

	t.Run("should set production 'true'", func(t *testing.T) {
		t.Setenv("SLACK_SIGNING_SECRET", "mock-secret")
		t.Setenv("SLACK_BOT_TOKEN", "mock-token")
		t.Setenv("BOT_MODE", "production")

		snaps.MatchSnapshot(t, NewConfiguration())
	})

	t.Run("should log fatalln", func(t *testing.T) {
		t.Cleanup(func() {
			logFatalln = log.Fatalln
		})
		called := false

		logFatalln = func(v ...any) {
			called = true
		}

		NewConfiguration()
		if !called {
			t.Error("fatal was not called")
		}
	})
}
