package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/apcera/nats"
	"github.com/mcuadros/natstress/runner"
)

var (
	flagServer           = flag.String("h", nats.DefaultURL, "")
	flagSubjects         = flag.Int("s", 5, "")
	flagRequests         = flag.Int("m", 2000, "")
	flagClients          = flag.Int("c", 5, "")
	flagRate             = flag.Int("r", 0, "")
	flagWarmupDuration   = flag.Duration("warmup", 50*time.Millisecond, "")
	flagShutdownDuration = flag.Duration("shutdown", 5*time.Second, "")
)

var usage = `Usage: natstress [options...]

Options:
  -h            NATS server url. (Default: nats://localhost:4222)
  -s            Number of subjects.
  -m            Number of message to send in each subject.
  -c            Number of clients to run concurrently.
  -r            Rate limit, in seconds.
  --warmup      Time to wait before start to deliver messages after connect
                to the server. (Default: 50ms)
  --shutdown    Wait time for received all the messages sent. (Default: 5s)
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	runner := &runner.Runner{
		Url:              *flagServer,
		NumClients:       *flagClients,
		NumSubjects:      *flagSubjects,
		NumMessages:      *flagRequests,
		Rate:             *flagRate,
		WarmupDuration:   *flagWarmupDuration,
		ShutdownDuration: *flagShutdownDuration,
	}

	runner.Run()
}
