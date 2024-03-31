package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/Rajprakashkarimsetti/apica-project/models"
)

// CORS is a middleware that sets headers to allow cross-origin requests.
// It also handles preflight requests.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers to allow cross-origin requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// CorrelationIDMiddleware is a middleware that checks if the request contains a Correlation ID.
// If no Correlation ID is found in the request header, it generates a new one and adds it to the request headers.
// It then passes the Correlation ID to the next handler.
func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request contains a Correlation ID
		correlationID := r.Header.Get(models.CorrelationIDHeader)
		if correlationID == "" {
			// If no Correlation ID is found in the request header, generate a new one
			correlationID = generateCorrelationID()
			// Add the Correlation ID to the request headers
			r.Header.Set(models.CorrelationIDHeader, correlationID)
		}

		// Pass the Correlation ID to the next handler
		ctx := r.Context()
		ctx = context.WithValue(ctx, models.CorrelationIDHeader, correlationID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateCorrelationID() string {
	// Generate a unique Correlation ID
	return uuid.New().String()
}

// RequestLogger is a middleware that logs information about incoming HTTP requests.
// It logs the HTTP method and the Correlation ID header (if present) of the incoming request.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received the Request, Method: %v Correlation-ID: %v", r.Method, r.Header.Get(models.CorrelationIDHeader))
		next.ServeHTTP(w, r)
	})
}

// SetResponseHeaders is a middleware that sets the response headers.
// It sets the Content-Type header to "application/json" and copies the Correlation ID header from the request to the response.
func SetResponseHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set(models.CorrelationIDHeader, r.Header.Get(models.CorrelationIDHeader))

		next.ServeHTTP(w, r)
	})
}
