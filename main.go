package main

import (
	"AI-Dietitian/handler"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//go:embed public
var FS embed.FS

func main() {

	if err := initEverything(); err != nil {
		log.Fatal(err)
	}
	router := chi.NewMux()
	router.Use(handler.WithUser)

	// to handle static files
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))	

	router.Get("/", handler.MakeHandler(handler.HandleHomeIndex))
	router.Get("/chat", handler.MakeHandler(handler.HandleChatIndex))
	router.Post("/chat", handler.MakeHandler(handler.HandleChatCreate))

	port := os.Getenv("HTTP_LISTEN_ADDRESS")
	slog.Info("Application is running on", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	return godotenv.Load()
}


