package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"Distributed/P2PClient/client"
	"Distributed/P2PClient/util"
)

// NewRouter : Creates a new router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/files", GetFileList).Methods("GET")
	router.HandleFunc("/register", RegisterNode).Methods("POST")
	router.HandleFunc("/unregister", UnregisterNode).Methods("DELETE")
	router.HandleFunc("/join", JoinNode).Methods("POST")
	router.HandleFunc("/routeTable", GetRouteTable).Methods("GET")
	return router
}

// GetFileList : Returns the file list in the node
func GetFileList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I have 10 files")
	w.Write([]byte("I have 10 files"))
}

// RegisterNode : Register node in network
func RegisterNode(w http.ResponseWriter, r *http.Request) {
	err := client.Register(util.Props.MustGetString("ip"), util.Props.MustGetString("port"), util.Props.MustGetString("username"))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully registered in network."))
	}
}

// UnregisterNode : Register node in network
func UnregisterNode(w http.ResponseWriter, r *http.Request) {
	err := client.Unregister(util.Props.MustGetString("ip"), util.Props.MustGetString("port"), util.Props.MustGetString("username"))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully unregistered in network."))
	}
}

// JoinNode : Join to a node in network
func JoinNode(w http.ResponseWriter, r *http.Request) {
	err := client.Join(util.Props.MustGetString("ip"), util.Props.MustGetString("port"))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully joined to network."))
	}
}

// GetRouteTable - Returns the route table
func GetRouteTable(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(util.RouteTable)
}
