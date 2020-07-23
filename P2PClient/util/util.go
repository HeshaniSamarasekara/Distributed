package util

import (
	"Distributed/P2PClient/model"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/magiconair/properties"
)

// Props - Properties
var Props *properties.Properties

// RouteTable - Route table
var routeTable model.RouteTable

// FileTable - File table
var fileTable model.FileTable

// NodeFiles - Files in the own node
var NodeFiles model.NodeFiles

// IP - My IP
var IP string

// Port - My Port
var Port string

// Name - My name
var Name string

// TTL - My TTL
var TTL int

// Hops - My Hop count
var Hops int

// MuFT - Mutex to update file table
var MuFT sync.Mutex

// MuRT - Mutex to update route table
var MuRT sync.Mutex

func init() {
	readConfigurations() // Read configuration files
	readFileNames()      // Read file names from list
	argIP := os.Args[1:]
	if len(argIP) >= 3 {
		IP = argIP[0]
		Port = argIP[1]
		Name = argIP[2]
	}
}

// Read configuration from file
func readConfigurations() {
	configFile := flag.String("configFile", "application.yaml", "Configuration File")
	flag.Parse()
	Props = properties.MustLoadFile(*configFile, properties.UTF8)
	IP = Props.MustGetString("ip")
	Port = Props.MustGetString("port")
	TTL = Props.MustGetInt("ttl")
	Hops = Props.MustGetInt("hopcount")
}

func readFileNames() {
	readFile := flag.String("fileNames", "FileNames.txt", "File names for nodes")
	flag.Parse()
	data, err := ioutil.ReadFile(*readFile)
	if err != nil {
		log.Println("Error reading file names...")
	}
	fileStrings := string(data)
	allFiles := strings.Split(strings.ReplaceAll(strings.ReplaceAll(fileStrings, "\r", ""), " ", "_"), "\n")
	from := float64(randomInt(0, len(allFiles)-1))
	to := float64(randomInt(0, len(allFiles)-1))
	NodeFiles.FileNames = allFiles[int(math.Min(from, to)):int(math.Max(from, to))]
	log.Println(len(NodeFiles.FileNames))
	log.Println(NodeFiles.FileNames)
}

