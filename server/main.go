package main

import (
	"fmt"
	"log"
	"net"
)

var (
	listenPort = "1337"
)

func main() {
	log.Println("Listening...")

	lobby := NewLobby()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", listenPort))
	if err != nil {
		panic(err)
	}

	for {
		conn, _ := ln.Accept()
		lobby.joins <- conn
	}

	fmt.Println("Done...")
}
