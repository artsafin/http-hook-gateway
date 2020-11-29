package main

import (
	"go.uber.org/zap"
	"http-hook-gateway/internal/application"
	"log"
	"net/http"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	app := application.NewApp(logger)
	if err := app.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", app.RootHandler)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
