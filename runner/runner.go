package runner

import (
	"fmt"
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/apcera/nats"
	"github.com/aybabtme/color/brush"
	"github.com/cheggaaa/pb"
)

type Runner struct {
	Url              string
	NumClients       int
	NumSubjects      int
	NumMessages      int
	ShutdownDuration time.Duration
	Rate             int
	clients          []*Client
	subjects         []string
	progressBar      *pb.ProgressBar
	profiler         *Profiler
	sync.WaitGroup
}

func (r *Runner) Run() {
	r.createProgressBarAndProfiler()
	r.buildSubjects()
	r.buildClients()
	r.subscribeClients()
	r.publishClients()
	r.printResume()
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
		cid:         cid,
		conn:        nc,
		subjects:    r.subjects,
		requests:    r.NumMessages / r.NumClients / r.NumSubjects,
		rate:        r.Rate,
		progressBar: r.progressBar,
		profiler:    r.profiler,
		received:    make(map[string]int32),
		delivered:   make(map[string]int32),
	}
}

func (r *Runner) subscribeClients() {
	fmt.Printf("Subscribing clients to %d subject(s) ... ", r.NumSubjects)
	r.profiler.Start()

	for _, client := range r.clients {
		r.Add(1)
		go func(c *Client) {
			c.Subscribe()
			r.Done()
		}(client)
	}
	r.Wait()
	fmt.Println(brush.Green("OK"))
}

func (r *Runner) publishClients() {
	fmt.Println("Publishing and receiving messages... ")
	r.progressBar.Start()

	r.profiler.Start()
	for _, client := range r.clients {
		r.Add(1)
		go func(c *Client) {
			c.Publish()
			r.Done()
		}(client)
	}

	r.Wait()
	r.profiler.Stop()
	time.Sleep(r.ShutdownDuration)
}

func (r *Runner) printResume() {
	fmt.Printf("\n\nPunlishing summary:\n")
	fmt.Printf("  Count:\t%d messages.\n", r.profiler.count)
	fmt.Printf("  Total:\t%4.4f secs.\n", r.profiler.duration.Seconds())
	fmt.Printf("  Slowest:\t%4.4d µs.\n", r.profiler.max.Nanoseconds()/1000)
	fmt.Printf("  Fastest:\t%4.4d µs.\n", r.profiler.min.Nanoseconds()/1000)
	fmt.Printf("  Average:\t%4.4d µs.\n", r.profiler.avg.Nanoseconds()/1000)
	fmt.Printf("  Messages/sec:\t%4.4f\n", r.profiler.rate)
}

func (r *Runner) createProgressBarAndProfiler() {
	r.progressBar = pb.New(r.NumMessages * r.NumClients)
	r.progressBar.ShowSpeed = true
	r.progressBar.SetWidth(80)
	r.progressBar.Prefix("Received ")

	r.profiler = &Profiler{}
}
