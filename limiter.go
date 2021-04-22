package server

import (
	"sync"
	"sync/atomic"
	"time"
)

const UNLIMITED = -1

type Limiter struct {
	Duration time.Duration
	timer    *time.Timer

	Limit   int32
	remains *int32

	EnableTimeChan chan time.Time
	resetChan      chan bool
	doneChan       chan bool

	wg sync.WaitGroup
}

func NewLimiter(duration time.Duration, limit int32) *Limiter {
	if duration <= 0 {
		return nil
	}

	if limit < UNLIMITED || limit == 0 {
		return nil
	}

	l := &Limiter{
		Duration: duration,

		Limit:   limit,
		remains: new(int32),

		EnableTimeChan: make(chan time.Time, 1),
		resetChan:      make(chan bool, 1),
		doneChan:       make(chan bool, 1),
	}
	l.timer = time.NewTimer(duration)
	l.timer.Stop()

	*l.remains = limit

	return l
}

func (l *Limiter) Start() {
	createDone := make(chan error, 1)
	defer close(createDone)

	l.wg.Add(1)
	go func(l *Limiter) {
		defer l.wg.Done()
		defer l.timer.Stop()

		createDone <- nil

		for {
			select {
			case enableTime := <-l.timer.C:
				atomic.StoreInt32(l.remains, l.Limit)
				l.EnableTimeChan <- enableTime
			case <-l.resetChan:
				l.timer.Reset(l.Duration)
			case <-l.doneChan:
				return
			}
		}
	}(l)
	<-createDone
}

func (l *Limiter) IsLimited() bool {
	var ret = false

	switch atomic.LoadInt32(l.remains) {
	case l.Limit, UNLIMITED:
		ret = false
	case 0:
		ret = true
	}

	return ret
}

func (l *Limiter) Decrease() {
	switch atomic.LoadInt32(l.remains) {
	case UNLIMITED:
	case l.Limit:
		l.resetChan <- true
		fallthrough
	default: // 0 < remains < limit
		atomic.AddInt32(l.remains, -1)
	}
}

func (l *Limiter) Stop() {
	l.timer.Stop()
	for len(l.EnableTimeChan) != 0 {
		<-l.EnableTimeChan
	}
	l.doneChan <- true

	l.wg.Wait()
}
