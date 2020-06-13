package client

import (
	"Distributed/P2PClient/model"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var conn *net.UDPConn
var routeTable model.RouteTable

// CreateConnection : Creates UDP connection
func CreateConnection() {
	s, err := net.ResolveUDPAddr("udp4", "localhost:55555")
	conn, err = net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
}

func closeConnection() {
	conn.Close()
}

// Register : Register to the network
func Register(ip string, port string, username string) error {
	cmd := " REG " + ip + " " + port + " " + username
	count := len(cmd) + 5
	regcmd := fmt.Sprintf("%04d", count) + cmd

	fmt.Println(regcmd)

	regbytes := []byte(regcmd)
	buffer := make([]byte, 1024)

	_, err := conn.Write(regbytes)
	if err != nil {
		fmt.Println(err)
		return err
	}

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println(err)
		return err
	}

	reply := string(buffer[0:n])
	fmt.Printf("Reply: %s\n", reply)

	storeInRT(reply)

	return nil
}

// Unregister : Unregister from the network
func Unregister(ip string, port string, username string) error {
	cmd := " UNREG " + ip + " " + port + " " + username
	count := len(cmd) + 5
	regcmd := fmt.Sprintf("%04d", count) + cmd

	fmt.Println(regcmd)

	regbytes := []byte(regcmd)
	buffer := make([]byte, 1024)

	_, err := conn.Write(regbytes)
	if err != nil {
		fmt.Println(err)
		return err
	}

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	reply := string(buffer[0:n])
	fmt.Printf("Reply: %s\n", reply)

	return nil
}

// Join : Join to a node in the network
func Join(ip string, port string) error {
	cmd := " JOIN " + ip + " " + port
	count := len(cmd) + 5
	regcmd := fmt.Sprintf("%04d", count) + cmd

	fmt.Println(regcmd)

	regbytes := []byte(regcmd)
	buffer := make([]byte, 1024)

	_, err := conn.Write(regbytes)
	if err != nil {
		fmt.Println(err)
		return err
	}

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Reply: %s\n", string(buffer[0:n]))
	return nil
}

func storeInRT(reply string) error {
	splittedReply := strings.Split(reply, " ")

	nodecount, err := strconv.Atoi(splittedReply[2])

	if err != nil {
		return err
	}

	for i := 0; i < nodecount; i++ {

		node := model.Node{}

		node.Ip = splittedReply[3+i]
		node.Port = splittedReply[4+i]

		err = Join(node.Ip, node.Port)

		if err != nil {
			return err
		}
		routeTable.Nodes = append(routeTable.Nodes, node)

	}

	if len(routeTable.Nodes) > 0 {
		fmt.Println(routeTable.Nodes[0].Ip + ":" + routeTable.Nodes[0].Port)
	}
	return nil
}
