package server

import (
	"fmt"
	"testing"
	"time"
)

func printRequestTest(request interface{}) {
	fmt.Println("request test")
}

func TestNewRequestServer(t *testing.T) {
	var requestServer *RequestServer

	requestServer = NewRequestServer(0)
	if requestServer != nil {
		t.Fatalf("requestServer is not nil")
	}

	requestServer = NewRequestServer(-1)
	if requestServer != nil {
		t.Fatalf("requestServer is not nil")
	}
}

func TestRequestServer_AddLimitersWithFunc(t *testing.T) {
	requestServer := NewRequestServer(1000)
	if requestServer == nil {
		t.Fatalf("requestServer is nil")
	}

	secLimiter := NewLimiter(time.Second, 8)
	if secLimiter == nil {
		t.Fatalf("secLimiter is nil")
	}

	minLimiter := NewLimiter(time.Minute, 200)
	if minLimiter == nil {
		t.Fatalf("minLimiter is nil")
	}

	if e := requestServer.AddLimitersWithFunc("request", []*Limiter{secLimiter, minLimiter}, nil); e == nil {
		t.Fatalf("func is nil, thus this errer must not be nil")
	}

	if e := requestServer.AddLimitersWithFunc("request", []*Limiter{secLimiter, minLimiter}, printRequestTest); e != nil {
		t.Fatalf(e.Error())
	}

	if e := requestServer.AddLimitersWithFunc("request", []*Limiter{secLimiter, minLimiter}, printRequestTest); e == nil {
		t.Fatalf("title is exist, thus this errer must not be nil")
	}
}

func TestRequestServer(t *testing.T) {
	requestServer := NewRequestServer(1000)
	if requestServer == nil {
		t.Fatalf("requestServer is nil")
	}

	secLimiter := NewLimiter(time.Second, 1)
	if secLimiter == nil {
		t.Fatalf("secLimiter is nil")
	}

	minLimiter := NewLimiter(time.Second*5, 3)
	if minLimiter == nil {
		t.Fatalf("minLimiter is nil")
	}

	if e := requestServer.AddLimitersWithFunc("request", []*Limiter{secLimiter, minLimiter}, printRequestTest); e != nil {
		t.Fatalf(e.Error())
	}

	requestServer.Start()

	for i := 0; i < 6; i++ {
		if e := requestServer.Request("request", nil); e != nil {
			t.Fatalf(e.Error())
		}
	}

	time.Sleep(time.Second * 10)

	requestServer.Shutdown()
}
