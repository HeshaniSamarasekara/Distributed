package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"Distributed/P2PClient/client"
)

func init() {
	client.CreateConnection()
}

// NewRouter : Creates a new router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/files", GetFileList).Methods("GET")
	router.HandleFunc("/register", RegisterNode).Methods("POST")
	router.HandleFunc("/unregister", UnregisterNode).Methods("DELETE")
	return router
}

// GetFileList : Returns the file list in the node
func GetFileList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I have 10 files")
	w.Write([]byte("I have 10 files"))
}

// RegisterNode : Register node in network
func RegisterNode(w http.ResponseWriter, r *http.Request) {
	err := client.Register("10.30.246.237", "8001", "heshani")
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully registered in network."))
		w.WriteHeader(http.StatusOK)
	}
}

// UnregisterNode : Register node in network
func UnregisterNode(w http.ResponseWriter, r *http.Request) {
	err := client.Unregister("10.30.246.237", "8001", "heshani")
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("Successfully unregistered in network."))
		w.WriteHeader(http.StatusOK)
	}
}
