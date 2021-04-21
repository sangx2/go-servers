package server

import (
	"fmt"
	"testing"
	"time"
)

func PrintTickerTest() {
	fmt.Println("ticker test")
}

func TestTickerServer_AddTickerWithFunc(t *testing.T) {
	tickerServer := NewTickerServer()
	if tickerServer == nil {
		t.Fatalf("tickerServer is nil")
	}

	if e := tickerServer.AddTickerWithFunc("ticker", 0, PrintTickerTest); e == nil {
		t.Fatalf("duration is 0, thus this error must not be nil")
	}

	if e := tickerServer.AddTickerWithFunc("ticker", time.Second, nil); e == nil {
		t.Fatalf("function is nil, thus this errer must not be nil")
	}

	if e := tickerServer.AddTickerWithFunc("ticker", time.Second, PrintTickerTest); e != nil {
		t.Fatalf(e.Error())
	}

	if e := tickerServer.AddTickerWithFunc("ticker", time.Second, PrintTickerTest); e == nil {
		t.Fatalf("title is exist, thus this errer must not be nil")
	}
}

func TestTickerServer(t *testing.T) {
	tickerServer := NewTickerServer()
	if tickerServer == nil {
		t.Fatalf("tickerServer is nil")
	}

	if e := tickerServer.AddTickerWithFunc("ticker", time.Second, PrintTickerTest); e != nil {
		t.Fatalf(e.Error())
	}

	tickerServer.Start()

	time.Sleep(time.Second * 1)

	if e := tickerServer.ResetTicker("ticker", time.Millisecond*500); e != nil {
		t.Fatalf(e.Error())
	}

	time.Sleep(time.Second * 2)

	tickerServer.Shutdown()
}
