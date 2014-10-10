package runner

import (
	"fmt"
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/apcera/nats"
	"github.com/aybabtme/color/brush"
	"github.com/mcuadros/pb"
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
	progressBar      *pb.ProgressBar
	sync.WaitGroup
}

func (r *Runner) Run() {
	r.createProgressBar()
	r.buildSubjects()
	r.buildClients()
	r.runClients()
}

func (r *Runner) buildSubjects() {
	r.subjects = make([]string, r.NumSubjects)
	for i := 0; i < r.NumSubjects; i++ {
		r.subjects[i] = uuid.New()
	}
}

func (r *Runner) buildClients() {
	fmt.Printf("Creating %d client(s) ... ", r.NumClients)
	r.clients = make([]*Client, r.NumClients)
	for i := 0; i < r.NumClients; i++ {
		r.clients[i] = r.buildClient(i)
	}

	fmt.Println(brush.Green("OK"))
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
		requests:         r.NumRequests / r.NumClients / r.NumSubjects,
		warmupDuration:   r.WarmupDuration,
		shutdownDuration: r.ShutdownDuration,
		rate:             r.Rate,
		progressBar:      r.progressBar,
		received:         make(map[string]int32),
		delivered:        make(map[string]int32),
	}
}

func (r *Runner) runClients() {
	fmt.Println("Sending and received messages...")
	r.progressBar.Start()
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

func (r *Runner) printResume() {

}

func (r *Runner) createProgressBar() {
	r.progressBar = pb.New(r.NumRequests * r.NumClients)
	r.progressBar.ShowSpeed = true
	r.progressBar.Width = 80
}
