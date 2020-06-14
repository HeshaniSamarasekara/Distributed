package util

import (
	"Distributed/P2PClient/model"
	"errors"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/magiconair/properties"
)

// Props - Properties
var Props *properties.Properties

// RouteTable - Route table
var RouteTable model.RouteTable

func init() {
	readConfigurations() // Read configuration files
}

// Read configuration from file
func readConfigurations() {
	configFile := flag.String("configFile", "application.yaml", "Configuration File")
	flag.Parse()
	Props = properties.MustLoadFile(*configFile, properties.UTF8)
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

// DecodeRequest - Decodes the response
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
		if n.Ip == node.Ip && n.Port == node.Port {
			return
		}
	}
	RouteTable.Nodes = append(RouteTable.Nodes, node)
}

// RemoveFromRT - Removes stored nodes in Routing table
func RemoveFromRT(node model.Node) {
	var removeNode int
	for i, n := range RouteTable.Nodes {
		if n.Ip == node.Ip && n.Port == node.Port {
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
