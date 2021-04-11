package server

import (
	"fmt"
	"sync"
)

type SchedulingServer struct {
	SchedulerMap map[string]*Scheduler
	functionMap  map[string]func()

	doneChans []chan bool

	wg sync.WaitGroup
}

func NewSchedulingServer() *SchedulingServer {
	return &SchedulingServer{
		SchedulerMap: make(map[string]*Scheduler),

		functionMap: make(map[string]func()),
	}
}

func (s *SchedulingServer) AddSchedulerWithFunc(title string, scheduler *Scheduler, function func()) error {
	if _, isExist := s.SchedulerMap[title]; isExist {
		return fmt.Errorf("%s is already exist", title)
	}

	if scheduler == nil {
		return fmt.Errorf("scheduler is nil")
	}

	if function == nil {
		return fmt.Errorf("function is nil")
	}

	s.SchedulerMap[title] = scheduler
	s.functionMap[title] = function

	return nil
}

func (s *SchedulingServer) Start() {
	for title, scheduler := range s.SchedulerMap {
		doneChan := make(chan bool, 1)
		s.doneChans = append(s.doneChans, doneChan)

		s.wg.Add(1)
		go func(scheduler *Scheduler, function func(), doneChan chan bool, wg *sync.WaitGroup) {
			defer wg.Done()

			timeChan := scheduler.Start()
			for {
				select {
				case <-timeChan:
					function()
				case <-doneChan:
					return
				}
			}
		}(scheduler, s.functionMap[title], doneChan, &s.wg)
	}
}

func (s *SchedulingServer) Shutdown() {
	for _, scheduler := range s.SchedulerMap {
		scheduler.Stop()
	}

	for _, doneChan := range s.doneChans {
		doneChan <- true
	}

	s.wg.Wait()
}
