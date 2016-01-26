// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of vonq.
//
// vonq is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// vonq is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with vonq.  If not, see <http://www.gnu.org/licenses/>.

// Package udpping is a point to multi-point ping test.
package udpping

import (
	"net"

	"github.com/tcolgate/vonq/probes/base"
	"github.com/tcolgate/vonq/probes/register"
)

func init() {
	register.Probe(Init)
}

type udpPingProbe struct {
	*base.Base
	key string
}

func Init() base.Probe {
	p := udpPingProbe{}
	p.Base, _ = base.New()
	p.StringVar(&p.key, "key", "secret", "Signing secret")
	return &p
}

func (p *udpPingProbe) Run(args ...string) {
	addr := net.IPv4(127, 0, 0, 1)
	uaddr := net.UDPAddr{IP: addr, Port: 5678}

	s := server{udpPingProbe: p, laddr: uaddr}
	go s.runServer()

	// Should be able to create multiple of these
	c := client{udpPingProbe: p}
	go c.runClient(uaddr)

}
