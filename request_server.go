// go-request
package server

import (
	"fmt"
	"sync"
)

type RequestServer struct {
	queueSize int

	requestChanMap map[string]chan interface{}
	functionMap    map[string]func()
	LimitersMap    map[string][]*Limiter

	doneChans []chan bool

	wg sync.WaitGroup
}

func NewRequestServer(qSize int) *RequestServer {
	if qSize <= 0 {
		return nil
	}

	return &RequestServer{
		queueSize: qSize,

		requestChanMap: make(map[string]chan interface{}),

		functionMap: make(map[string]func()),

		LimitersMap: make(map[string][]*Limiter),
	}
}

func (r *RequestServer) AddLimitersWithFunc(title string, limiters []*Limiter, f func()) error {
	if _, isExist := r.LimitersMap[title]; isExist {
		return fmt.Errorf("%s is already exist", title)
	}

	if f == nil {
		return fmt.Errorf("func is nil")
	}

	r.LimitersMap[title] = limiters
	r.functionMap[title] = f
	r.requestChanMap[title] = make(chan interface{}, r.queueSize)

	return nil
}

func (r *RequestServer) Start() {
	createDoneChan := make(chan error, 1)
	defer close(createDoneChan)

	doneChan := make(chan bool, 1)
	r.doneChans = append(r.doneChans, doneChan)

	for title, limiters := range r.LimitersMap {
		for _, limiter := range limiters {
			limiter.Start()
		}

		r.wg.Add(1)
		go func(requestChan chan interface{}, f func(), limiters []*Limiter, doneChan chan bool) {
			defer r.wg.Done()
			createDoneChan <- nil

			for {
				select {
				case <-requestChan:
					for _, limiter := range limiters {
						if limiter.IsLimited() {
							<-limiter.EnableTimeChan
						}
					}

					f()

					for _, limiter := range limiters {
						limiter.Decrease()
					}
				case <-doneChan:
					return
				}
			}
		}(r.requestChanMap[title], r.functionMap[title], limiters, doneChan)
		<-createDoneChan
	}
}

func (r *RequestServer) Request(title string) error {
	if len(r.requestChanMap[title]) == r.queueSize {
		return fmt.Errorf("queue is full")
	}

	r.requestChanMap[title] <- nil

	return nil
}

func (r *RequestServer) Shutdown() {
	for _, limiters := range r.LimitersMap {
		for _, limiter := range limiters {
			limiter.Stop()
		}
	}

	for _, doneChan := range r.doneChans {
		doneChan <- true
	}

	r.wg.Wait()
}
