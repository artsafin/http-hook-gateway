package main

import (
	"go.uber.org/zap"
	"http-hook-gateway/internal/application"
	"log"
	"net/http"
)

func main() {
	logger, _ := zap.NewDevelopment(zap.WithCaller(false))
	defer logger.Sync()

	app := application.NewApp(logger)
	if err := app.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	addr := ":8080"
	logger.Info("Starting web server at " + addr)
	http.HandleFunc("/", app.RootHandler)
	err := http.ListenAndServe(addr, nil)
	log.Fatal(err)
}
