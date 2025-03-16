package main

import (
	"HTTP/infoServer"
	"HTTP/middleware"
	"HTTP/weather/handlers"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
	}
	port := os.Getenv("PORT")
	addr := ":" + port

	router := http.NewServeMux()
	stack := middleware.CreateStack(middleware.Logging)
	server := http.Server{
		Addr:    addr,
		Handler: stack(router),
	}

	router.HandleFunc("/weather/open/{city}", handlers.HandleGetWeatherAPIOpenWeather)
	router.HandleFunc("/weather/mail/{city}", handlers.HandleGetWeatherParseMailRu)
	router.HandleFunc("GET /infoserver", infoServer.HandleGetInfoServer)

	fmt.Printf("сервер стартует на %s порту \n", addr)
	server.ListenAndServe()

}

//TODO: завернкть сервис в докер
