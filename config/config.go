package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
)

// Config describes the config of the system.
type Config struct {
	// Net should be tcp4 or tcp6.
	Net string
	// AddrStr is the local address string.
	AddrStr string
	// Peers is peer list.
	Peers []string
	// LocalTCPAddr is TCP address parsed from
	// Net and AddrStr.
	LocalTCPAddr *net.TCPAddr
	// AViewMinSize is the minimum size of the active view.
	AViewMinSize int
	// AViewMaxSize is the maximum size of the active view.
	AViewMaxSize int
	// PViewSize is the size of the passive view.
	PViewSize int
	// Ka is the number of nodes to choose from active view
	// when shuffling views.
	Ka int
	// Kp is the number of nodes to choose from passive view
	// when shuffling views.
	Kp int
	// Active Random Walk Length.
	ARWL int
	// Passive Random Walk Length.
	PRWL int
	// Shuffle Random Walk Length.
	SRWL int
	// Message life.
	MLife int
	// Shuffle Duration in seconds.
	ShuffleDuration int
	// Heal Duration in seconds.
	HealDuration int
	// The REST server address.
	RESTAddrStr string
	// The path to user message handler(script).
	UserMsgHandler string
}

func ParseConfig() (*Config, error) {
	var peerStr string
	var peerFile string

	cfg := new(Config)

	flag.StringVar(&cfg.Net, "net", "tcp", "The network protocol")
	flag.StringVar(&cfg.AddrStr, "addr", "localhost:8424", "The address the agent listens on")

	flag.StringVar(&peerFile, "peer-file", "", "Peer list file")
	flag.StringVar(&peerStr, "peers", "", "Comma-separated list of peers")

	flag.IntVar(&cfg.AViewMinSize, "min-aview-size", 3, "The minimum size of the active view")
	flag.IntVar(&cfg.AViewMaxSize, "max-aview-size", 5, "The maximum size of the active view")
	flag.IntVar(&cfg.PViewSize, "pview-size", 5, "The size of the passive view")

	flag.IntVar(&cfg.Ka, "ka", 1, "The number of active nodes to shuffle")
	flag.IntVar(&cfg.Kp, "kp", 3, "The number of passive nodes to shuffle")

	flag.IntVar(&cfg.ARWL, "arwl", 5, "The active random walk length")
	flag.IntVar(&cfg.PRWL, "prwl", 5, "The passive random walk length")
	flag.IntVar(&cfg.SRWL, "srwl", 3, "The shuffle random walk length")

	flag.IntVar(&cfg.MLife, "msg-life", 5000, "The default message life (milliconds)")
	flag.IntVar(&cfg.ShuffleDuration, "shuffle-duration", 5, "The default shuffle duration (seconds)")
	flag.IntVar(&cfg.HealDuration, "heal", 5, "The default heal duration (seconds)")
	flag.StringVar(&cfg.RESTAddrStr, "rest-addr", "localhost:8425", "The address of the REST server")
	flag.StringVar(&cfg.UserMsgHandler, "user-message-handler", "", "The path to the user message handler script")

	flag.Parse()

	// Check configuration.
	if peerStr != "" {
		cfg.Peers = strings.Split(peerStr, ",")
	}
	if peerFile != "" {
		peers, err := parsePeerFile(peerFile)
		if err != nil {
			return nil, err
		}
		cfg.Peers = peers
	}

	// Check agent server address.
	tcpAddr, err := net.ResolveTCPAddr(cfg.Net, cfg.AddrStr)
	if err != nil {
		return nil, err
	}
	cfg.LocalTCPAddr = tcpAddr

	// Check REST API address.
	_, err = net.ResolveTCPAddr(cfg.Net, cfg.RESTAddrStr)
	if err != nil {
		return nil, err
	}

	// Check User Message Handler.
	_, err = exec.LookPath(cfg.UserMsgHandler)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func parsePeerFile(path string) ([]string, error) {
	var peers []string
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &peers); err != nil {
		return nil, err
	}
	return peers, nil
}

func (cfg *Config) ShufflePeers() []string {
	shuffledPeers := make([]string, len(cfg.Peers))
	copy(shuffledPeers, cfg.Peers)
	for i := range shuffledPeers {
		if i == 0 {
			continue
		}
		swapIndex := rand.Intn(i)
		shuffledPeers[i], shuffledPeers[swapIndex] = shuffledPeers[swapIndex], shuffledPeers[i]
	}
	return shuffledPeers
}
