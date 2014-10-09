package main

import (
	"flag"
	"fmt"

	"github.com/go-distributed/gog/agent"
	"github.com/go-distributed/gog/config"
)

func main() {
	var peer string
	flag.StringVar(&peer, "join_node", "", "contact node")
	config, err := config.ParseConfig()
	if err != nil {
		fmt.Println("Failed to parse configuration", err)
		return
	}
	ag := agent.NewAgent(config)
	if peer != "" {
		fmt.Println("ready to join", peer)
		ag.Join(peer)
	}
	fmt.Printf("serving at %v...\n", config.AddrStr)
	ag.Serve()
}
