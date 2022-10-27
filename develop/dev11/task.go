package dev11

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const configPath = "config.json"

func main() {
	cfg, err := ParseConfig(configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	cache := NewCache()
	service := NewService(cache)
	validator := NewValidator()

	handler := NewHandler(service, validator)
	handler.Run()

	server := &http.Server{Addr: cfg.Port}
	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	logrus.WithField("port", cfg.Port).Info("server started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logrus.Info("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatal(err)
	}
}
