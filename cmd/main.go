package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"os"
	"os/signal"
	"time"
	"wildberries_test_task/internal/handler"
	"wildberries_test_task/internal/service"
	"wildberries_test_task/internal/storage"

	"fmt"
	"log"
	"net/http"
	"wildberries_test_task/internal/config"
)

type Handler interface {
	Method() string
	Path() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func registerHandler(router chi.Router, handler Handler) {
	router.Method(handler.Method(), handler.Path(), handler)
}

func connectionsClosedForServer(servers []*http.Server) chan struct{} {
	connectionsClosed := make(chan struct{})
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)
		defer signal.Stop(shutdown)
		<-shutdown

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		log.Println("Closing connections")
		for _, server := range servers {
			if err := server.Shutdown(ctx); err != nil {
				log.Println(err)
			}
		}
		close(connectionsClosed)
	}()
	return connectionsClosed
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	storage := storage.MemStorage{}
	service := service.NewService(&storage)

	getRouter := chi.NewRouter()
	getRouter.Use(middleware.RequestID)
	getRouter.Use(middleware.Logger)
	getRouter.Use(middleware.Recoverer)
	getRouter.Use(cors.AllowAll().Handler)

	setRouter := chi.NewRouter()
	setRouter.Use(middleware.RequestID)
	setRouter.Use(middleware.Logger)
	setRouter.Use(middleware.Recoverer)
	setRouter.Use(cors.AllowAll().Handler)

	getRouter.Group(func(router chi.Router) {
		registerHandler(router, &handler.GetUserGradeHandler{Service: service})
	})

	setRouter.Group(func(router chi.Router) {
		registerHandler(router, &handler.SetUserGradeHandler{Service: service})
	})

	addrGetPort := fmt.Sprintf(":%s", cfg.GetPort)
	addrSetPort := fmt.Sprintf(":%s", cfg.SetPort)

	servers := []*http.Server{
		{
			Addr:    addrGetPort,
			Handler: getRouter,
		},
		{
			Addr:    addrSetPort,
			Handler: setRouter,
		},
	}

	connectionsClosed := connectionsClosedForServer(servers)
	log.Println("Server with get is listening on " + addrGetPort)
	log.Println("Server with set is listening on " + addrSetPort)
	for _, server := range servers {
		server := server
		go func() {
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				log.Println(err)
			}
		}()
	}
	<-connectionsClosed
}
