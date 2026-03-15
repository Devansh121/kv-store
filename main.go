package main

import (
	"flag"
	"log"

	"github.com/Devansh121/kv-store/config"
	"github.com/Devansh121/kv-store/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host for the KV store server")
	flag.IntVar(&config.Port, "port", 7379, "Port for the KV store server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Starting KV store server...")
	server.RunSyncTCPServer()
}
