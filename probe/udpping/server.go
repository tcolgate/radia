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

	pb "github.com/tcolgate/vonq/probe/udpping/proto"
)

type server struct {
	laddr net.UDPAddr
	key   []byte
}

func (s *server) run() {
	so, err := net.ListenUDP("udp", &s.laddr)
	if err != nil {
		os.Exit(1)
	}

	for {
		b := make([]byte, 128)
		i, sa, e := so.ReadFromUDP(b)
		log.Println(i, sa, e)

		tns := uint64(time.Now().UnixNano())
		macLen := 256 / 8
		mac := b[:macLen]
		message := b[macLen:i]

		if len(mac) != macLen || !checkMAC(message, mac, s.key) {
			log.Println("Server: Bad HMAC ", macLen)
			continue
		}

		req := pb.PingRequest{}
		proto.Unmarshal(message, &req)
		go s.process(so, sa, &req, tns)
	}

	so.Close()
}

func (s *server) process(so *net.UDPConn, sa *net.UDPAddr, r *pb.PingRequest, timeIn uint64) {
	tns := uint64(time.Now().UnixNano())

	thing := pb.PingReply{}
	thing.TimeSent = r.Time
	thing.TimeIn = &timeIn
	thing.Count = r.Count
	thing.TimeOut = &tns

	bs, err := proto.Marshal(&thing)
	if err != nil {
		panic(err)
	}

	mac := genMAC(bs, s.key)

	b := append(mac, bs...)
	so.WriteToUDP(b, sa)

	return
}
