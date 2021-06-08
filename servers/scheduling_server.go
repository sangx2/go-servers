package servers

import (
	"fmt"
	"sync"
)

type SchedulingServer struct {
	SchedulerMap map[string]*Scheduler
	cbFuncMap    map[string]func()

	doneChans []chan bool

	wg sync.WaitGroup
}

func NewSchedulingServer() *SchedulingServer {
	return &SchedulingServer{
		SchedulerMap: make(map[string]*Scheduler),

		cbFuncMap: make(map[string]func()),
	}
}

func (s *SchedulingServer) AddSchedulerWithFunc(title string, scheduler *Scheduler, cbFunc func()) error {
	if _, isExist := s.SchedulerMap[title]; isExist {
		return fmt.Errorf("%s is already exist", title)
	}

	if scheduler == nil {
		return fmt.Errorf("scheduler is nil")
	}

	if cbFunc == nil {
		return fmt.Errorf("cbFunc is nil")
	}

	s.SchedulerMap[title] = scheduler
	s.cbFuncMap[title] = cbFunc

	return nil
}

func (s *SchedulingServer) Start() {
	for title, scheduler := range s.SchedulerMap {
		doneChan := make(chan bool, 1)
		s.doneChans = append(s.doneChans, doneChan)

		s.wg.Add(1)
		go func(scheduler *Scheduler, cbFunc func(), doneChan chan bool, wg *sync.WaitGroup) {
			defer wg.Done()
			defer scheduler.Stop()

			timeChan := scheduler.Start()
			for {
				select {
				case <-timeChan:
					cbFunc()
				case <-doneChan:
					return
				}
			}
		}(scheduler, s.cbFuncMap[title], doneChan, &s.wg)
	}
}

func (s *SchedulingServer) Shutdown() {
	for _, doneChan := range s.doneChans {
		doneChan <- true
	}

	s.wg.Wait()
}
