package main

// VoNQ - View of Network Quality
//
// A tool for monitoring network quality between multiple
// locations

import (
	_ "expvar"

	"github.com/tcolgate/vonq/internal/probes"
	"github.com/tcolgate/vonq/internal/reporter"
	_ "github.com/tcolgate/vonq/probe"
)

func main() {
	rr := reporter.New()
	probes.RunAll(rr)
	select {}
}
