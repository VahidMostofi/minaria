package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"github.com/vahidmostofi/minaria/handlers"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {
	env.Parse()

	l := log.New(os.Stdout, "minaria", log.LstdFlags)

	sm := mux.NewRouter()

	// health checks
	hh := handlers.NewHealthCheck(l)

	healthRouter := sm.PathPrefix("/health").Subrouter()
	healthRouter.HandleFunc("", hh.CheckHealthStatus).Methods(http.MethodGet)

	// Swagger documentations
	opts := middleware.RedocOpts{SpecURL: "/swagger.yml"}
	sh := middleware.Redoc(opts, nil)

	docsRouter := sm.NewRoute().Subrouter()
	docsRouter.Handle("/docs", sh).Methods(http.MethodGet)
	docsRouter.Handle("/swagger.yml", http.FileServer(http.Dir("./"))).Methods(http.MethodGet)

	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
