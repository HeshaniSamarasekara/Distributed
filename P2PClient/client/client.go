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

// CreateBootstrapConnection : Creates UDP connection with Bootstrap
func CreateBootstrapConnection() {
	s, err := net.ResolveUDPAddr("udp4", util.Props.MustGetString("bootstrapIp")+":"+util.Props.MustGetString("bootstrapPort"))
	conn, err = net.DialUDP("udp4", nil, s)
	if err != nil {
		log.Println(err)
		return
	}
	AutomaticRegister()
	log.Println("The Bootstrap server is ", conn.RemoteAddr().String())
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

	rtLength := len(util.GetRT().Nodes)

	// Remove from Routing table
	for i := 0; i < rtLength; i++ {

		nodeToRemove := util.GetRT().Nodes[i]

		var node model.Node

		node.IP = nodeToRemove.IP
		node.Port = nodeToRemove.Port

		err := Leave(node.IP, node.Port)

		if err != nil {
			return err
		}
	}

	util.SetRT(model.RouteTable{})

	updateFileEntryTable(ip, port)

	return nil
}

// Join : Join to a node in the network
func Join(ip string, port string) error {

	s, err1 := net.ResolveUDPAddr("udp4", ip+":"+port)
	if err1 != nil {
		return err1
	}
	peerConn, err2 := net.DialUDP("udp4", nil, s)
	if err2 != nil {
		return err2
	}

	log.Println("The UDP server is ", peerConn.RemoteAddr().String())
	defer closeConnection(peerConn)

	cmd := " JOIN " + util.IP + " " + util.Port
	count := len(cmd) + 5
	regcmd := fmt.Sprintf("%04d", count) + cmd

	log.Println(regcmd)

	_, err := util.ReadWriteUDP(regcmd, peerConn)

	if err != nil {
		return err
	}
	return nil
}

// Leave : Leave from the network
func Leave(ip string, port string) error {

	s, err1 := net.ResolveUDPAddr("udp4", ip+":"+port)
	if err1 != nil {
		return err1
	}
	peerConn, err2 := net.DialUDP("udp4", nil, s)
	if err2 != nil {
		return err2
	}
	log.Println("The UDP server is ", peerConn.RemoteAddr().String())
	defer closeConnection(peerConn)

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

	for _, neighbor := range util.GetRT().Nodes {
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

	done := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("SEARCH DONE IN NETWORK")
		return searchInNode(searchString), nil
	case <-ctx.Done():
		fmt.Println("TIME OUT")
		cancel()
		cmd := " SEROK 0 " + util.IP + " " + util.Port + " 0 "
		count := len(cmd) + 5
		returnCmd := fmt.Sprintf("%04d", count) + cmd
		log.Println(returnCmd)
		return returnCmd, nil
	}
}

func searchInNetwork(wg *sync.WaitGroup, ip string, port string, filename string, hops int) {
	stored := false
	s, err1 := net.ResolveUDPAddr("udp4", ip+":"+port)
	peerConn, err2 := net.DialUDP("udp4", nil, s)
	log.Println("The UDP server is ", peerConn.RemoteAddr().String())

	if err1 != nil || err2 != nil {
		log.Println(err1)
		log.Println(err2)
		updateRoutingTable(ip, port)
		updateFileEntryTable(ip, port)
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
		updateFileEntryTable(ip, port)
		log.Println(err)
		return
	}

	if strings.TrimSpace(resp) != "" {
		searchResp, _ := util.DecodeSearchResponse(resp)
		if searchResp.Count > 0 {
			stored = true
			util.StoreInFT(wg, searchResp)
		}
	}

	if err != nil {
		log.Println(err)
	}

	if !stored {
		wg.Done()
	}

}

func updateFileEntryTable(ip string, port string) {
	localFt := util.GetFT()
	if 0 == len(localFt.Files) {
		return
	}
	if 1 == len(localFt.Files) {
		localFt.Files = localFt.Files[:0]
	}
	for i, node := range util.GetFT().Files {
		if node.IP+":"+node.Port == ip+":"+port {
			localFt.Files = append(localFt.Files[:i], localFt.Files[i+1:]...)
			log.Println("removing entry from File table " + ip + ":" + port)
			break
		}
	}
	util.SetFT(localFt)
}

func updateRoutingTable(ip string, port string) {
	localRT := util.GetRT()
	if 0 == len(localRT.Nodes) {
		return
	}
	if 1 == len(localRT.Nodes) {
		localRT.Nodes = localRT.Nodes[:0]
	}
	for i, node := range util.GetRT().Nodes {
		if node.IP+":"+node.Port == ip+":"+port {
			localRT.Nodes = append(localRT.Nodes[:i], localRT.Nodes[i+1:]...)
			log.Println("removing entry from Route table " + ip + ":" + port)
			break
		}
	}
	util.SetRT(localRT)
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

	localFt := util.GetFT()
	if localFt.Files != nil && len(localFt.Files) > 0 {
		for _, ftEntry := range localFt.Files {
			if strings.Contains(ftEntry.FileStrings, searchString) {
				fileNames := strings.Split(ftEntry.FileStrings, ",")
				correcteNames := ""
				for _, n := range fileNames {
					if strings.Contains(n, searchString) {
						correcteNames += n
					}
				}
				log.Println("File " + searchString + " can be found in " + ftEntry.IP + ":" + ftEntry.Port)
				cmd := " SEROK " + fmt.Sprintf("%d", 1) + " " + ftEntry.IP + " " + ftEntry.Port + " 0 " + correcteNames
				count := len(cmd) + 5
				returnCmd := fmt.Sprintf("%04d", count) + cmd
				return returnCmd
			}
		}
	}
	return ""
}

func AutomaticRegister() {
	err2 := Register(util.IP, util.Port, util.Name)
	if err2 != nil {
		log.Println(err2)
		err1 := Unregister(util.IP, util.Port, util.Name)
		if err1 != nil {
			log.Println(err1)
		}
		err3 := Register(util.IP, util.Port, util.Name)
		if err3 != nil {
			log.Println(err3)
		}
	}
}
