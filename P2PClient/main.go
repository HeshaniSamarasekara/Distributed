package main

import (
	"fmt"
	"log"
	"net/http"

	"Distributed/P2PClient/client"
	"Distributed/P2PClient/router"
	"Distributed/P2PClient/server"
	"Distributed/P2PClient/util"
)

func main() {

	// The server is run on a go routine so as not to block.
	go func() {
		client.CreateConnection()
		server.CreateServer()
	}()

	// Create a server listening on port 8000
	s := &http.Server{
		Addr:    ":" + util.Props.MustGetString("serverport"),
		Handler: router.NewRouter(),
	}

	fmt.Println("Starting TCP client at port " + util.Props.MustGetString("serverport"))

	// closeConnection()

	// Continue to process new requests until an error occurs
	log.Fatal(s.ListenAndServe())
}
