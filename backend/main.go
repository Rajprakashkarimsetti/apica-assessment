package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Rajprakashkarimsetti/apica-project/cacher"
	hTTP "github.com/Rajprakashkarimsetti/apica-project/handler"
	"github.com/Rajprakashkarimsetti/apica-project/middlewares"
	cacheService "github.com/Rajprakashkarimsetti/apica-project/service"
	cacheStore "github.com/Rajprakashkarimsetti/apica-project/store"
)

func main() {
	// Initialize Cache
	c := cacher.NewCache(1024)

	// Initialize Routing
	router := mux.NewRouter()

	// Initialize Middlewares
	router.Use(middlewares.CorrelationIDMiddleware, middlewares.CORS, middlewares.RequestLogger, middlewares.SetResponseHeaders)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initialize Handler, Service, Store
	cacheStr := cacheStore.New(c)
	cacheSvc := cacheService.New(cacheStr)
	handler := hTTP.New(cacheSvc)

	router.HandleFunc("/set", handler.Set) // if data is already present it updates it
	router.HandleFunc("/get/{key}", handler.Get)

	log.Printf("Running the server at Port: %v", "8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Error while running the server, Err: %v", err)
		return
	}
}
