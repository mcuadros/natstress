package runner

import (
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/apcera/nats"
	"github.com/mcuadros/pb"
)

type Client struct {
	cid         int
	conn        *nats.Conn
	subjects    []string
	requests    int
	rate        int
	progressBar *pb.ProgressBar
	profiler    *Profiler
	received    map[string]int32
	delivered   map[string]int32
	sync.Mutex
}

func (c *Client) Subscribe() {
	for _, subject := range c.subjects {
		c.received[subject] = 0
		c.conn.Subscribe(subject, func(m *nats.Msg) {
			c.Lock()
			defer c.Unlock()

			c.progressBar.Increment()
			c.received[m.Subject]++
		})
	}
}

func (c *Client) Publish() {
	for _, subject := range c.subjects {
		c.delivered[subject] = 0
	}

	for i := 0; i < c.requests; i++ {
		c.publishToSubjects()
	}
}

func (c *Client) publishToSubjects() {
	var throttle <-chan time.Time
	if c.rate > 0 {
		throttle = time.Tick(time.Duration(1e6/(c.rate)) * time.Microsecond)
	}

	for _, subject := range c.subjects {
		if c.rate > 0 {
			<-throttle
		}

		c.delivered[subject]++
		c.profiler.Profile(func() {
			err := c.conn.Publish(subject, []byte(uuid.New()))
			if err != nil {
				panic(err)
			}
		})
	}
}
