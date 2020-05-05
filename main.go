package main

import (
	"flag"
	"log"

	"github.com/wear_underpants/client"
	"github.com/wear_underpants/server"
)

var (
	port           = "8082"
	h              = false
	wearServerAddr = "127.0.0.1:8082"
	mode           = "server"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&port, "p", "8082", "port")
	flag.StringVar(&mode, "m", "server", "server")
	flag.StringVar(&wearServerAddr, "addr", "127.0.0.1:8083", "remote addr:port")
	flag.Parse()
}

func main() {
	if mode == "server" {
		server.StartServer(port)
	} else if mode == "client" {
		client.StartClient(port, wearServerAddr)
	}
}
