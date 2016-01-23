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
	addr := net.UDPAddr{Port: 5678}
	s, err := net.ListenUDP("udp", &addr)
	if err != nil {
		os.Exit(1)
	}

	for {
		b := make([]byte, 128)
		i, sa, e := s.ReadFrom(b)
		log.Println(i, sa, e)

		macLen := 256 / 8
		mac := b[:macLen]
		message := b[macLen:i]

		log.Println(mac, string(message), len(message))

		if len(mac) != macLen {
			log.Println("Short HMAC")
			continue
		}
		if !vonq.CheckMAC(message, mac, []byte("1234")) {
			log.Println("Bad HMAC ", macLen)
			continue
		}

		log.Println("Good HMAC ", macLen)

		req := pb.PingRequest{}
		proto.Unmarshal(message, &req)
		then := time.Unix(0, int64(*req.Time))
		log.Println(req, *req.Test, then)
	}

	s.Close()
}
