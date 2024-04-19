package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"shopscraper/pkg/database"

	_ "github.com/lib/pq"
)

var db database.Database

func apiKeyMiddleware(next http.Handler, expectedApiKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		requestApiKey := r.Header.Get("X-API-KEY")

		if requestApiKey != expectedApiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-KEY")
}

func preflightHandler(w http.ResponseWriter, _ *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	products, err := db.GetAllProducts()
	if err != nil {
		log.Printf("Failed to retrieve products: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("Failed to encode products: %v", err)
		http.Error(w, "Error encoding products", http.StatusInternalServerError)
	}
}

func main() {
	connectionString := os.Getenv("SHOPSCRAPER_DB_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatalf("SHOPSCRAPER_DB_CONNECTION_STRING not provided")
	}
	db = database.NewPostgresDB()
	db.Initialize(connectionString, "products")
	defer db.Close()

	db.EnsureProductTableExists()

	expectedApiKey := os.Getenv("SHOPSCRAPER_API_KEY")
	if expectedApiKey == "" {
		log.Fatalf("No API key specified")
	}

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			preflightHandler(w, r)
		} else {
			apiKeyMiddleware(http.HandlerFunc(getProducts), expectedApiKey).ServeHTTP(w, r)
		}
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
