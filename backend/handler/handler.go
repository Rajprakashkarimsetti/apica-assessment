package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Rajprakashkarimsetti/apica-project/models"
	"github.com/Rajprakashkarimsetti/apica-project/service"
)

type handler struct {
	s service.LruCacher
}

func New(lruCacher service.LruCacher) handler {
	return handler{s: lruCacher}
}

// Get handles HTTP GET requests for retrieving a value associated with a key. It responds with errors if the key is missing or no value is found.
func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		models.SetError(w, http.StatusBadRequest, "Parameter key is required")

		return
	}

	value := h.s.Get(key)
	if value == "" {
		models.SetError(w, http.StatusNotFound, "No key found")

		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(value)
}

// Set handles HTTP POST/PUT requests to set a key-value pair. It responds with success or error messages based on the request body validity.
func (h handler) Set(w http.ResponseWriter, r *http.Request) {
	cache, err := io.ReadAll(r.Body)
	if err != nil {
		models.SetError(w, http.StatusBadRequest, "Invalid request body")

		return
	}

	var reqData models.CacheData

	err = json.Unmarshal(cache, &reqData)
	if err != nil {
		models.SetError(w, http.StatusBadRequest, "Invalid Body")

		return
	}

	if reqData.Key == "" {
		models.SetError(w, http.StatusBadRequest, "Key is required")

		return
	}

	if reqData.Value == "" {
		models.SetError(w, http.StatusBadRequest, "Value is required")

		return
	}

	h.s.Set(reqData.Key, reqData.Value)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode("Successfully inserted")
}
