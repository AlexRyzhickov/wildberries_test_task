package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"wildberries_test_task/internal/config"
	"wildberries_test_task/internal/handler"
	"wildberries_test_task/internal/service"
	"wildberries_test_task/internal/storage"
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

func createRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)
	return router
}

func createRouterWithBasicAuth() *chi.Mux {
	creds := map[string]string{
		"alex": "1234",
		"mike": "4321",
	}
	router := createRouter()
	router.Use(middleware.BasicAuth("Give username and password", creds))
	return router
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	storage := storage.InitializeMemoryStorage()
	service := service.NewService(storage, nc, cfg.Priority)

	getRouter := createRouter()
	setRouter := createRouterWithBasicAuth()

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
