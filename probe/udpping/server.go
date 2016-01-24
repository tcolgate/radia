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

		if len(mac) != macLen {
			log.Println("Server: Short HMAC")
			continue
		}
		if !checkMAC(message, mac, s.key) {
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
