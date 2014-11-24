package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/go-distributed/gog/agent"
	"github.com/go-distributed/gog/config"
)

var cfg *config.Config
var err error

func main() {
	cfg, err = config.ParseConfig()
	if err != nil {
		fmt.Println("Failed to parse configuration", err)
		return
	}
	ag := agent.NewAgent(cfg)
	ag.RegisterMessageHandler(msgCallBack)
	fmt.Printf("serving at %v...\n", cfg.AddrStr)
	go ag.Serve()

	if cfg.Peers != nil {
		if err := ag.Join(cfg.Peers); err != nil {
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
		fmt.Println(bs)
		ag.Broadcast([]byte(bs))
	}
}

func msgCallBack(msg []byte) {
	fmt.Println(string(msg))
	resp, err := http.Get(cfg.MeasureServer + "/received")
	if err != nil {
		fmt.Println("Failed to send received", err)
		return
	}
	defer resp.Body.Close()
}
