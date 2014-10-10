package runner

import (
	"sync"
	"time"
)

type Profiler struct {
	count    int64
	max      time.Duration
	min      time.Duration
	total    time.Duration
	start    time.Time
	duration time.Duration
	avg      time.Duration
	rate     float64
	sync.Mutex
}

func (p *Profiler) Start() {
	p.start = time.Now()
}

func (p *Profiler) Profile(f func()) {
	s := time.Now()
	f()
	p.Add(time.Now().Sub(s))
}

func (p *Profiler) Add(d time.Duration) {
	p.Lock()
	defer p.Unlock()

	p.count++
	p.total += d

	if p.max < d {
		p.max = d
	}

	if p.min > d || p.min == 0 {
		p.min = d
	}
}

func (p *Profiler) Stop() {
	p.avg = time.Duration(int64(p.total) / p.count)
	p.duration = time.Now().Sub(p.start)
	p.rate = float64(p.duration.Seconds()) / float64(p.count)
}
