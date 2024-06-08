package database

import (
	"fmt"
	"math/rand"
	"os"
	"shopscraper/pkg/models"
	"testing"
	"time"

	_ "github.com/lib/pq"
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

var db *PostgresDB

func TestMain(m *testing.M) {
	runID := generateRandomString(16)

	db = NewPostgresDB()
	db.Initialize("postgresql://test:Test1234@localhost:5432/test?sslmode=disable", "test_products_"+runID)
	defer db.Close()

	result := m.Run()

	os.Exit(result)
}

func setup(t *testing.T) {
	err := db.EnsureProductTableExists()
	if err != nil {
		t.Errorf("failed to ensure table exists %v", err)
	}
}

func teardown(t *testing.T) {
	err := db.DropProductTable()
	if err != nil {
		t.Errorf("failed to drop table %v", err)
	}
}

func TestGetAllProducts(t *testing.T) {
	setup(t)
	defer teardown(t)

	// Insert test data
	query := fmt.Sprintf(`
	INSERT INTO %s (name, shop, price, link, first_seen, last_seen, notified)
	VALUES
		('Product 1', 'Shop 1', '10', 'https://example.com/product1', $1, $1, true),
		('Product 2', 'Shop 2', '19', 'https://example.com/product2', $1, $1, false)
	`, db.productTableName)

	_, err := db.db.Exec(query, time.Now().UTC())
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	products, err := db.GetAllProducts()
	if err != nil {
		t.Fatalf("Failed to get all products: %v", err)
	}

	expectedCount := 2
	if len(products) != expectedCount {
		t.Errorf("Expected %d products, but got %d", expectedCount, len(products))
	}

	expectedProducts := []models.Product{
		{Name: "Product 1", Shop: "Shop 1", Price: 10, Link: "https://example.com/product1", Notified: true},
		{Name: "Product 2", Shop: "Shop 2", Price: 19, Link: "https://example.com/product2", Notified: false},
	}
	for i, expected := range expectedProducts {
		assert.Equal(t, expected.Name, products[i].Name, "Product %d Name mismatch", i+1)
		assert.Equal(t, expected.Shop, products[i].Shop, "Product %d Shop mismatch", i+1)
		assert.Equal(t, expected.Price, products[i].Price, "Product %d Price mismatch", i+1)
		assert.Equal(t, expected.Link, products[i].Link, "Product %d Link mismatch", i+1)
		assert.Equal(t, expected.Notified, products[i].Notified, "Product %d Notified status mismatch", i+1)
	}
}

func TestGetNonNotifiedProducts(t *testing.T) {
	setup(t)
	defer teardown(t)

	// Insert test data
	query := fmt.Sprintf(`
		INSERT INTO %s (name, shop, price, link, first_seen, last_seen, notified)
		VALUES
			('Product 1', 'Shop 1', '10', 'https://example.com/product1', $1, $1, true),
			('Product 2', 'Shop 2', '19', 'https://example.com/product2', $1, $1, false),
			('Product 3', 'Shop 3', '5', 'https://example.com/product3', $1, $1, false)
		`, db.productTableName)

	_, err := db.db.Exec(query, time.Now().UTC())
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	products, err := db.GetNonNotifiedProducts()
	if err != nil {
		t.Fatalf("Failed to get non-notified products: %v", err)
	}

	expectedCount := 2
	if len(products) != expectedCount {
		t.Errorf("Expected %d non-notified products, but got %d", expectedCount, len(products))
	}

	expectedProducts := []models.Product{
		{Name: "Product 2", Shop: "Shop 2", Price: 19, Link: "https://example.com/product2", Notified: false},
		{Name: "Product 3", Shop: "Shop 3", Price: 5, Link: "https://example.com/product3", Notified: false},
	}
	for i, expected := range expectedProducts {
		assert.Equal(t, expected.Name, products[i].Name, "Name mismatch for product %d", i+1)
		assert.Equal(t, expected.Shop, products[i].Shop, "Shop mismatch for product %d", i+1)
		assert.Equal(t, expected.Price, products[i].Price, "Price mismatch for product %d", i+1)
		assert.Equal(t, expected.Link, products[i].Link, "Link mismatch for product %d", i+1)
		assert.Equal(t, expected.Notified, products[i].Notified, "Notified status mismatch for product %d", i+1)
	}
}

func TestSaveProducts(t *testing.T) {
	setup(t)
	defer teardown(t)

	// Test Case 1, add new products

	initalTime := time.Now().UTC()
	products := []models.Product{
		{Name: "Product 1", Shop: "Shop 1", Price: 10, Link: "https://example.com/product1", LastSeen: initalTime, Notified: true},
		{Name: "Product 2", Shop: "Shop 2", Price: 19, Link: "https://example.com/product2", LastSeen: initalTime, Notified: true},
		{Name: "Product 3", Shop: "Shop 3", Price: 5, Link: "https://example.com/product3", LastSeen: initalTime, Notified: false},
	}

	newProducts, err := db.SaveProducts(products)
	if err != nil {
		t.Fatalf("Failed to save products: %v", err)
	}

	// Assert the expected number of new products
	expectedNewCount := 3
	if len(newProducts) != expectedNewCount {
		t.Errorf("Expected %d new products, but got %d", expectedNewCount, len(newProducts))
	}

	// Assert the expected new product data
	for i, expected := range products {
		assert.Equal(t, expected.Name, newProducts[i].Name, "Product Name mismatch")
		assert.Equal(t, expected.Shop, newProducts[i].Shop, "Product Shop mismatch")
		assert.Equal(t, expected.Price, newProducts[i].Price, "Product Price mismatch")
		assert.Equal(t, expected.Link, newProducts[i].Link, "Product Link mismatch")
		expectedRounded := expected.LastSeen.UTC().Round(time.Millisecond)
		newRounded := newProducts[i].LastSeen.UTC().Round(time.Millisecond)
		assert.Equal(t, expectedRounded.String(), newRounded.String(), "Product last seen timestamps do not match")
		assert.Equal(t, expected.Notified, newProducts[i].Notified, "Product Notified status mismatch")
	}

	// Test Case 2, Product with a price update
	newTime := time.Now().UTC()

	// Update an existing product price
	updatedProduct := models.Product{Name: "Product 1", Shop: "Shop 1", Price: 12, Link: "https://example.com/product1", LastSeen: newTime, Notified: true}
	updatedProducts := []models.Product{updatedProduct}

	newProducts, err = db.SaveProducts(updatedProducts)
	if err != nil {
		t.Fatalf("Failed to save updated product: %v", err)
	}

	// Assert that no new products were returned
	if len(newProducts) != 0 {
		t.Errorf("Expected to see no new products, but got %d", len(newProducts))
	}

	// Retrieve the updated product from the database
	var retrievedProduct models.Product
	err = db.db.QueryRow("SELECT name, shop, previous_price, price, link, notified, first_seen, last_seen FROM "+db.productTableName+" WHERE name = $1", updatedProduct.Name).
		Scan(&retrievedProduct.Name, &retrievedProduct.Shop, &retrievedProduct.PreviousPrice, &retrievedProduct.Price, &retrievedProduct.Link, &retrievedProduct.Notified, &retrievedProduct.FirstSeen, &retrievedProduct.LastSeen)
	if err != nil {
		t.Fatalf("Failed to retrieve updated product: %v", err)
	}

	// Assert the updated product data
	assert.Equal(t, updatedProduct.Name, retrievedProduct.Name, "Name mismatch")
	assert.Equal(t, updatedProduct.Shop, retrievedProduct.Shop, "Shop mismatch")
	assert.Equal(t, updatedProduct.Price, retrievedProduct.Price, "Price mismatch")
	if assert.True(t, retrievedProduct.PreviousPrice.Valid, "Expected PreviousPrice to be valid") {
		assert.Equal(t, retrievedProduct.PreviousPrice.Int64, int64(10), "Price mismatch")
	}
	assert.Equal(t, updatedProduct.Link, retrievedProduct.Link, "Link mismatch")

	updatedRounded := updatedProduct.LastSeen.UTC().Round(time.Millisecond)
	retrievedRounded := retrievedProduct.LastSeen.UTC().Round(time.Millisecond)
	assert.Equal(t, updatedRounded.String(), retrievedRounded.String(), "Product last seen timestamps do not match")

	retrievedRoundedFirst := retrievedProduct.FirstSeen.UTC().Round(time.Millisecond)
	// First seen should not be equal to Last Seen
	assert.NotEqual(t, retrievedRounded, retrievedRoundedFirst)
	initalTimeRounded := initalTime.UTC().Round(time.Millisecond)
	// First seen should be from the first insert
	assert.Equal(t, initalTimeRounded, retrievedRoundedFirst)

	// Notified should be reset to false when product changes price
	assert.Equal(t, false, retrievedProduct.Notified, "Notified status mismatch")

	// Test Case 3, Previous Price stays if price doesn't change
	updatedProduct = models.Product{Name: "Product 1", Shop: "Shop 1", Price: 12, Link: "https://example.com/product1", LastSeen: newTime, Notified: true}
	updatedProducts = []models.Product{updatedProduct}

	newProducts, err = db.SaveProducts(updatedProducts)
	if err != nil {
		t.Fatalf("Failed to save updated product: %v", err)
	}

	// Assert that no new products were returned
	if len(newProducts) != 0 {
		t.Errorf("Expected to see no new products, but got %d", len(newProducts))
	}

	// Retrieve the updated product from the database
	err = db.db.QueryRow("SELECT name, shop, previous_price, price, link, notified, last_seen FROM "+db.productTableName+" WHERE name = $1", updatedProduct.Name).
		Scan(&retrievedProduct.Name, &retrievedProduct.Shop, &retrievedProduct.PreviousPrice, &retrievedProduct.Price, &retrievedProduct.Link, &retrievedProduct.Notified, &retrievedProduct.LastSeen)
	if err != nil {
		t.Fatalf("Failed to retrieve updated product: %v", err)
	}

	// Assert the updated product data
	assert.Equal(t, updatedProduct.Name, retrievedProduct.Name, "Name mismatch")
	assert.Equal(t, updatedProduct.Shop, retrievedProduct.Shop, "Shop mismatch")
	assert.Equal(t, updatedProduct.Price, retrievedProduct.Price, "Price mismatch")
	if assert.True(t, retrievedProduct.PreviousPrice.Valid, "Expected PreviousPrice to be valid") {
		assert.Equal(t, retrievedProduct.PreviousPrice.Int64, int64(10), "Price mismatch")
	}

	assert.Equal(t, updatedProduct.Link, retrievedProduct.Link, "Link mismatch")
	updatedRounded = updatedProduct.LastSeen.UTC().Round(time.Millisecond)
	retrievedRounded = retrievedProduct.LastSeen.UTC().Round(time.Millisecond)
	assert.Equal(t, updatedRounded.String(), retrievedRounded.String(), "Product last seen timestamps do not match")
	// Notified should be reset to false when product changes price
	assert.Equal(t, false, retrievedProduct.Notified, "Notified status mismatch")

	// Test Case 4, Product with no actual changes - Notified should stay True
	newTime = time.Now().UTC()
	updatedProduct = models.Product{Name: "Product 2", Shop: "Shop 2", Price: 19, Link: "https://example.com/product2", LastSeen: newTime, Notified: false}
	updatedProducts = []models.Product{updatedProduct}

	// Call the SaveProducts function with the updated product
	newProducts, err = db.SaveProducts(updatedProducts)
	if err != nil {
		t.Fatalf("Failed to save updated product: %v", err)
	}

	// Assert that no new products were returned
	if len(newProducts) != 0 {
		t.Errorf("Expected no new products, but got %d", len(newProducts))
	}

	// Retrieve the updated product from the database
	err = db.db.QueryRow("SELECT name, shop, price, link, notified, last_seen FROM "+db.productTableName+" WHERE name = $1", updatedProduct.Name).
		Scan(&retrievedProduct.Name, &retrievedProduct.Shop, &retrievedProduct.Price, &retrievedProduct.Link, &retrievedProduct.Notified, &retrievedProduct.LastSeen)
	if err != nil {
		t.Fatalf("Failed to retrieve updated product: %v", err)
	}

	// Assert the updated product data
	assert.Equal(t, updatedProduct.Name, retrievedProduct.Name, "Name mismatch")
	assert.Equal(t, updatedProduct.Shop, retrievedProduct.Shop, "Shop mismatch")
	assert.Equal(t, updatedProduct.Price, retrievedProduct.Price, "Price mismatch")
	assert.Equal(t, updatedProduct.Link, retrievedProduct.Link, "Link mismatch")
	updatedRounded = updatedProduct.LastSeen.UTC().Round(time.Millisecond)
	retrievedRounded = retrievedProduct.LastSeen.UTC().Round(time.Millisecond)
	assert.Equal(t, updatedRounded.String(), retrievedRounded.String(), "Product last seen timestamps do not match")
	// Notified should stay true when no "real" data changes
	assert.Equal(t, true, retrievedProduct.Notified, "Notified status mismatch")
}
