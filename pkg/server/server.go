package server

import (
	"net/http"

	"github.com/go-logr/logr"
)

type Server struct {
	CompletedConfig

	Handler http.Handler
	Log     logr.Logger
}

type preparedServer struct {
	*Server
}

func New(c CompletedConfig, handler http.Handler, log logr.Logger) (*Server, error) {
	return &Server{
		CompletedConfig: c,
		Handler:         handler,
		Log:             log,
	}, nil
}

func (s *Server) PrepareRun() preparedServer {
	return preparedServer{s}
}

func (s preparedServer) Run() error {
	if s.Options.SecureServing {
		return http.ListenAndServeTLS(s.Options.Address, s.Options.CertFile, s.Options.KeyFile, s.Handler)
	}

	s.Log.V(0).Info("Listening on", "address", s.Options.Address)
	return http.ListenAndServe(s.Options.Address, s.Handler)
}
