package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"shopscraper/pkg/database"
	"shopscraper/pkg/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func TestMain(m *testing.M) {
	os.Setenv("SHOPSCRAPER_API_KEY", "test-api-key")
	runID := generateRandomString(16)
	productTableName := "test_products_" + runID

	db = database.NewPostgresDB()
	db.Initialize("postgresql://test:Test1234@localhost:5432/test?sslmode=disable", productTableName)
	defer db.Close()

	code := m.Run()
	os.Exit(code)
}

func setup(t *testing.T) *[]models.Product {
	products := []models.Product{
		{Name: "Product 1", Shop: "Shop 1", Price: 10, Link: "https://example.com/product1", LastSeen: time.Now().UTC(), Notified: true},
		{Name: "Product 2", Shop: "Shop 2", Price: 19, Link: "https://example.com/product2", LastSeen: time.Now().UTC(), Notified: true},
		{Name: "Product 3", Shop: "Shop 3", Price: 5, Link: "https://example.com/product3", LastSeen: time.Now().UTC(), Notified: false},
	}

	err := db.EnsureProductTableExists()
	assert.NoError(t, err, "Ensuring product table exists should not produce an error")

	_, err = db.SaveProducts(products)
	assert.NoError(t, err, "Saving products should not produce an error")

	return &products
}

func teardown(t *testing.T) {
	err := db.DropProductTable()
	assert.NoError(t, err, "Dropping product table should not produce an error")
}

func TestApiKeyMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testCases := []struct {
		apiKey         string
		expectedStatus int
		requestWithKey string
	}{
		{"valid-api-key", http.StatusOK, "valid-api-key"},
		{"valid-api-key", http.StatusUnauthorized, "invalid-api-key"},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Add("X-API-KEY", tc.requestWithKey)
		rr := httptest.NewRecorder()

		handler := apiKeyMiddleware(nextHandler, tc.apiKey)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, tc.expectedStatus, rr.Code, "Handler returned wrong status code")
	}
}

func TestGetProductsEndpoint(t *testing.T) {
	expectedProducts := setup(t)
	defer teardown(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			preflightHandler(w, r)
		} else {
			apiKeyMiddleware(http.HandlerFunc(getProducts), "test-api-key").ServeHTTP(w, r)
		}
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/products", nil)
	assert.NoError(t, err, "Creating request should not produce an error")
	req.Header.Add("X-API-KEY", os.Getenv("SHOPSCRAPER_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Executing request should not produce an error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	var returnedProducts []models.Product
	err = json.NewDecoder(resp.Body).Decode(&returnedProducts)
	assert.NoError(t, err, "Decoding response should not produce an error")

	assert.Equal(t, len(*expectedProducts), len(returnedProducts), "Number of returned products should match")

	for i, expected := range *expectedProducts {
		assert.Equal(t, expected.Name, returnedProducts[i].Name, "Product names should match")
		assert.Equal(t, expected.Shop, returnedProducts[i].Shop, "Product shops should match")
		assert.Equal(t, expected.Price, returnedProducts[i].Price, "Product prices should match")
		assert.Equal(t, expected.Link, returnedProducts[i].Link, "Product links should match")
		assert.Equal(t, expected.Notified, returnedProducts[i].Notified, "Product notification statuses should match")

		expectedRounded := expected.LastSeen.UTC().Round(time.Millisecond)
		retrievedRounded := returnedProducts[i].LastSeen.UTC().Round(time.Millisecond)
		assert.Equal(t, expectedRounded, retrievedRounded, "Product last seen timestamps should match")
	}
}
