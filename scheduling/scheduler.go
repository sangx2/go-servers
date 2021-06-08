package scheduling

import (
	"sync"
	"time"
)

type Scheduler struct {
	Hour, Min, Sec int
	Weekday        map[time.Weekday]bool

	TimeChan chan time.Time
	doneChan chan bool

	wg sync.WaitGroup
}

func NewScheduler(hour, min, sec int, weekdays []time.Weekday) *Scheduler {
	if weekdays == nil {
		return nil
	}

	if hour < 0 || hour > 23 || min < 0 || min > 59 || sec < 0 || sec > 59 || len(weekdays) == 0 {
		return nil
	}

	s := &Scheduler{
		Hour:    hour,
		Min:     min,
		Sec:     sec,
		Weekday: make(map[time.Weekday]bool),

		TimeChan: make(chan time.Time, 1),
		doneChan: make(chan bool, 1),
	}

	for _, day := range weekdays {
		s.Weekday[day] = true
	}

	return s
}

// Start 정해진 시간에 TimeChan 을 통해 알림
func (s *Scheduler) Start() chan time.Time {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), s.Hour, s.Min, s.Sec, 0, now.Location())

	if !t.After(now) {
		t = time.Date(now.Year(), now.Month(), now.Day()+1, s.Hour, s.Min, s.Sec, 0, now.Location())
	}

	s.wg.Add(1)
	go func(s *Scheduler) {
		defer s.wg.Done()

		ticker := time.NewTicker(t.Sub(now))

		select {
		case t := <-ticker.C:
			ticker.Stop()

			ticker = time.NewTicker(24 * time.Hour)

			s.TimeChan <- t

			for {
				select {
				case t := <-ticker.C:
					if s.Weekday[t.Weekday()] {
						s.TimeChan <- t
					}
				case <-s.doneChan:
					return
				}
			}
		case <-s.doneChan:
			ticker.Stop()
			return
		}
	}(s)

	return s.TimeChan
}

// Stop 알림 중지
func (s *Scheduler) Stop() {
	s.doneChan <- true

	s.wg.Wait()
}
