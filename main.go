package main

import (
	"fmt"
	"log"
	"os"
	"workspace/shop/server"
	"workspace/shop/utilities"
)

func main() {

	fmt.Println("Hello, World")

	// set some global settings
	// set the random seed
	utilities.SetRandomSeed()

	// create the log file
	logfileHandle, err := os.OpenFile("logs/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("%s Fatal: Fatal Error Signal")
	}
	defer logfileHandle.Close()
	log.SetOutput(logfileHandle)

	// start the http server for incoming requests
	server.Start()
}
