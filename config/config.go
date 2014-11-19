package config

import (
	"bufio"
	"flag"
	"io"
	"math/rand"
	"net"
	"os"
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
	// Network Delay Max in milliseconds.
	NetworkDelayMax int
	// Network Drop Rate (0-1)
	NetworkDropRate int
}

func ParseConfig() (*Config, error) {
	var peerStr string
	var peerFile string

	cfg := new(Config)

	flag.StringVar(&cfg.Net, "net", "tcp", "The network protocol")
	flag.StringVar(&cfg.AddrStr, "addr", "localhost:8424", "The address the agent listens on")

	flag.StringVar(&peerFile, "peer-file", "", "Peer list file")
	flag.StringVar(&peerStr, "peers", "", "Comma-separated list of peers")

	flag.IntVar(&cfg.AViewMinSize, "min_active_view_size", 3, "The minimum size of the active view")
	flag.IntVar(&cfg.AViewMaxSize, "max_active_view_size", 5, "The maximum size of the active view")
	flag.IntVar(&cfg.PViewSize, "passive_view_size", 5, "The size of the passive view")

	flag.IntVar(&cfg.Ka, "ka", 1, "The number of active nodes to shuffle")
	flag.IntVar(&cfg.Kp, "kp", 3, "The number of passive nodes to shuffle")

	flag.IntVar(&cfg.ARWL, "arwl", 5, "The active random walk length")
	flag.IntVar(&cfg.PRWL, "prwl", 5, "The passive random walk length")
	flag.IntVar(&cfg.SRWL, "srwl", 3, "The shuffle random walk length")

	flag.IntVar(&cfg.MLife, "msg_life", 5000, "The default message life (milliconds)")
	flag.IntVar(&cfg.ShuffleDuration, "shuffle_duration", 5, "The default shuffle duration (seconds)")
	flag.IntVar(&cfg.HealDuration, "heal", 5, "The default heal duration (seconds)")
	flag.IntVar(&cfg.NetworkDelayMax, "delay", 0, "The network delay for simulation (milliseconds)")
	flag.IntVar(&cfg.NetworkDropRate, "droprate", 0, "The network drop rate")

	flag.Parse()

	// TODO check config.
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
	tcpAddr, err := net.ResolveTCPAddr(cfg.Net, cfg.AddrStr)
	if err != nil {
		return nil, err
	}
	cfg.LocalTCPAddr = tcpAddr
	return cfg, nil
}

func parsePeerFile(path string) ([]string, error) {
	peers := make([]string, 0)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	br := bufio.NewReader(f)

	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		peers = append(peers, line[:len(line)-1])
	}
	return peers, nil
}

func (cfg *Config) ShufflePeers() []string {
	shuffledPeers := make([]string, len(cfg.Peers))
	if copy(shuffledPeers, cfg.Peers) != len(cfg.Peers) {
		panic("Failed to copy")
	}
	for i := range shuffledPeers {
		if i == 0 {
			continue
		}
		swapIndex := rand.Intn(i)
		shuffledPeers[i], shuffledPeers[swapIndex] = shuffledPeers[swapIndex], shuffledPeers[i]
	}
	return shuffledPeers
}
