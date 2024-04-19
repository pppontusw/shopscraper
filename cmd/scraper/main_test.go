package main

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"shopscraper/pkg/config"
	"shopscraper/pkg/database"
	"shopscraper/pkg/models"
	"shopscraper/pkg/scraper"

	"github.com/PuerkitoBio/goquery"
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
	runID := generateRandomString(16)
	productTableName := "test_products_" + runID

	db = database.NewPostgresDB()
	db.Initialize("postgresql://test:Test1234@localhost:5432/test?sslmode=disable", productTableName)
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

func TestRunScrapersWithMockScraper(t *testing.T) {
	setup(t)
	defer teardown(t)

	var currentTime = time.Now().UTC()

	var product1 = models.Product{Shop: "Shop1", Name: "Product1", Price: 10, Link: "https://example.com/product1", LastSeen: currentTime}
	var product2 = models.Product{Shop: "Shop1", Name: "Product2", Price: 20, Link: "https://example.com/product2", LastSeen: currentTime}
	var product3 = models.Product{Shop: "Shop2", Name: "Product3", Price: 30, Link: "https://example.com/product3", LastSeen: currentTime}
	var product4 = models.Product{Shop: "Shop2", Name: "Product4", Price: 40, Link: "https://example.com/product4", LastSeen: currentTime}

	var expectedProducts = []models.Product{product1, product2, product3, product4}

	mockScrapers := []scraper.Scraper{
		&mockScraper{products: []models.Product{product1, product2}},
		&mockScraper{products: []models.Product{product3, product4}},
	}

	maxWorkers := 2
	runScrapers(mockScrapers, maxWorkers)

	retrievedProducts, err := db.GetAllProducts()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	assert.Equal(t, len(expectedProducts), len(retrievedProducts), "The number of products does not match")

	for i, expected := range expectedProducts {
		assert.Equal(t, expected.Name, retrievedProducts[i].Name, "Product names do not match")
		assert.Equal(t, expected.Shop, retrievedProducts[i].Shop, "Product shops do not match")
		assert.Equal(t, expected.Price, retrievedProducts[i].Price, "Product prices do not match")
		assert.Equal(t, expected.Link, retrievedProducts[i].Link, "Product links do not match")
		expectedRounded := expected.LastSeen.UTC().Round(time.Millisecond)
		retrievedRounded := retrievedProducts[i].LastSeen.UTC().Round(time.Millisecond)
		assert.Equal(t, expectedRounded.String(), retrievedRounded.String(), "Product last seen timestamps do not match")
	}
}
func TestRunScrapersWithMockGetHTML(t *testing.T) {
	setup(t)
	defer teardown(t)

	htmlContent := `
		<div class="item">
			<div class="name">Product 1</div>
			<div class="price">1 499,00€</div>
			<a class="link" href="/product1">Product 1 Link</a>
		</div>
		<div class="item">
			<div class="name">Product 2</div>
			<div class="price">2 999,00€</div>
			<a class="link" href="/product2">Product 2 Link</a>
		</div>
	`

	config := config.ScraperConfig{
		URLs:             []string{"http://example.com"},
		ItemSelector:     ".item",
		NameSelector:     ".name",
		LinkSelector:     ".link",
		NextPageSelector: ".next",
		PriceSelector:    []string{".price"},
		ShopName:         "Test Shop",
	}

	scrapers := []scraper.Scraper{
		&scraper.BaseScraper{
			Config: config,
			HTMLGetter: &mockHTMLGetter{
				HTMLContent: htmlContent,
			},
		},
	}

	maxWorkers := 1
	runScrapers(scrapers, maxWorkers)

	products, err := db.GetAllProducts()
	if err != nil {
		log.Fatalf("error %v", err)
	}

	expectedProducts := []models.Product{
		{
			Name:     "Product 1",
			Shop:     "Test Shop",
			Price:    1499,
			Link:     "http://example.com/product1",
			LastSeen: time.Now(),
			Notified: false,
		},
		{
			Name:     "Product 2",
			Shop:     "Test Shop",
			Price:    2999,
			Link:     "http://example.com/product2",
			LastSeen: time.Now(),
			Notified: false,
		},
	}
	for i, expected := range expectedProducts {
		assert.Equal(t, expected.Name, products[i].Name, "Product %d Name mismatch", i+1)
		assert.Equal(t, expected.Shop, products[i].Shop, "Product %d Shop mismatch", i+1)
		assert.Equal(t, expected.Price, products[i].Price, "Product %d Price mismatch", i+1)
		assert.Equal(t, expected.Link, products[i].Link, "Product %d Link mismatch", i+1)
		assert.Equal(t, expected.Notified, products[i].Notified, "Product %d Notified status mismatch", i+1)
	}

}

// Mock scraper implementation returns specific products it's set up with
type mockScraper struct {
	products []models.Product
}

func (s *mockScraper) Scrape(maxWorkers int) ([]models.Product, error) {
	return s.products, nil
}

func (s *mockScraper) ParseHTML(htmlContent, fetchedUrl string) ([]models.Product, string, error) {
	return s.products, "", nil
}

func (s *mockScraper) GetPrice(q *goquery.Selection) (int, error) {
	return 0, nil
}

func (s *mockScraper) ParsePrice(itemPrice string) string {
	return ""
}

// Mock HTML getter returns specific HTMLContent it's set up with
type mockHTMLGetter struct {
	HTMLContent string
}

func (c *mockHTMLGetter) GetHTML(currentURL string, attempts ...int) (string, error) {
	return c.HTMLContent, nil
}
