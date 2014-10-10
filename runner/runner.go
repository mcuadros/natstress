package runner

import (
	"fmt"
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/apcera/nats"
)

type Runner struct {
	Url              string
	NumClients       int
	NumSubjects      int
	NumRequests      int
	WarmupDuration   time.Duration
	ShutdownDuration time.Duration
	Rate             int
	clients          []*Client
	subjects         []string
	sync.WaitGroup
}

func (r *Runner) Run() {
	r.buildSubjects()
	r.buildClients()

	for _, client := range r.clients {
		r.Add(1)
		go func(c *Client) {
			c.Run()
			r.Done()
		}(client)
	}

	r.Wait()
	r.printResume()
}

func (r *Runner) buildSubjects() {
	r.subjects = make([]string, r.NumSubjects)
	for i := 0; i < r.NumSubjects; i++ {
		r.subjects[i] = uuid.New()
	}

	r.subjects[0] = "foo"
}

func (r *Runner) buildClients() {
	r.clients = make([]*Client, r.NumClients)
	for i := 0; i < r.NumClients; i++ {
		r.clients[i] = r.buildClient(i)
	}
}

func (r *Runner) buildClient(cid int) *Client {
	nc, err := nats.Connect(r.Url)
	if err != nil {
		panic(err)
	}

	return &Client{
		cid:              cid,
		conn:             nc,
		subjects:         r.subjects,
		requests:         r.NumRequests,
		warmupDuration:   r.WarmupDuration,
		shutdownDuration: r.ShutdownDuration,
		rate:             r.Rate,
		received:         make(map[string]int32),
		delivered:        make(map[string]int32),
	}
}

func (r *Runner) printResume() {
	for _, client := range r.clients {
		fmt.Println(client)
	}
}
