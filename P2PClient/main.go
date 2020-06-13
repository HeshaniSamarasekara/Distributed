package main

import (
	"log"
	"net/http"

	"Distributed/P2PClient/router"
)

func main() {
	// Create a server listening on port 8000
	s := &http.Server{
		Addr:    ":8001",
		Handler: router.NewRouter(),
	}

	// closeConnection()

	// Continue to process new requests until an error occurs
	log.Fatal(s.ListenAndServe())
}
