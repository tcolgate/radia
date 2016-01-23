package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/tcolgate/vonq"

	pb "github.com/tcolgate/vonq/proto"
)

func main() {
	addr := net.UDPAddr{}
	s, err := net.ListenUDP("udp", &addr)
	if err != nil {
		os.Exit(1)
	}
	daddr := net.UDPAddr{Port: 5678}
	daddr.IP = net.IPv4(127, 0, 0, 1)

	for {
		thing := pb.PingRequest{}
		thing.Test = proto.String("hello")
		tns := uint64(time.Now().UnixNano())
		thing.Time = &tns
		bs, err := proto.Marshal(&thing)
		if err != nil {
			os.Exit(1)
		}

		mac := vonq.GenMAC(bs, []byte("1234"))
		log.Println(len(mac))

		b := append(mac, bs...)
		i, e := s.WriteToUDP(b, &daddr)
		log.Println(i, b, e)

		time.Sleep(time.Second)

	}

	s.Close()
}
