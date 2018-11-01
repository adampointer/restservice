package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adampointer/restservice/data"
	"github.com/adampointer/restservice/handlers"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func routes(h handlers.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/payments", h.GetAll).Methods("GET")
	router.HandleFunc("/payments/{id}", h.GetOne).Methods("GET")
	router.HandleFunc("/payments/{id}", h.Create).Methods("PUT")
	router.HandleFunc("/payments/{id}", h.Update).Methods("POST")
	router.HandleFunc("/payments/{id}", h.Delete).Methods("DELETE")
	return router
}

func listenTCP(stop chan bool, errs chan error) {
	// The db path and listen address should really be from flags/config
	dbClient, err := data.NewClient("data.db")
	if err != nil {
		log.Fatalf("unable to initialise database: %s", err)
	}
	defer dbClient.Close()
	paymentsHandler := handlers.NewPayments(dbClient)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(paymentsHandler),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// Send serve errors to our errors channel
			errs <- err
		}
	}()
	<-stop
	log.Debug("starting shutdown of serve goroutine")
	ctx, cancel := context.WithTimeout(context.Background(), (500 * time.Millisecond))
	defer cancel()
	// Allow existing connections 500ms to drain, if there are any
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down http server: %s", err)
	}
}

func main() {
	// Trap kill signals
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	// Listen in a new goroutine
	log.Info("starting HTTP server on :8080")
	stop := make(chan bool)
	errs := make(chan error)
	go listenTCP(stop, errs)
	select {
	case err := <-errs:
		log.Errorf("serve error: %s", err)
	case <-terminate:
		log.Info("starting 1 second timeout for graceful shutdowns")
		close(stop)
		time.Sleep(time.Second)
		break
	}
	log.Info("terminating")
}
