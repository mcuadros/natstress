package main

import (
	"flag"
	"time"

	"github.com/apcera/nats"
	"github.com/mcuadros/natstress/runner"
)

var (
	flagServer           = flag.String("h", nats.DefaultURL, "")
	flagSubjects         = flag.Int("s", 5, "")
	flagRequests         = flag.Int("r", 2000, "")
	flagClients          = flag.Int("c", 5, "")
	flagWarmupDuration   = flag.Duration("warmup", 50*time.Millisecond, "")
	flagShutdownDuration = flag.Duration("shutdown", 5*time.Second, "")
)

func main() {
	flag.Parse()

	runner := &runner.Runner{
		Url:              *flagServer,
		NumClients:       *flagClients,
		NumSubjects:      *flagSubjects,
		NumRequests:      *flagRequests,
		WarmupDuration:   *flagWarmupDuration,
		ShutdownDuration: *flagShutdownDuration,
	}

	runner.Run()
}
