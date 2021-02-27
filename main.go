package main

import (
	"fmt"
	"log"

	"github.com/michelaquino/golang_api_skeleton/cmd"
)

// Version var is used to retrieve from binary git version of application
var Version = "development"

func main() {
	fmt.Println("Version:\t", Version)
	if err := cmd.Execute(); err != nil {
		log.Fatalf("cannot start application, an error has ocurred %s", err.Error())
	}
}
