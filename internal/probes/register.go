package probes

import "github.com/tcolgate/vonq/internal/reporter"

type Runner interface {
	Run(reporter.Reporter)
}

var probes []Runner

func Register(p Runner) {
	probes = append(probes, p)
}

func RunAll(r reporter.Reporter) {
	for _, p := range probes {
		go p.Run(r)
	}
}
