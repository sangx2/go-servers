package ticker

import (
	"fmt"
	"sync"
	"time"
)

type Server struct {
	TickerMap map[string]*time.Ticker
	cbFuncMap map[string]func()

	doneChans []chan bool

	wg sync.WaitGroup
}

func NewServer() *Server {
	return &Server{
		TickerMap: make(map[string]*time.Ticker),

		cbFuncMap: make(map[string]func()),
	}
}

func (t *Server) AddTickerWithFunc(title string, duration time.Duration, cbFunc func()) error {
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

func (t *Server) ResetTicker(title string, duration time.Duration) error {
	ticker, isExist := t.TickerMap[title]
	if !isExist {
		return fmt.Errorf("title is not exist")
	}
	ticker.Reset(duration)

	return nil
}

func (t *Server) Start() {
	for title, ticker := range t.TickerMap {
		doneChan := make(chan bool, 1)
		t.doneChans = append(t.doneChans, doneChan)

		t.wg.Add(1)
		go func(ticker *time.Ticker, cbFunc func(), doneChan chan bool, wg *sync.WaitGroup) {
			defer wg.Done()
			defer ticker.Stop()

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

func (t *Server) Shutdown() {
	for _, doneChan := range t.doneChans {
		doneChan <- true
	}

	t.wg.Wait()
}
