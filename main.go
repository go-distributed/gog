package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/go-distributed/gog/agent"
	"github.com/go-distributed/gog/config"
)

func main() {
	config, err := config.ParseConfig()
	if err != nil {
		fmt.Println("Failed to parse configuration", err)
		return
	}
	ag := agent.NewAgent(config)
	ag.RegisterMessageHandler(msgCallBack)
	fmt.Printf("serving at %v...\n", config.AddrStr)
	go ag.Serve()

	if config.Peers != nil {
		if err := ag.Join(config.Peers); err != nil {
			fmt.Println("No available peers")
			//return
		}
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
