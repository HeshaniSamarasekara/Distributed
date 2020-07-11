package router

import (
	"encoding/json"
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
	router.HandleFunc("/fileTable", GetFileTable).Methods("GET")
	return router
}

// GetFileList : Returns the file list in the node
func GetFileList(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(util.NodeFiles.FileNames)
}

// RegisterNode : Register node in network
func RegisterNode(w http.ResponseWriter, r *http.Request) {
	err := client.Register(util.IP, util.Port, util.Name)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully registered in network...\n"))
	}
}

// UnregisterNode : Unregister node in network
func UnregisterNode(w http.ResponseWriter, r *http.Request) {
	err := client.Unregister(util.IP, util.Port, util.Name)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully unregistered in network...\n"))
	}
}

// GetRouteTable - Returns the route table
func GetRouteTable(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(util.GetRT())
}

// SearchFile - Search for a file in the network
func SearchFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp, err := client.Search(vars["file_name"], "", 9999)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte(resp))
	}
}

// GetFileTable : Get the file table for node
func GetFileTable(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(util.GetFT())
}
