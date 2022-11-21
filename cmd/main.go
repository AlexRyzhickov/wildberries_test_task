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

func connectionsClosedForServer(server *http.Server) chan struct{} {
	connectionsClosed := make(chan struct{})
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)
		defer signal.Stop(shutdown)
		<-shutdown

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		log.Println("Closing connections")
		if err := server.Shutdown(ctx); err != nil {
			log.Println(err)
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

	s := service.NewService()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)

	router.Group(func(router chi.Router) {
		registerHandler(router, &handler.SetUserGradeHandler{Service: s})
		registerHandler(router, &handler.GetUserGradeHandler{Service: s})
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	connectionsClosed := connectionsClosedForServer(&server)
	log.Println("Server is listening on " + addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Println(err)
	}
	<-connectionsClosed
}
