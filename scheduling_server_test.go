package server

import (
	"fmt"
	"testing"
	"time"
)

func printSchedulingTest() {
	fmt.Println("scheduling test")
}

func TestSchedulingServer_AddSchedulerWithFunc(t *testing.T) {
	schedulingServer := NewSchedulingServer()
	if schedulingServer == nil {
		t.Fatalf("schedulingServer is nil")
	}

	n := time.Now()
	scheduler := NewScheduler(n.Hour(), n.Minute(), n.Second()+1,
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday})

	if e := schedulingServer.AddSchedulerWithFunc("scheduling", nil, printSchedulingTest); e == nil {
		t.Fatalf("scheduler is nil, thus this eror must not be nil")
	}

	if e := schedulingServer.AddSchedulerWithFunc("scheduling", scheduler, nil); e == nil {
		t.Fatalf("function is nil, thus this errer must not be nil")
	}

	if e := schedulingServer.AddSchedulerWithFunc("scheduling", scheduler, printSchedulingTest); e != nil {
		t.Fatalf(e.Error())
	}

	if e := schedulingServer.AddSchedulerWithFunc("scheduling", scheduler, printSchedulingTest); e == nil {
		t.Fatalf("title is exist, thus this errer must not be nil")
	}
}

func TestSchedulingServer(t *testing.T) {
	schedulingServer := NewSchedulingServer()
	if schedulingServer == nil {
		t.Fatalf("schedulingServer is nil")
	}

	n := time.Now()
	scheduler := NewScheduler(n.Hour(), n.Minute(), n.Second()+1,
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday})

	if e := schedulingServer.AddSchedulerWithFunc("printTest", scheduler, printSchedulingTest); e != nil {
		t.Fatalf(e.Error())
	}

	schedulingServer.Start()

	time.Sleep(time.Second * 3)

	schedulingServer.Shutdown()
}
