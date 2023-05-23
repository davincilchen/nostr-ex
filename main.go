package main

import (
	"nostr-ex/pkg/app/server"
	"nostr-ex/pkg/config"
)

const confPath = "./config.json"

func main() {

	//cfg, err := config.New(confPath)
	cfg, _ := config.New(confPath)

	svr := server.New(cfg)
	svr.Serve()

}
