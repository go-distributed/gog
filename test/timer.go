package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex
var hostport string

var receivedNum int
var startTime time.Time
var elaspedTime time.Duration

var vmu sync.Mutex
var aViewCnt = make(map[string]int)
var pViewCnt = make(map[string]int)

func init() {
	flag.StringVar(&hostport, "-hostport", ":11000", "The server's address")
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	startTime = time.Now()
	elaspedTime = 0
	receivedNum = 0
}

func handleClean(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	vmu.Lock()
	defer mu.Unlock()
	defer vmu.Unlock()

	startTime = time.Now()
	elaspedTime = 0
	receivedNum = 0
	aViewCnt = make(map[string]int)
	pViewCnt = make(map[string]int)
}

func handleReceived(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	receivedNum++
	elaspedTime = time.Now().Sub(startTime)
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Received: %v, time: %v\n", receivedNum, elaspedTime)

	vmu.Lock()
	defer vmu.Unlock()

	avg, std := computeView(aViewCnt)
	fmt.Fprintf(w, "Aview avg: %v, std: %v\n", avg, std)
	avg, std = computeView(pViewCnt)
	fmt.Fprintf(w, "Pview avg: %v, std: %v\n", avg, std)
}

func handleView(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// "aviewCnt:pviewCnt:id"
	view := strings.SplitAfterN(string(b), ":", 3)

	av, err := strconv.Atoi(view[0][:len(view[0])-1])
	if err != nil {
		fmt.Println(err)
		return
	}
	pv, err := strconv.Atoi(view[1][:len(view[1])-1])
	if err != nil {
		fmt.Println(err)
		return
	}

	vmu.Lock()
	defer vmu.Unlock()
	aViewCnt[view[2]] = av
	pViewCnt[view[2]] = pv
}

// Return average, std
func computeView(view map[string]int) (float64, float64) {
	node := 0
	total := 0
	for _, v := range view {
		node++
		total = total + v
	}
	avg := float64(total) / float64(node)

	varstd := 0.0
	for _, v := range view {
		varstd = varstd + (float64(v)-avg)*(float64(v)-avg)
	}
	varstd = varstd / float64(node)
	std := math.Sqrt(varstd)
	return avg, std
}

func main() {
	fmt.Println("Start server...")
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/received", handleReceived)
	http.HandleFunc("/query", handleQuery)
	http.HandleFunc("/view", handleView)

	if err := http.ListenAndServe(hostport, nil); err != nil {
		fmt.Println(err)
	}
}
