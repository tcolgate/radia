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
		i, sa, e := so.ReadFrom(b)
		log.Println(i, sa, e)

		macLen := 256 / 8
		mac := b[:macLen]
		message := b[macLen:i]

		if len(mac) != macLen {
			log.Println("Short HMAC")
			continue
		}
		if !checkMAC(message, mac, s.key) {
			log.Println("Bad HMAC ", macLen)
			continue
		}

		log.Println("Good HMAC ", macLen)

		req := pb.PingRequest{}
		proto.Unmarshal(message, &req)
		then := time.Unix(0, int64(*req.Time))
		log.Println(req, *req.Test, then)
	}

	so.Close()
}
