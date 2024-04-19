package scraper

import (
	"fmt"
	"shopscraper/pkg/config"
	"shopscraper/pkg/models"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestParsePrice(t *testing.T) {
	bs := &BaseScraper{
		Config: config.ScraperConfig{
			PriceFormat: "reverse",
		},
	}

	// Test case 1: Price format is "reverse"
	itemPrice := "1.499,00€"
	expectedResult := "1499"
	result := bs.ParsePrice(itemPrice)
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 2: Price format is not set
	bs.Config.PriceFormat = ""
	itemPrice = "1 499.00 EUR"
	expectedResult = "1499"
	result = bs.ParsePrice(itemPrice)
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 3: Price format is unknown
	bs.Config.PriceFormat = "unknown"
	itemPrice = "1499,00€"
	expectedResult = "1499"
	result = bs.ParsePrice(itemPrice)
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

	// Test case 4: Price format is double_eur
	bs.Config.PriceFormat = "double_eur"
	itemPrice = "1 499 EUR 2 055 EUR"
	expectedResult = "1499"
	result = bs.ParsePrice(itemPrice)
	if result != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, result)
	}

}

func TestGetPrice(t *testing.T) {
	bs := &BaseScraper{
		Config: config.ScraperConfig{
			PriceSelector: []string{".price"},
		},
	}

	// Test case 1: Single price element
	doc := `
		<div class="price">1 499,00€</div>
	`
	expectedResult := 1499
	docSelection, err := goquery.NewDocumentFromReader(strings.NewReader(doc))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	result, _ := bs.GetPrice(docSelection.Selection)
	if result != expectedResult {
		t.Errorf("Expected %d, but got %d", expectedResult, result)
	}

	// Test case 2: Multiple price elements, lowest price is selected
	doc = `
		<div class="price">2 499,00€</div>
		<div class="price">1 999,00€</div>
		<div class="price">1 499,00€</div>
	`
	expectedResult = 1499
	docSelection, err = goquery.NewDocumentFromReader(strings.NewReader(doc))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	result, _ = bs.GetPrice(docSelection.Selection)
	if result != expectedResult {
		t.Errorf("Expected %d, but got %d", expectedResult, result)
	}

	// Test case 3: No price element found
	doc = `
		<div class="name">Product Name</div>
	`
	expectedResult = 0
	docSelection, err = goquery.NewDocumentFromReader(strings.NewReader(doc))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	result, err = bs.GetPrice(docSelection.Selection)
	if err != nil {
		if !strings.Contains(err.Error(), fmt.Sprintf("unable to extract price from %s", docSelection.Text())) {
			t.Errorf("Expected error message 'unable to extract price from %s', but got: %s", docSelection.Text(), err)
		}
	} else if result != expectedResult {
		t.Errorf("Expected %d, but got %d", expectedResult, result)
	}

	// Test case 4: Multiple price elements, including broken
	doc = `
		<div class="price">2 499,00€</div>
		<div class="price">1e999,00€</div>
		<div class="price">Fake</div>
		<div class="price">1 499,00€</div>
	`
	expectedResult = 1499
	docSelection, err = goquery.NewDocumentFromReader(strings.NewReader(doc))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	result, _ = bs.GetPrice(docSelection.Selection)
	if result != expectedResult {
		t.Errorf("Expected %d, but got %d", expectedResult, result)
	}
}

func TestParseHTML(t *testing.T) {
	bs := &BaseScraper{
		Config: config.ScraperConfig{
			ItemSelector:     ".item",
			NameSelector:     ".name",
			LinkSelector:     ".link",
			NextPageSelector: ".next",
			PriceSelector:    []string{".price"},
			ShopName:         "Test Shop",
		},
	}

	// Test case 1: Valid HTML content
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
	fetchedUrl := "https://example.com"
	expectedProducts := []models.Product{
		{
			Name:     "Product 1",
			Shop:     "Test Shop",
			Price:    1499,
			Link:     "https://example.com/product1",
			LastSeen: time.Now(),
			Notified: false,
		},
		{
			Name:     "Product 2",
			Shop:     "Test Shop",
			Price:    2999,
			Link:     "https://example.com/product2",
			LastSeen: time.Now(),
			Notified: false,
		},
	}
	expectedNextURL := ""

	products, nextURL, err := bs.ParseHTML(htmlContent, fetchedUrl)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the parsed products match the expected products
	if len(products) != len(expectedProducts) {
		t.Errorf("Expected %d products, but got %d", len(expectedProducts), len(products))
	}

	for i, expected := range expectedProducts {
		assert.Equal(t, expected.Name, products[i].Name, "Product %d Name mismatch", i+1)
		assert.Equal(t, expected.Shop, products[i].Shop, "Product %d Shop mismatch", i+1)
		assert.Equal(t, expected.Price, products[i].Price, "Product %d Price mismatch", i+1)
		assert.Equal(t, expected.Link, products[i].Link, "Product %d Link mismatch", i+1)
		assert.Equal(t, expected.Notified, products[i].Notified, "Product %d Notified status mismatch", i+1)
	}

	// Check if the next URL matches the expected next URL
	if nextURL != expectedNextURL {
		t.Errorf("Expected next URL %s, but got %s", expectedNextURL, nextURL)
	}

	// Test case 2: Valid HTML content with next page
	htmlContent = `
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
		<a class="next" href="/next">Next Page</a>
	`
	expectedNextURL = "https://example.com/next"

	products, nextURL, err = bs.ParseHTML(htmlContent, fetchedUrl)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the parsed products match the expected products
	if len(products) != len(expectedProducts) {
		t.Errorf("Expected %d products, but got %d", len(expectedProducts), len(products))
	}

	for i, expected := range expectedProducts {
		assert.Equal(t, expected.Name, products[i].Name, "Product %d Name mismatch", i+1)
		assert.Equal(t, expected.Shop, products[i].Shop, "Product %d Shop mismatch", i+1)
		assert.Equal(t, expected.Price, products[i].Price, "Product %d Price mismatch", i+1)
		assert.Equal(t, expected.Link, products[i].Link, "Product %d Link mismatch", i+1)
		assert.Equal(t, expected.Notified, products[i].Notified, "Product %d Notified status mismatch", i+1)
	}

	// Check if the next URL matches the expected next URL
	if nextURL != expectedNextURL {
		t.Errorf("Expected next URL %s, but got %s", expectedNextURL, nextURL)
	}

}
