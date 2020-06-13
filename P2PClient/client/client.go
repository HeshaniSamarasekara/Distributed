package client

import (
	"fmt"
	"net"
	"strconv"
)

var conn *net.UDPConn

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
	regcmd := "00" + strconv.Itoa(count) + cmd

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

// Unregister : Unregister from the network
func Unregister(ip string, port string, username string) error {
	cmd := " UNREG " + ip + " " + port + " " + username
	count := len(cmd) + 5
	regcmd := "00" + strconv.Itoa(count) + cmd

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
