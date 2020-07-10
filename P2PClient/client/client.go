package client

import (
	"Distributed/P2PClient/model"
	"Distributed/P2PClient/util"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
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
	log.Println("The Bootstrap server is ", conn.RemoteAddr().String())
}

// CreatePeerConnection : Creates UDP connection
func createPeerConnection(ip string, port string) error {
	s, err := net.ResolveUDPAddr("udp4", ip+":"+port)
	if err != nil {
		log.Println(err)
		return err
	}
	peerConn, err = net.DialUDP("udp4", nil, s)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("The UDP server is ", peerConn.RemoteAddr().String())

	return nil
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

	cmd := " JOIN " + util.IP + " " + util.Port
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

	cmd := " LEAVE " + util.IP + " " + util.Port
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
func Search(searchString string, incomingHostPort string, hopCount int) (string, error) {
	var wg sync.WaitGroup
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Duration(util.TTL)*time.Second)
	defer cancel()

	if isFileInNode := searchInNode(searchString); isFileInNode != "" {
		return isFileInNode, nil
	}

	for _, neighbor := range util.RouteTable.Nodes {
		// This is goroutine. Concurrently executes.
		if incomingHostPort == "" || (neighbor.IP+":"+neighbor.Port) != incomingHostPort {
			var hops int
			if hopCount == 9999 {
				hops = util.Hops
			} else {
				hops = hopCount
			}
			wg.Add(1)
			go searchInNetwork(&wg, neighbor.IP, neighbor.Port, searchString, hops)
		}
	}

	wg.Wait()

	select {
	case <-ctx.Done():
		fmt.Println("TIME OUT")
		cancel()
		return "TIME OUT", ctx.Err()
	default:
		fmt.Println("ALL DONE")
		return searchInNode(searchString), nil
	}
}

func searchInNetwork(wg *sync.WaitGroup, ip string, port string, filename string, hops int) {

	defer wg.Done()
	errm := createPeerConnection(ip, port)

	if errm != nil {
		updateRoutingTable(ip, port)
		return
	}

	defer closeConnection(peerConn)

	cmd := " SER " + util.IP + " " + util.Port + " " + filename + " " + fmt.Sprintf("%d", hops)
	count := len(cmd) + 5
	sercmd := fmt.Sprintf("%04d", count) + cmd

	log.Println(sercmd)

	resp, err := util.ReadWriteUDP(sercmd, peerConn)

	if err != nil {
		updateRoutingTable(ip, port)
		log.Println(err)
		return
	}

	if strings.TrimSpace(resp) != "" {
		util.MU.Lock()
		defer util.MU.Unlock()
		log.Println(resp)
		searchResp, _ := util.DecodeSearchResponse(resp)
		if searchResp.Count > 0 {
			util.StoreInFT(searchResp)
		}
	}

	if err != nil {
		log.Println(err)
	}
}

func updateRoutingTable(ip string, port string) {
	for i, node := range util.RouteTable.Nodes {
		if node.IP == ip && node.Port == port {
			util.RouteTable.Nodes = append(util.RouteTable.Nodes[:i], util.RouteTable.Nodes[i+1:]...)
			log.Println("removing node" + ip + port)
			break
		}
	}
}

func searchInNode(searchString string) string {
	contains := ""

	for _, file := range util.NodeFiles.FileNames {
		if strings.Contains(file, searchString) {
			log.Println("File found in this node")
			contains += file + ","
		}
	}

	if contains != "" {
		cmd := " SEROK " + fmt.Sprintf("%d", len(strings.Split(contains, ","))) + " " + util.IP + " " + util.Port + " 0 " + strings.ReplaceAll(contains, ",", " ")
		count := len(cmd) + 5
		returnCmd := fmt.Sprintf("%04d", count) + cmd
		return returnCmd
	}

	util.MU.Lock()
	defer util.MU.Unlock()
	if util.FileTable.Files != nil && len(util.FileTable.Files) > 0 {
		for _, ftEntry := range util.FileTable.Files {
			if strings.Contains(ftEntry.FileStrings, searchString) {
				log.Println("File " + searchString + "can be found in " + ftEntry.IP + ":" + ftEntry.Port)
				cmd := " SEROK " + fmt.Sprintf("%d", 1) + " " + ftEntry.IP + " " + ftEntry.Port + " 0 " + ftEntry.FileStrings
				count := len(cmd) + 5
				returnCmd := fmt.Sprintf("%04d", count) + cmd
				return returnCmd
			}
		}
	}
	return ""
}
