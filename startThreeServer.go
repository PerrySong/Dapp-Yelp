package main

import (
	"fmt"
	"github.com/cs686-blockchain-p3-PerrySong/p3"
	"log"
	"net/http"
)

func main() {
	router := p3.NewRouter()

	fmt.Println("Server listening on 6688")
	go http.ListenAndServe(":6688", router)
	fmt.Println("Server listening on 6689")
	go http.ListenAndServe(":6689", router)
	fmt.Println("Server listening on 6690")
	go log.Fatal(http.ListenAndServe(":6690", router))
}
