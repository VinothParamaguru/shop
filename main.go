package main

import (
	"fmt"
	"shop/server"
	"shop/utilities"
)

func main() {

	// application init
	fmt.Println("Hello, World, new!")

	// set some global settings
	// set the random seed
	utilities.SetRandomSeed()

	// start the http server for incoming requests
	server.Start()
}
