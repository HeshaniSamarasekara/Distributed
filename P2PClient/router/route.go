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
	router.HandleFunc("/routeTable", GetRouteTable).Methods("GET")
	router.HandleFunc("/search/{file_name}", SearchFile).Methods("GET")
	return router
}

// GetFileList : Returns the file list in the node
func GetFileList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I have 10 files")
	w.Write([]byte("I have 10 files"))
}

// RegisterNode : Register node in network
func RegisterNode(w http.ResponseWriter, r *http.Request) {
	err := client.Register(util.GetCommandLineArgument().IP, util.GetCommandLineArgument().Port, util.Props.MustGetString("username"))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully registered in network."))
	}
}

// UnregisterNode : Register node in network
func UnregisterNode(w http.ResponseWriter, r *http.Request) {
	err := client.Unregister(util.GetCommandLineArgument().IP, util.GetCommandLineArgument().Port, util.Props.MustGetString("username"))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully unregistered in network."))
	}
}

// GetRouteTable - Returns the route table
func GetRouteTable(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(util.RouteTable)
}

// SearchFile - Search for a file in the network
func SearchFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := client.Search(vars["file_name"])
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("File is in network."))
	}
}
