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

package udpping

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	pb "github.com/tcolgate/vonq/probes/udpping/proto"
	"github.com/tcolgate/vonq/reporter"
)

type client struct {
	addr net.UDPAddr
	key  []byte
	r    reporter.Reporter
}

func (cl *client) sendPing(c *net.UDPConn, count uint32) error {
	thing := pb.PingRequest{}
	tns := uint64(time.Now().UnixNano())
	thing.Time = &tns
	thing.Count = &count
	bs, err := proto.Marshal(&thing)
	if err != nil {
		panic(err)
	}

	mac := genMAC(bs, cl.key)

	b := append(mac, bs...)
	_, e := c.Write(b)

	return e
}

func (c *client) run(daddr net.UDPAddr) {
	var count uint32
	s, err := net.DialUDP("udp", &c.addr, &daddr)
	if err != nil {
		os.Exit(1)
	}
	go c.getReplies(s)

	for {
		c.sendPing(s, count)
		time.Sleep(time.Second)
		count++
	}

	s.Close()
}

func (c *client) getReplies(so *net.UDPConn) {
	for {
		b := make([]byte, 128)
		i, sa, err := so.ReadFromUDP(b)
		if err != nil {
			log.Println("failed to read from socket, ", err.Error())
			continue
		}

		tns := uint64(time.Now().UnixNano())
		macLen := 256 / 8
		mac := b[:macLen]
		message := b[macLen:i]

		if len(mac) != macLen || !checkMAC(message, mac, c.key) {
			//Should count bad hmacs
			continue
		}

		rep := pb.PingReply{}
		proto.Unmarshal(message, &rep)
		go c.processReply(&rep, tns, sa)
	}
}

func (c *client) processReply(r *pb.PingReply, timeIn uint64, sa *net.UDPAddr) {
	t1 := time.Unix(0, int64(*r.TimeSent))
	t2 := time.Unix(0, int64(*r.TimeIn))
	t3 := time.Unix(0, int64(*r.TimeOut))
	t4 := time.Unix(0, int64(timeIn))

	dout := t2.Sub(t1)
	dsrv := t3.Sub(t2)
	dback := t4.Sub(t3)
	rtt := t4.Sub(t1) - dsrv

	tags := map[string]string{
		"local":  c.addr.IP.String(),
		"remote": sa.IP.String(),
	}

	c.r.Report([]reporter.Metric{
		{Name: "udpping.rttns", Tags: tags, Value: float64(rtt.Nanoseconds())},
		{Name: "udpping.doutns", Tags: tags, Value: float64(dout.Nanoseconds())},
		{Name: "udpping.dsrvns", Tags: tags, Value: float64(dsrv.Nanoseconds())},
		{Name: "udpping.dbackns", Tags: tags, Value: float64(dback.Nanoseconds())},
	})
	return
}
