package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/minaria/common"
	"github.com/vahidmostofi/minaria/domain"
	"github.com/vahidmostofi/minaria/handlers"
	"github.com/vahidmostofi/minaria/repositories"
	"github.com/vahidmostofi/minaria/usecase"
)

type Server struct {
	Router      *mux.Router
	HTTPServer  http.Server
	l           *log.Logger
	bindAddress string
}

func NewServer() *Server {
	s := &Server{}
	s.l = common.GetLogger()

	s.Router = mux.NewRouter()

	// health checks
	hh := handlers.NewHealthCheck(s.l)
	hh.AttachRouter(s.Router)

	// auth handlers
	ur, _ := repositories.NewUserRepository( // TODO
		repositories.InMemoryKind,
		nil,
	)
	uc := usecase.NewUser(s.l, ur, usecase.UserOptions{})   // TODO
	ah := handlers.NewAuth(s.l, uc, domain.NewValidation()) // TODO
	ah.AttachRouter(s.Router)

	// Swagger documentations
	opts := middleware.RedocOpts{SpecURL: "/swagger.yml"}
	sh := middleware.Redoc(opts, nil)

	docsRouter := s.Router.NewRoute().Subrouter()
	docsRouter.Handle("/docs", sh).Methods(http.MethodGet)
	docsRouter.Handle("/swagger.yml", http.FileServer(http.Dir("./"))).Methods(http.MethodGet)

	s.bindAddress = viper.GetString(common.BIND_PORT)
	s.HTTPServer = http.Server{
		Addr:         s.bindAddress,     // configure the bind address
		Handler:      s.Router,          // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	return s
}

func (s *Server) Start() {
	// start the server
	go func() {
		s.l.Println("Starting server on port", s.bindAddress)

		err := s.HTTPServer.ListenAndServe()
		if err != nil {
			s.l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()
}

func (s *Server) ShutDown() {
	s.l.Println("Shutting down the server.")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.HTTPServer.Shutdown(ctx)
}
