package internal

import (
	"github.com/getclasslabs/go-chat/internal/handlers"
	"github.com/getclasslabs/go-tools/pkg/request"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	Router *mux.Router
}

func NewServer() *Server {
	r := mux.NewRouter()
	s := Server{r}

	s.serve()

	return &s
}

func (s *Server) serve() {

	s.Router.Path("/heartbeat").HandlerFunc(request.PreRequest(handlers.Heartbeat)).Methods(http.MethodGet)
	s.Router.Path("/create").HandlerFunc(request.PreRequest(handlers.CreateRoom)).Methods(http.MethodPost)
	s.Router.Path("/connect/{room}").HandlerFunc(request.PreRequest(handlers.Connect))
}