// ValidateErrorCode - Error from BS
func validateErrorCode(code string) error {
	switch code {
	case "9999":
		log.Println("failed, there is some error in the command...")
		return errors.New("failed, there is some error in the command")
	case "9998":
		log.Println("failed, already registered to you, unregister first...")
		return errors.New("failed, already registered to you, unregister first")
	case "9997":
		log.Println("failed, registered to another user, try a different IP and port...")
		return errors.New("failed, registered to another user, try a different IP and port")
	case "9996":
		log.Println("failed, can’t register. BS full...")
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

// DecodeSearchResponse - Decodes the search response
func DecodeSearchResponse(reply string) (model.SearchResponse, error) {
	splittedReply := strings.Split(reply, " ")

	err := validateErrorCode(splittedReply[2])

	if err != nil {
		return model.SearchResponse{}, err
	}

	response := model.SearchResponse{}
	response.Length = splittedReply[0]
	response.Code = splittedReply[1]
	response.Count, _ = strconv.Atoi(splittedReply[2])
	response.IP = splittedReply[3]
	response.Port = splittedReply[4]
	if response.Count > 0 {
		response.Hops = splittedReply[5]
		for i := 6; i < (6 + response.Count); i++ {
			response.Files = append(response.Files, splittedReply[i])
		}
	}
	log.Println(response)

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
	// MuRT.Lock()
	// defer MuRT.Unlock()
	for _, n := range GetRT().Nodes {
		if n.IP == node.IP && n.Port == node.Port {
			return
		}
	}
	localRT := GetRT()
	localRT.Nodes = append(localRT.Nodes, node)
	SetRT(localRT)
}

// RemoveFromRT - Removes stored nodes in Routing table
func RemoveFromRT(node model.Node) {
	// MuRT.Lock()
	// defer MuRT.Unlock()
	var removeNode int
	for i, n := range GetRT().Nodes {
		if n.IP == node.IP && n.Port == node.Port {
			removeNode = i
			break
		}
	}
	localRT := GetRT()
	localRT.Nodes = append(localRT.Nodes[:removeNode], localRT.Nodes[removeNode+1:]...)
	SetRT(localRT)
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
func ReadWriteUDP(regcmd string, peerConn *net.UDPConn) (string, error) {
	regbytes := []byte(regcmd)
	buffer := make([]byte, 1024)

	_, err := peerConn.Write(regbytes)
	if err != nil {
		log.Println(err)
		return "", err
	}

	n, _, err := peerConn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return "", err
	}
	reply := string(buffer[0:n])
	log.Println("Reply: ", reply)
	return reply, nil
}

// StoreInFT - Stores the files in File table
func StoreInFT(wg *sync.WaitGroup, response model.SearchResponse) {
	defer wg.Done()
	localFT := GetFT()
	for i, f := range localFT.Files {
		stringsToAdd := ""
		if f.IP+":"+f.Port == response.IP+":"+response.Port {
			for _, incoming := range response.Files {
				if !strings.Contains(f.FileStrings, incoming) {
					stringsToAdd += "," + incoming
				}
			}
			localFT.Files[i].FileStrings += stringsToAdd
			SetFT(localFT)
			return
		}
	}
	newFileEntry := model.FileTableEntry{}
	newFileEntry.IP = response.IP
	newFileEntry.Port = response.Port
	newFileEntry.FileStrings = strings.Join(response.Files, ",")

	localFT.Files = append(localFT.Files, newFileEntry)
	log.Println("Stored in FT : ", localFT)
	SetFT(localFT)
}

// GetRT - Return route table
func GetRT() model.RouteTable {
	MuRT.Lock()
	defer MuRT.Unlock()
	return routeTable
}

// SetRT - Set new route table
func SetRT(rt model.RouteTable) {
	MuRT.Lock()
	defer MuRT.Unlock()
	routeTable = rt
}

// GetFT - Return file table
func GetFT() model.FileTable {
	MuFT.Lock()
	defer MuFT.Unlock()
	return fileTable
}

// SetFT - Set new file table
func SetFT(ft model.FileTable) {
	MuFT.Lock()
	defer MuFT.Unlock()
	fileTable = ft
}

// PrepareFile - Prepares random file
func PrepareFile(fileName string) (string, error) {
	if contains(NodeFiles.FileNames, fileName) {
		value := randomInt(2, 10) * 1024 * 1024
		randFile := make([]byte, value)
		if _, err := os.Stat(Name); os.IsNotExist(err) {
			os.Mkdir(Name, 0755)
		}
		err := ioutil.WriteFile(Name+"/"+fileName, randFile, 0755)
		sha := CalculateHash(Name + "/" + fileName)
		log.Println("Hash of sent file : " + sha)
		log.Println("File size of sent file : ", value/(1024*1024), "Mb")
		if err != nil {
			return "", err
		}
		return sha, nil
	}
	return "", errors.New("File not found")
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// CalculateHash - Calculates hash value of a given file
func CalculateHash(filePath string) string {
	hasher := sha256.New()
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}
	sha := hex.EncodeToString(hasher.Sum(nil))
	return sha
}

// UpdateNodeFiles - Update the file list inside node
func UpdateNodeFiles(newFile string) {
	NodeFiles.FileNames = append(NodeFiles.FileNames, newFile)
}

// UpdateFileEntryTable - Updates file table
func UpdateFileEntryTable(ip string, port string) {
	localFt := GetFT()
	if 0 == len(localFt.Files) {
		return
	}
	if 1 == len(localFt.Files) {
		localFt.Files = localFt.Files[:0]
	}
	for i, node := range GetFT().Files {
		if node.IP+":"+node.Port == ip+":"+port {
			localFt.Files = append(localFt.Files[:i], localFt.Files[i+1:]...)
			log.Println("removing entry from File table " + ip + ":" + port)
			break
		}
	}
	SetFT(localFt)
}

// AddToFileEntryTable - Updates file table
func AddToFileEntryTable(ip string, port string, fileString string) {
	localFt := GetFT()
	ftEntry := model.FileTableEntry{}
	ftEntry.IP = ip
	ftEntry.Port = port
	ftEntry.FileStrings = fileString
	localFt.Files = append(localFt.Files, ftEntry)
	log.Println("Adding entry to File table " + ip + ":" + port)
	SetFT(localFt)
}
