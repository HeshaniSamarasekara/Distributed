package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"Distributed/P2PClient/router"

	"github.com/magiconair/properties"
)

var props *properties.Properties

func init() {
	configFile := flag.String("configFile", "application.yaml", "Configuration File")
	flag.Parse()
	props = properties.MustLoadFile(*configFile, properties.UTF8)
}

func main() {
	// Create a server listening on port 8000
	s := &http.Server{
		Addr:    ":" + props.MustGetString("port"),
		Handler: router.NewRouter(props),
	}

	fmt.Println("Starting UDP client at port " + props.MustGetString("port"))

	// closeConnection()

	// Continue to process new requests until an error occurs
	log.Fatal(s.ListenAndServe())
}
