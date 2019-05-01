package main

import (
	"./p3"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	router := p3.NewRouter()
	// args: 1. port {id = port number}
	if len(os.Args) > 1 {
		fmt.Println("Server listening on ", os.Args[1])
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		fmt.Println("Server listening on 6688")
		log.Fatal(http.ListenAndServe(":6688", router))
	}
}
