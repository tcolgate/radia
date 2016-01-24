package udpping

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	pb "github.com/tcolgate/vonq/probe/udpping/proto"
)

type client struct {
	addr net.UDPAddr
}

func sendPing(c *net.UDPConn) error {
	thing := pb.PingRequest{}
	tns := uint64(time.Now().UnixNano())
	thing.Time = &tns
	bs, err := proto.Marshal(&thing)
	if err != nil {
		panic(err)
	}

	mac := genMAC(bs, []byte("1234"))
	log.Println(len(mac))

	b := append(mac, bs...)
	_, e := c.Write(b)

	return e
}

func (c *client) run(daddr net.UDPAddr) {
	s, err := net.DialUDP("udp", &c.addr, &daddr)
	if err != nil {
		os.Exit(1)
	}

	for {
		sendPing(s)
		time.Sleep(time.Second)
	}

	s.Close()
}
