package main

import (
	"fmt"
	"workspace/shop/server"
	"workspace/shop/utilities"
)

func main() {
	fmt.Println("Hello, World")
	// set some global settings
	utilities.SetRandomSeed()
	server.Start()
}
