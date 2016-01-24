// Package udpping is a point to multi-point ping test.
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
	addr := net.IPv4(127, 0, 0, 1)
	uaddr := net.UDPAddr{IP: addr, Port: 5678}

	s := server{key: []byte("1234"), laddr: uaddr}
	go s.run()

	// Should be able to create multiple of these
	c := client{r: r, key: []byte("1234")}
	go c.run(uaddr)

}
