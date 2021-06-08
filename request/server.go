package request

import (
	"fmt"
	"sync"
)

type Server struct {
	queueSize int

	requestChanMap map[string]chan interface{}
	functionMap    map[string]func(interface{})
	LimitersMap    map[string][]*Limiter

	doneChans []chan bool

	wg sync.WaitGroup
}

func NewServer(qSize int) *Server {
	if qSize <= 0 {
		return nil
	}

	return &Server{
		queueSize: qSize,

		requestChanMap: make(map[string]chan interface{}),

		functionMap: make(map[string]func(interface{})),

		LimitersMap: make(map[string][]*Limiter),
	}
}

func (r *Server) AddLimitersWithFunc(title string, limiters []*Limiter, f func(interface{})) error {
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

func (r *Server) Start() {
	createDoneChan := make(chan error, 1)
	defer close(createDoneChan)

	for title, limiters := range r.LimitersMap {
		for _, limiter := range limiters {
			limiter.Start()
		}

		doneChan := make(chan bool, 1)
		r.doneChans = append(r.doneChans, doneChan)

		r.wg.Add(1)
		go func(requestChan chan interface{}, f func(interface{}), limiters []*Limiter, doneChan chan bool, wg *sync.WaitGroup) {
			defer wg.Done()
			createDoneChan <- nil

			for {
				select {
				case request := <-requestChan:
					for _, limiter := range limiters {
						if limiter.IsLimited() {
							select {
							case <-limiter.EnableTimeChan:
							case <-doneChan:
								return
							}
						} else {
							for len(limiter.EnableTimeChan) != 0 {
								<-limiter.EnableTimeChan
							}
						}
					}

					f(request)

					for _, limiter := range limiters {
						limiter.Decrease()
					}
				case <-doneChan:
					return
				}
			}
		}(r.requestChanMap[title], r.functionMap[title], limiters, doneChan, &r.wg)
		<-createDoneChan
	}
}

func (r *Server) Request(title string, request interface{}) error {
	if len(r.requestChanMap[title]) == r.queueSize {
		return fmt.Errorf("queue is full")
	}

	r.requestChanMap[title] <- request

	return nil
}

func (r *Server) Shutdown() {
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
