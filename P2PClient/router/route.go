package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/magiconair/properties"

	"Distributed/P2PClient/client"
)

var props *properties.Properties

func init() {
	client.CreateConnection()
	// server.Server()
}

// NewRouter : Creates a new router
func NewRouter(prop *properties.Properties) *mux.Router {
	props = prop
	router := mux.NewRouter()
	router.HandleFunc("/files", GetFileList).Methods("GET")
	router.HandleFunc("/register", RegisterNode).Methods("POST")
	router.HandleFunc("/unregister", UnregisterNode).Methods("DELETE")
	router.HandleFunc("/join", JoinNode).Methods("POST")
	return router
}

// GetFileList : Returns the file list in the node
func GetFileList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I have 10 files")
	w.Write([]byte("I have 10 files"))
}

// RegisterNode : Register node in network
func RegisterNode(w http.ResponseWriter, r *http.Request) {
	err := client.Register(props.MustGetString("ip"), props.MustGetString("port"), props.MustGetString("username"))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully registered in network."))
	}
}

// UnregisterNode : Register node in network
func UnregisterNode(w http.ResponseWriter, r *http.Request) {
	err := client.Unregister(props.MustGetString("ip"), props.MustGetString("port"), props.MustGetString("username"))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully unregistered in network."))
	}
}

// JoinNode : Join to a node in network
func JoinNode(w http.ResponseWriter, r *http.Request) {
	err := client.Join(props.MustGetString("ip"), props.MustGetString("port"))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully joined to network."))
	}
}
