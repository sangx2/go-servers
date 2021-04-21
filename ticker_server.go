package server

import (
	"fmt"
	"sync"
	"time"
)

type TickerServer struct {
	TickerMap map[string]*time.Ticker
	cbFuncMap map[string]func()

	doneChans []chan bool

	wg sync.WaitGroup
}

func NewTickerServer() *TickerServer {
	return &TickerServer{
		TickerMap: make(map[string]*time.Ticker),

		cbFuncMap: make(map[string]func()),
	}
}

func (t *TickerServer) AddTickerWithFunc(title string, duration time.Duration, cbFunc func()) error {
	if _, isExist := t.TickerMap[title]; isExist {
		return fmt.Errorf("%s is already exist", title)
	}

	if duration <= 0 {
		return fmt.Errorf("duration(%d) is invalid", duration)
	}

	if cbFunc == nil {
		return fmt.Errorf("cbFunction is nil")
	}

	t.TickerMap[title] = time.NewTicker(duration)
	t.cbFuncMap[title] = cbFunc

	return nil
}

func (t *TickerServer) ResetTicker(title string, duration time.Duration) error {
	ticker, isExist := t.TickerMap[title]
	if !isExist {
		return fmt.Errorf("title is not exist")
	}
	ticker.Reset(duration)

	return nil
}

func (t *TickerServer) Start() {
	for title, ticker := range t.TickerMap {
		doneChan := make(chan bool, 1)
		t.doneChans = append(t.doneChans, doneChan)

		t.wg.Add(1)
		go func(ticker *time.Ticker, cbFunc func(), doneChan chan bool, wg *sync.WaitGroup) {
			defer wg.Done()

			for {
				select {
				case <-ticker.C:
					cbFunc()
				case <-doneChan:
					return
				}
			}
		}(ticker, t.cbFuncMap[title], doneChan, &t.wg)
	}
}

func (t *TickerServer) Shutdown() {
	for _, ticker := range t.TickerMap {
		ticker.Stop()
	}

	for _, doneChan := range t.doneChans {
		doneChan <- true
	}

	t.wg.Wait()
}
