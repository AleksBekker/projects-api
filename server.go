package main

import (
	"context"
	"log"
	"net/http"

	"github.com/AleksBekker/project-api/database"
)

type Server struct {
	addr   string
	logger *log.Logger
	server *http.Server
	db     *db.Database
}

func NewServer(addr string, logger *log.Logger, db *db.Database) *Server {
	if logger == nil {
		logger = log.Default()
	}

	return &Server{
		addr:   addr,
		logger: logger,
		db:     db,
	}
}

func (s *Server) Run() error {

	router := http.NewServeMux()

	router.Handle("GET /projects", handleGetProjects(s.db, s.logger))

	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}
	s.server = &server

	s.logger.Printf("server started: %s\n", s.addr)

	return server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Printf("server shutdown failed:%+s\n", err)
		return err
	}

	s.logger.Println("server shutdown succeeded")
	return nil
}
