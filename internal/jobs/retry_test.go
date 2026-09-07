package jobs

import (
	"testing"
	"time"
)

func TestRetryDelay(t *testing.T) {
	got := RetryDelay(1)
	want := 2 * time.Second

	if got != want {
		t.Errorf("RetryDelay(1) = %v, want %v", got, want)
	}
}
