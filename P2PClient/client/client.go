package client

import (
	"Distributed/P2PClient/model"
	"Distributed/P2PClient/util"
	"fmt"
	"log"
	"net"
	"strings"
)

var conn *net.UDPConn
var peerConn *net.UDPConn

// CreateConnection : Creates UDP connection with Bootstrap
func CreateConnection() {
	s, err := net.ResolveUDPAddr("udp4", util.Props.MustGetString("bootstrapIp")+":"+util.Props.MustGetString("bootstrapPort"))
	conn, err = net.DialUDP("udp4", nil, s)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("The UDP server is ", conn.RemoteAddr().String())
}

// CreatePeerConnection : Creates UDP connection
func createPeerConnection(ip string, port string) {
	s, err := net.ResolveUDPAddr("udp4", ip+":"+port)
	peerConn, err = net.DialUDP("udp4", nil, s)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("The UDP server is ", peerConn.RemoteAddr().String())
}

func closeConnection(connect *net.UDPConn) {
	connect.Close()
}

// Register : Register to the network
func Register(ip string, port string, username string) error {
	cmd := " REG " + ip + " " + port + " " + username
	count := len(cmd) + 5
	regcmd := fmt.Sprintf("%04d", count) + cmd

	log.Println(regcmd)

	regbytes := []byte(regcmd)
	buffer := make([]byte, 1024)

	_, err := conn.Write(regbytes)
	if err != nil {
		log.Println(err)
		return err
	}

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return err
	}

	reply := string(buffer[0:n])
	log.Println("Reply: ", reply)

	response, err := util.DecodeResponse(reply)
	if err != nil {
		return err
	}

	peersToConnect := util.RandomPeer(response)

	for i := 0; i < len(peersToConnect); i++ {

		var node model.Node

		hostPort := strings.Split(peersToConnect[i], ":")

		node.Ip = hostPort[0]
		node.Port = hostPort[1]

		err = Join(node.Ip, node.Port)

		if err != nil {
			return err
		}

		util.StoreInRT(node)
	}

	return nil
}

// Unregister : Unregister from the network
func Unregister(ip string, port string, username string) error {
	cmd := " UNREG " + ip + " " + port + " " + username
	count := len(cmd) + 5
	regcmd := fmt.Sprintf("%04d", count) + cmd

	log.Println(regcmd)

	regbytes := []byte(regcmd)
	buffer := make([]byte, 1024)

	_, err := conn.Write(regbytes)
	if err != nil {
		log.Println(err)
		return err
	}

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return err
	}
	reply := string(buffer[0:n])
	log.Println("Reply: ", reply)

	return nil
}

// Join : Join to a node in the network
func Join(ip string, port string) error {

	createPeerConnection(ip, port)

	cmd := " JOIN " + util.Props.MustGetString("ip") + " " + util.Props.MustGetString("port")
	count := len(cmd) + 5
	regcmd := fmt.Sprintf("%04d", count) + cmd

	log.Println(regcmd)

	regbytes := []byte(regcmd)
	buffer := make([]byte, 1024)

	_, err := peerConn.Write(regbytes)
	if err != nil {
		log.Println(err)
		return err
	}

	n, _, err := peerConn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Reply: ", string(buffer[0:n]))

	closeConnection(peerConn)
	return nil
}
