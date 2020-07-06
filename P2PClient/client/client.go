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

		node.IP = hostPort[0]
		node.Port = hostPort[1]

		err = Join(node.IP, node.Port)

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

	rtLength := len(util.RouteTable.Nodes)

	// Remove from Routing table
	for i := 0; i < rtLength; i++ {

		nodeToRemove := util.RouteTable.Nodes[i]

		var node model.Node

		node.IP = nodeToRemove.IP
		node.Port = nodeToRemove.Port

		err := Leave(node.IP, node.Port)

		if err != nil {
			return err
		}
	}

	util.RouteTable = model.RouteTable{}

	return nil
}

// Join : Join to a node in the network
func Join(ip string, port string) error {

	createPeerConnection(ip, port)

	cmd := " JOIN " + util.Props.MustGetString("ip") + " " + util.Props.MustGetString("port")
	count := len(cmd) + 5
	regcmd := fmt.Sprintf("%04d", count) + cmd

	log.Println(regcmd)

	_, err := util.ReadWriteUDP(regcmd, peerConn)

	if err != nil {
		return err
	}

	closeConnection(peerConn)
	return nil
}

// Leave : Leave from the network
func Leave(ip string, port string) error {

	createPeerConnection(ip, port)

	cmd := " LEAVE " + util.Props.MustGetString("ip") + " " + util.Props.MustGetString("port")
	count := len(cmd) + 5
	regcmd := fmt.Sprintf("%04d", count) + cmd

	log.Println(regcmd)

	resp, err := util.ReadWriteUDP(regcmd, peerConn)

	if err != nil {
		return err
	}

	_, err = util.DecodeResponse(resp)

	if err != nil {
		return err
	}

	log.Println(resp)

	closeConnection(peerConn)
	return nil
}

// Search : Search file in the network
func Search(searchString string) error {
	//ttl := util.Props.MustGetInt("ttl")

	for _, file := range util.NodeFiles.FileNames {
		if strings.Contains(file, searchString) {
			log.Println("File found in this node")
			return nil
		}
	}

	if util.FileTable.Files != nil && len(util.FileTable.Files) > 0 {
		for _, ftEntry := range util.FileTable.Files {
			for _, fileString := range ftEntry.FileStrings {
				if strings.Contains(fileString, searchString) {
					log.Println("File found in  searchString " + ftEntry.IP + ":" + ftEntry.Port)
					return nil
				}

			}
		}
	}

	for _, neighbor := range util.RouteTable.Nodes {
		searchInNetwork(neighbor.IP, neighbor.Port, searchString, "2")
	}

	return nil
}

func searchInNetwork(ip string, port string, filename string, ttl string) error {

	createPeerConnection(ip, port)

	cmd := " SER " + util.Props.MustGetString("ip") + " " + util.Props.MustGetString("port") + " " + filename + " " + ttl
	count := len(cmd) + 5
	regcmd := fmt.Sprintf("%04d", count) + cmd

	log.Println(regcmd)

	resp, err := util.ReadWriteUDP(regcmd, peerConn)

	if err != nil {
		return err
	}

	if resp != "" {
		searchResp, err = util.DecodeResponse(resp)
	}

	if err != nil {
		return err
	}

	util.StoreInFT(searchResp)

	closeConnection(peerConn)
	return nil
}
