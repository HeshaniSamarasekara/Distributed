package main

import (
	"Distributed/P2PClient/model"
	"fmt"
	"log"
	"net/http"
	"os"

	"Distributed/P2PClient/client"
	"Distributed/P2PClient/router"
	"Distributed/P2PClient/server"
	"Distributed/P2PClient/util"
)

func main() {

	argIp := os.Args[1:]
	var arg model.Argument
	arg.IP = argIp[0]
	arg.Port = argIp[1]
	util.SetCommandLineArgument(arg)
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
