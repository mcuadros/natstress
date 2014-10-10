package runner

import (
	"fmt"
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/apcera/nats"
	"github.com/mcuadros/pb"
)

type Client struct {
	cid              int
	conn             *nats.Conn
	subjects         []string
	requests         int
	warmupDuration   time.Duration
	shutdownDuration time.Duration
	rate             int
	progressBar      *pb.ProgressBar
	received         map[string]int32
	delivered        map[string]int32
	sync.Mutex
}

func (c *Client) Run() {
	c.subscribe()
	time.Sleep(c.warmupDuration)
	c.publish()
	time.Sleep(c.shutdownDuration)
}

func (c *Client) subscribe() {
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

func (c *Client) publish() {
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
		err := c.conn.Publish(subject, []byte(uuid.New()))
		if err != nil {
			panic(err)
		}
	}
}

func (c *Client) String() string {
	var received int32
	for _, r := range c.received {
		received += r
	}

	var delivered int32
	for _, d := range c.delivered {
		delivered += d
	}

	return fmt.Sprintf(
		"cid: %d, Received: %d, Delivered: %d",
		c.cid, received, delivered,
	)
}
