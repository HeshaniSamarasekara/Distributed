package util

import (
	"Distributed/P2PClient/model"
	"errors"
	"flag"
	"fmt"
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

// StoreInRT - Stores the joined nodes in Routing table
func StoreInRT(node model.Node) error {

	RouteTable.Nodes = append(RouteTable.Nodes, node)

	if len(RouteTable.Nodes) > 0 {
		fmt.Println(RouteTable.Nodes[0].Ip + ":" + RouteTable.Nodes[0].Port)
	}
	return nil
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
		addresses = removeIndex(addresses, choose)
		choose = rand.Intn(count - 1)
		randAddresses = append(randAddresses, addresses[choose])
		return randAddresses
	}
	return addresses

}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
