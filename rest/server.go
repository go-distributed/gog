package rest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-distributed/gog/agent"
	"github.com/go-distributed/gog/config"
	log "github.com/go-distributed/gog/log"
)

const (
	listURL      = "/api/list"
	joinURL      = "/api/join"
	broadcastURL = "/api/broadcast"
	configURL    = "/api/config"
)

var (
	errInvalidMethod = errors.New("server: Invalid method")
)

// RESTServer handles RESTful requests for gog agent.
type RESTServer struct {
	ag  agent.Agent
	mux *http.ServeMux
}

// NewServer creates a new RESTful server for gog agent.
// It will also starts the agent server.
func NewServer(cfg *config.Config) *http.Server {
	ag := agent.NewAgent(cfg)
	go func() {
		if err := ag.Serve(); err != nil {
			log.Fatalf("server.NewServer(): Agent failed to serve: %v\n", err)
		}
	}()

	handler := NewRESTServer(ag)

	// Register a user message handler.
	ag.RegisterMessageHandler(handler.UserMessagHandler)

	return &http.Server{
		Addr:    cfg.RESTAddrStr,
		Handler: handler,
	}
}

// NewRESTServer creates an http.Handler to handle HTTP requests.
func NewRESTServer(ag agent.Agent) *RESTServer {
	mux := http.NewServeMux()
	rh := &RESTServer{ag, mux}
	rh.RegisterAPI(mux)
	return rh
}

// registerAPI registers the api urls.
func (rh *RESTServer) RegisterAPI(mux *http.ServeMux) {
	mux.HandleFunc(listURL, rh.List)
	mux.HandleFunc(joinURL, rh.Join)
	mux.HandleFunc(broadcastURL, rh.Broadcast)
	mux.HandleFunc(configURL, rh.Config)
	return
}

// List lists the views
func (rh *RESTServer) List(w http.ResponseWriter, r *http.Request) {
	b, err := rh.ag.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(b))
	return
}

// Join joins the agent to a cluster.
func (rh *RESTServer) Join(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, errInvalidMethod.Error(), http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	peer := r.Form.Get("peer")

	// Join a single peer
	if peer != "" {
		log.Infof("Joining: %s\n", peer)
		if err := rh.ag.Join(peer); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Infof("join, %v\n", string(b))
}

// Broadcast broadcasts the message to the cluster
func (rh *RESTServer) Broadcast(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := r.Form.Get("message")
	if msg != "" {
		log.Infof("Broadcasting: %s\n", msg)
		if err := rh.ag.Broadcast([]byte(msg)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	return
}

// Config get/set the current configuration.
func (rh *RESTServer) Config(w http.ResponseWriter, r *http.Request) {
	log.Infof("config")
}

// UserMessagHandler is the handler for user messages. It will run a script
// specified by the configuration.
func (rh *RESTServer) UserMessagHandler(msg []byte) {
	log.Infof("user message")
}

// ServeHTTP implements the http.Handler for RESTServer.
// It will get the handler from mux and invoke the handler.
func (rh *RESTServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, _ := rh.mux.Handler(r)
	h.ServeHTTP(w, r)
}
