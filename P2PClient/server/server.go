package server

import (
	"Distributed/P2PClient/model"
	"Distributed/P2PClient/util"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

var connection *net.UDPConn

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

		data, err = processRequest(message)

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

func processRequest(message string) ([]byte, error) {

	decodeMessage, err := util.DecodeRequest(message)

	var node model.Node

	var returnErr error
	var returnMessage string

	node.IP = decodeMessage.Ips[0]
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
	case "SER":
		if err != nil {
			returnErr = err
		}
		resp := search(decodeMessage.Ips[2])
		if resp != "" {
			returnMessage = resp
		}
		log.Println("Search received")
	}
	return []byte(returnMessage), returnErr
}

func search(searchString string) string {
	var containFiles []string
	resp := " "

	for _, file := range util.NodeFiles.FileNames {
		if strings.Contains(file, searchString) {
			containFiles = append(containFiles, file)
			resp = resp + file + " "
		}
	}

	if len(containFiles) > 0 {
		cmd := " SEROK " + util.IP + " " + util.Port + resp
		count := len(cmd) + 5
		return fmt.Sprintf("%04d", count) + cmd
	}
	return ""

}
