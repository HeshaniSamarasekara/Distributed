package util

import (
	"Distributed/P2PClient/model"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/magiconair/properties"
)

// Props - Properties
var Props *properties.Properties

// RouteTable - Route table
var RouteTable model.RouteTable

// FileTable - File table
var FileTable model.FileTable

// NodeFiles - Files in the own node
var NodeFiles model.NodeFiles

func init() {
	readConfigurations() // Read configuration files
	readFileNames()      // Read file names from list
}

// Read configuration from file
func readConfigurations() {
	configFile := flag.String("configFile", "application.yaml", "Configuration File")
	flag.Parse()
	Props = properties.MustLoadFile(*configFile, properties.UTF8)
}

func readFileNames() {
	readFile := flag.String("fileNames", "FileNames.txt", "File names for nodes")
	flag.Parse()
	data, err := ioutil.ReadFile(*readFile)
	if err != nil {
		log.Println("Error reading file names.")
	}
	allFiles := strings.Split(string(data), "\n")
	from := float64(randomInt(0, len(allFiles)-1))
	to := float64(randomInt(0, len(allFiles)-1))
	NodeFiles.FileNames = allFiles[int(math.Min(from, to)):int(math.Max(from, to))]
	log.Println(strings.Join(NodeFiles.FileNames, ", "))
}

// ValidateErrorCode - Error from BS
func validateErrorCode(code string) error {
	switch code {
	case "9999":
		log.Println("failed, there is some error in the command")
		return errors.New("failed, there is some error in the command")
	case "9998":
		log.Println("failed, already registered to you, unregister first")
		return errors.New("failed, already registered to you, unregister first")
	case "9997":
		log.Println("failed, registered to another user, try a different IP and port")
		return errors.New("failed, registered to another user, try a different IP and port")
	case "9996":
		log.Println("failed, can’t register. BS full")
		return errors.New("failed, can’t register. BS full")
	}
	return nil
}

// DecodeResponse - Decodes the response
func DecodeResponse(reply string) (model.Response, error) {
	splittedReply := strings.Split(reply, " ")

	err := validateErrorCode(splittedReply[2])

	if err != nil {
		return model.Response{}, err
	}

	response := model.Response{}
	response.Length = splittedReply[0]
	response.Code = splittedReply[1]
	response.Count = splittedReply[2]
	for i := 3; i < len(splittedReply); i++ {
		response.Ips = append(response.Ips, splittedReply[i])
	}

	return response, nil
}

// DecodeRequest - Decodes the request
func DecodeRequest(reply string) (model.Response, error) {
	splittedReply := strings.Split(reply, " ")

	response := model.Response{}
	response.Length = splittedReply[0]
	response.Code = splittedReply[1]
	for i := 2; i < len(splittedReply); i++ {
		response.Ips = append(response.Ips, splittedReply[i])
	}

	return response, nil
}

// StoreInRT - Stores the joined nodes in Routing table
func StoreInRT(node model.Node) {
	for _, n := range RouteTable.Nodes {
		if n.IP == node.IP && n.Port == node.Port {
			return
		}
	}
	RouteTable.Nodes = append(RouteTable.Nodes, node)
}

// RemoveFromRT - Removes stored nodes in Routing table
func RemoveFromRT(node model.Node) {
	var removeNode int
	for i, n := range RouteTable.Nodes {
		if n.IP == node.IP && n.Port == node.Port {
			removeNode = i
			break
		}
	}
	RouteTable.Nodes = append(RouteTable.Nodes[:removeNode], RouteTable.Nodes[removeNode+1:]...)
}

// RandomPeer - Select random peers
func RandomPeer(reply model.Response) []string {
	count, _ := strconv.Atoi(reply.Count)
	addresses := []string{}
	randAddresses := []string{}

	for i := 0; i < len(reply.Ips); i += 2 {
		addresses = append(addresses, reply.Ips[i]+":"+reply.Ips[i+1])
	}

	if count > 2 {
		choose := rand.Intn(count)
		randAddresses = append(randAddresses, addresses[choose])
		addresses = append(addresses[:choose], addresses[choose+1:]...)
		choose = rand.Intn(count - 1)
		randAddresses = append(randAddresses, addresses[choose])
		return randAddresses
	}
	return addresses
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// ReadWriteUDP - Send UDP request and get the response back
func ReadWriteUDP(regcmd string, peerConn *net.UDPConn) ([]byte, int, error) {
	regbytes := []byte(regcmd)
	buffer := make([]byte, 1024)

	_, err := peerConn.Write(regbytes)
	if err != nil {
		log.Println(err)
		return []byte{}, 0, err
	}

	n, _, err := peerConn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return []byte{}, 0, err
	}
	log.Println("Reply: ", string(buffer[0:n]))
	return regbytes, n, nil
}
