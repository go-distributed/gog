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
	"sync/atomic"
	"time"
)

var hostport string

var receivedNum int32
var startTime time.Time
var elaspedTime time.Duration

var vmu sync.Mutex
var aViewCnt = make(map[string]int)
var pViewCnt = make(map[string]int)
var times = make([]time.Duration, 21)

var testSerf bool

func handleStart(w http.ResponseWriter, r *http.Request) {
	startTime = time.Now()
	elaspedTime = 0
	receivedNum = 0
	times = make([]time.Duration, 21)
	aViewCnt = make(map[string]int)
	pViewCnt = make(map[string]int)
}

func handleReceived(w http.ResponseWriter, r *http.Request) {
	reqstart := time.Now()
	num := atomic.AddInt32(&receivedNum, 1)
	if num == 1 {
		startTime = time.Now()
	}
	elaspedTime = time.Now().Sub(startTime)
	if !testSerf {
		if num%100 == 0 {
			times[num/100] = elaspedTime
		}
	} else {
		if num%20 == 0 {
			times[num/20] = elaspedTime
		}
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	fmt.Fprintf(w, "hello %d, time: %v, req time: %v", num, elaspedTime, time.Now().Sub(reqstart))
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	if !testSerf {
		for i := range times {
			fmt.Fprintf(w, "Received: %d, time: %v\n", i*100, times[i])
		}
	} else {
		for i := range times {
			fmt.Fprintf(w, "Received: %d, time: %v\n", i*20, times[i])
		}
	}
	fmt.Fprintf(w, "total received: %d, elasped time: %v\n", receivedNum, elaspedTime)

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

func init() {
	flag.StringVar(&hostport, "hostport", ":11000", "The server's address")
	flag.BoolVar(&testSerf, "testserf", false, "If testing serf")
}

func main() {
	flag.Parse()

	fmt.Println("Start server...")
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/received", handleReceived)
	http.HandleFunc("/query", handleQuery)
	http.HandleFunc("/view", handleView)

	if err := http.ListenAndServe(hostport, nil); err != nil {
		fmt.Println(err)
	}
}
