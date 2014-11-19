package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var mu sync.Mutex
var hostport string

var receivedNum int
var startTime time.Time
var elaspedTime time.Duration

func init() {
	flag.StringVar(&hostport, "-hostport", ":8080", "The server's address")
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	startTime = time.Now()
	elaspedTime = 0
	receivedNum = 0
}

func handleReceived(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	receivedNum++
	elaspedTime = time.Now().Sub(startTime)
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Received: %v, time: %v\n", receivedNum, elaspedTime)
}

func main() {
	fmt.Println("Start server...")
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/received", handleReceived)
	http.HandleFunc("/query", handleQuery)

	if err := http.ListenAndServe(hostport, nil); err != nil {
		fmt.Println(err)
	}
}
