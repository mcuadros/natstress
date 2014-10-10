package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/apcera/nats"
)

var (
	flagServer   = flag.String("h", nats.DefaultURL, "")
	flagSubjects = flag.Int("s", 5, "")
	flagRequests = flag.Int("r", 2000, "")
	flagClients  = flag.Int("c", 5, "")
)

func main() {
	flag.Parse()

	runner := &Natstress{
		Url:         *flagServer,
		NumClients:  *flagClients,
		NumSubjects: *flagSubjects,
		NumRequests: *flagRequests,
	}

	runner.Run()
}

type Natstress struct {
	Url         string
	NumClients  int
	NumSubjects int
	NumRequests int
	clients     []*Client
	subjects    []string
	sync.WaitGroup
}

func (n *Natstress) Run() {
	n.buildSubjects()
	n.buildClients()

	for _, client := range n.clients {
		n.Add(1)
		go func(c *Client) {
			c.Run()
			n.Done()
		}(client)
	}

	n.Wait()
	n.printResume()
}

func (n *Natstress) buildSubjects() {
	n.subjects = make([]string, n.NumSubjects)
	for i := 0; i < n.NumSubjects; i++ {
		n.subjects[i] = uuid.New()
	}

	n.subjects[0] = "foo"
}

func (n *Natstress) buildClients() {
	n.clients = make([]*Client, n.NumClients)
	for i := 0; i < n.NumClients; i++ {
		n.clients[i] = n.buildClient(i)
	}
}

func (n *Natstress) buildClient(cid int) *Client {
	nc, err := nats.Connect(n.Url)
	if err != nil {
		panic(err)
	}

	return &Client{
		cid:       cid,
		conn:      nc,
		subjects:  n.subjects,
		requests:  n.NumRequests,
		received:  make(map[string]int32),
		delivered: make(map[string]int32),
	}
}

func (n *Natstress) printResume() {
	for _, client := range n.clients {
		fmt.Println(client)
	}
}

type Client struct {
	cid       int
	conn      *nats.Conn
	subjects  []string
	requests  int
	received  map[string]int32
	delivered map[string]int32
	sync.Mutex
}

func (c *Client) Run() {
	c.subscribe()
	c.publish()

	time.Sleep(5 * time.Second)
}

func (c *Client) subscribe() {
	for _, subject := range c.subjects {
		c.received[subject] = 0
		c.conn.Subscribe(subject, func(m *nats.Msg) {
			c.Lock()
			defer c.Unlock()

			c.received[m.Subject]++
			//fmt.Printf("cid %d, Received a message: %s\n", c.cid, string(m.Data))
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
	for _, subject := range c.subjects {
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
