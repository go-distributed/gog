package rest

import (
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

// RESTServer handles RESTful requests for gog agent.
type RESTServer struct {
	ag  agent.Agent
	mux *http.ServeMux
}

// NewServer creates a new RESTful server for gog agent.
func NewServer(cfg *config.Config) *http.Server {
	ag := agent.NewAgent(cfg)
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
	log.Infof("list")
}

// Join joins the agent to a cluster.
func (rh *RESTServer) Join(w http.ResponseWriter, r *http.Request) {
	log.Infof("join")
}

// Broadcast broadcasts the message to the cluster
func (rh *RESTServer) Broadcast(w http.ResponseWriter, r *http.Request) {
	log.Infof("broadcast")
}

// Config get/set the curreng configuration.
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
