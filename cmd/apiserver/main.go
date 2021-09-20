package main

import (
	"log"
	"pingram/pkg/apiserver"
)

func main() {
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
