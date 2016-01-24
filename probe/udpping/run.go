package udpping

import (
	"log"
	"net"

	"github.com/tcolgate/vonq/internal/probes"
	"github.com/tcolgate/vonq/internal/reporter"
)

func init() {
	probes.Register(&probe{})
}

type probe struct {
}

func (p *probe) Run(r reporter.Reporter) {
	log.Println("Got HERE")
	c := client{}
	s := server{}

	addr := net.IPv4(127, 0, 0, 1)
	go c.run(net.UDPAddr{IP: addr})
	go s.run()
}
