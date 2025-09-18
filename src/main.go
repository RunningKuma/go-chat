package main

import (
	"github.com/RunningKuma/It-My-First-GO/src/server"

)

func main() {
	server := server.NewServer("127.0.0.1",8080)
	server.Start()
}

