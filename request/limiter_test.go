package request

import (
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	var limiter *Limiter

	limiter = NewLimiter(time.Second*1, 0)
	if limiter != nil {
		t.Fatalf("limiter is not nil")
	}

	limiter = NewLimiter(time.Second*1, -2)
	if limiter != nil {
		t.Fatalf("limiter is not nil")
	}
}

func TestLimiter(t *testing.T) {
	var limiter = NewLimiter(time.Second*1, 3)
	if limiter == nil {
		t.Fatalf("limiter is nil")
	}

	limiter.Start()

	for i := 0; i < 5; i++ {
		if limiter.IsLimited() {
			<-limiter.EnableTimeChan
		}

		t.Logf("call func")

		limiter.Decrease()
	}

	limiter.Stop()
}
