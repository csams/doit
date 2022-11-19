package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	CompletedConfig

	Handler http.Handler
}

type preparedServer struct {
	*Server
}

func New(c CompletedConfig, handler http.Handler) (*Server, error) {
	return &Server{
		CompletedConfig: c,
		Handler:         handler,
	}, nil
}

func (s *Server) PrepareRun() preparedServer {
	return preparedServer{s}
}

func (s preparedServer) Run() error {
	if s.SecureServing {
		return http.ListenAndServeTLS(s.Address, s.CertFile, s.KeyFile, s.Handler)
	}
	fmt.Printf("Listening on %s\n", s.Address)
	return http.ListenAndServe(s.Address, s.Handler)
}
