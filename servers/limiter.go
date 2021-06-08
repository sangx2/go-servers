package servers

import (
	"sync"
	"sync/atomic"
	"time"
)

const UNLIMITED = -1

type Limiter struct {
	Duration time.Duration

	Limit   int32
	remains *int32

	EnableTimeChan chan time.Time
	resetChan      chan bool
	doneChan       chan bool

	wg sync.WaitGroup
}

func NewLimiter(Duration time.Duration, limit int32) *Limiter {
	if limit < UNLIMITED || limit == 0 {
		return nil
	}

	l := &Limiter{
		Duration: Duration,

		Limit:   limit,
		remains: new(int32),

		EnableTimeChan: make(chan time.Time, 1),
		resetChan:      make(chan bool, 1),
		doneChan:       make(chan bool, 1),
	}
	*l.remains = limit

	return l
}

func (l *Limiter) Start() {
	createDone := make(chan error, 1)
	defer close(createDone)

	l.wg.Add(1)
	go func(l *Limiter) {
		defer l.wg.Done()

		var timer = time.NewTimer(l.Duration)
		timer.Stop()

		createDone <- nil

		for {
			select {
			case enableTime := <-timer.C:
				atomic.StoreInt32(l.remains, l.Limit)
				l.EnableTimeChan <- enableTime
			case <-l.resetChan:
				timer.Reset(l.Duration)
			case <-l.doneChan:
				timer.Stop()
				close(l.EnableTimeChan)

				return
			}
		}
	}(l)
	<-createDone
}

func (l *Limiter) IsLimited() bool {
	var ret = false

	switch atomic.LoadInt32(l.remains) {
	case 0:
		ret = true
		//default: // UNLIMITED, > 0
		//	ret = false
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
	l.doneChan <- true
	close(l.doneChan)
}
