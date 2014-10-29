package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/go-distributed/gog/agent"
	"github.com/go-distributed/gog/config"
)

func main() {
	var peer string
	flag.StringVar(&peer, "peer_addr", "", "contact node")
	config, err := config.ParseConfig()
	if err != nil {
		fmt.Println("Failed to parse configuration", err)
		return
	}
	ag := agent.NewAgent(config)
	ag.RegisterMessageHandler(msgCallBack)
	fmt.Printf("serving at %v...\n", config.AddrStr)
	go ag.Serve()
	if peer != "" {
		fmt.Println("ready to join", peer)
		ag.Join(peer)
	}

	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("input a message:")
		bs, err := input.ReadString('\n')
		if err != nil {
			fmt.Println("error reading:", err)
			break
		}
		if bs == "list\n" {
			ag.List()
			continue
		}
		ag.Broadcast([]byte(bs))
	}
}

func msgCallBack(msg []byte) {
	fmt.Println(string(msg))
}
