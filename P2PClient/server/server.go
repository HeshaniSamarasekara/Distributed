package server

import (
	"Distributed/P2PClient/model"
	"Distributed/P2PClient/util"
	"log"
	"math/rand"
	"net"
	"time"
)

var connection *net.UDPConn

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// CreateServer - Creates UDP server
func CreateServer() {
	s, err := net.ResolveUDPAddr("udp4", util.Props.MustGetString("ip")+":"+util.Props.MustGetString("port"))
	if err != nil {
		log.Println(err)
		return
	}

	connection, err = net.ListenUDP("udp4", s)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Starting UDP server at port " + util.Props.MustGetString("port"))

	defer connection.Close()

	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())

	for {
		n, addr, err := connection.ReadFromUDP(buffer)

		message := string(buffer[0:n])

		log.Print("-> ", message)

		data := []byte("0014 JOINOK 0")

		decodeMessage, err := util.DecodeResponse(message)

		if err != nil {
			data = []byte("0014 JOINOK 9999")
		}

		var node model.Node

		node.Ip = decodeMessage.Ips[0]
		node.Port = decodeMessage.Ips[1]

		err = util.StoreInRT(node)

		if err != nil {
			data = []byte("0014 JOINOK 9999")
		}

		log.Println("data: ", string(data))

		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func Shutdown() {
	connection.Close()
}
