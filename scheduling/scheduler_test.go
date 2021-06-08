package scheduling

import (
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	var scheduler *Scheduler

	scheduler = NewScheduler(0, 0, 0, nil)
	if scheduler != nil {
		t.Fatalf("scheduler is not nil")
	}

	scheduler = NewScheduler(-1, 0, 0,
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday})
	if scheduler != nil {
		t.Fatalf("scheduler is not nil")
	}

	scheduler = NewScheduler(24, 0, 0,
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday})
	if scheduler != nil {
		t.Fatalf("scheduler is not nil")
	}

	scheduler = NewScheduler(0, -1, 0,
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday})
	if scheduler != nil {
		t.Fatalf("scheduler is not nil")
	}

	scheduler = NewScheduler(0, 60, 0,
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday})
	if scheduler != nil {
		t.Fatalf("scheduler is not nil")
	}

	scheduler = NewScheduler(0, 0, -1,
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday})
	if scheduler != nil {
		t.Fatalf("scheduler is not nil")
	}

	scheduler = NewScheduler(0, 0, 60,
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday})
	if scheduler != nil {
		t.Fatalf("scheduler is not nil")
	}

	scheduler = NewScheduler(0, 0, 60, []time.Weekday{})
	if scheduler != nil {
		t.Fatalf("scheduler is not nil")
	}
}

func TestScheduler(t *testing.T) {
	var scheduler *Scheduler

	n := time.Now()
	scheduler = NewScheduler(n.Hour(), n.Minute(), n.Second()+1,
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday})
	if scheduler == nil {
		t.Fatalf("scheduler is nil")
	}

	<-scheduler.Start()

	t.Logf("call func")

	scheduler.Stop()
}
