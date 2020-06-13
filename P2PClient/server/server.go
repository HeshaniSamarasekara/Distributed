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

		data := []byte("")

		data, err = updateRT(message)

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

func updateRT(message string) ([]byte, error) {

	decodeMessage, err := util.DecodeRequest(message)

	var node model.Node

	var returnErr error
	var returnMessage string

	node.Ip = decodeMessage.Ips[0]
	node.Port = decodeMessage.Ips[1]

	switch decodeMessage.Code {
	case "JOIN":
		if err != nil {
			returnErr = err
			returnMessage = "0016 JOINOK 9999"
		}
		util.StoreInRT(node)
		returnMessage = "0013 JOINOK 0"
		break
	case "LEAVE":
		if err != nil {
			returnErr = err
			returnMessage = "0017 LEAVEOK 9999"
		}
		util.RemoveFromRT(node)
		returnMessage = "0014 LEAVEOK 0"
		break
	}
	return []byte(returnMessage), returnErr
}
