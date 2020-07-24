package server

import (
	"Distributed/P2PClient/client"
	"Distributed/P2PClient/model"
	"Distributed/P2PClient/util"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var connection *net.UDPConn

// CreateServer - Creates UDP server
func CreateServer() {

	s, err := net.ResolveUDPAddr("udp4", util.IP+":"+util.Port)

	if err != nil {
		log.Println(err)
		return
	}

	connection, err = net.ListenUDP("udp4", s)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Starting UDP server at port " + util.Port)

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
		searchStr := decodeMessage.Ips[2]
		hopCount, _ := strconv.Atoi(decodeMessage.Ips[3])
		log.Println(hopCount)
		if hopCount > 0 {
			resp := search(searchStr, hopCount-1, node.IP, node.Port)
			if resp != "" {
				returnMessage = resp
			}
			log.Println("Search received")
		} else {
			cmd := " SEROK 0 " + util.IP + " " + util.Port
			returnMessage = fmt.Sprintf("%04d", len(cmd)+5) + cmd
			log.Println("Hop count exceeded.")
		}
	}
	return []byte(returnMessage), returnErr
}

func search(searchString string, hopCount int, incomingIP string, incomingPort string) string {
	var containFiles []string
	resp := ""

	for _, file := range util.NodeFiles.FileNames {
		if strings.Contains(strings.ToLower(file), strings.ToLower(searchString)) {
			containFiles = append(containFiles, file)
			resp = resp + " " + strings.Join(strings.Split(file, " "), "_")
		}
	}
	if len(containFiles) == 0 {
		log.Println("No files found on " + util.IP + ":" + util.Port)
		resp, err := client.Search(searchString, incomingIP+":"+incomingPort, hopCount)
		if err != nil {
			log.Println(err)
			return ""
		}
		if strings.TrimSpace(resp) != "" {
			responseModel, _ := util.DecodeSearchResponse(resp)
			if responseModel.Count > 0 {
				return resp
			}
		}
	}

	if len(containFiles) > 0 {
		cmd := " SEROK " + fmt.Sprintf("%d", len(containFiles)) + " " + util.IP + " " + util.Port + " " + fmt.Sprintf("%d", hopCount) + resp
		count := len(cmd) + 5
		returnCmd := fmt.Sprintf("%04d", count) + cmd
		return returnCmd
	}
	cmd := " SEROK 0 " + util.IP + " " + util.Port + " 0 "
	count := len(cmd) + 5
	returnCmd := fmt.Sprintf("%04d", count) + cmd
	return returnCmd

}
