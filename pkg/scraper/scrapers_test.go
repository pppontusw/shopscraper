package scraper

import (
	"shopscraper/pkg/config"
	"testing"
)

func TestCreateScrapers(t *testing.T) {
	// Create a test configuration
	testConfig := &config.ProgramConfig{
		Scrapers: []config.ScraperConfig{
			{
				Type: "WebShopScraper",
			},
			{
				Type: "JavaScriptWebShopScraper",
			},
		},
	}

	// Call the CreateScrapers function with the test configuration
	scrapers, _ := CreateScrapers(testConfig.Scrapers)

	// Assert the expected number of scrapers
	expectedScrapers := 2
	if len(scrapers) != expectedScrapers {
		t.Errorf("Expected %d scrapers, but got %d", expectedScrapers, len(scrapers))
	}

	// Assert the types of the created scrapers
	if _, ok := scrapers[0].(*WebShopScraper); !ok {
		t.Errorf("Expected scraper at index 0 to be of type WebShopScraper")
	}
	if _, ok := scrapers[1].(*JavaScriptWebShopScraper); !ok {
		t.Errorf("Expected scraper at index 1 to be of type JavaScriptWebShopScraper")
	}

	unknownScraper := config.ScraperConfig{
		Type: "UnknownScraper",
	}
	testConfig.Scrapers = append(testConfig.Scrapers, unknownScraper)
	_, err := CreateScrapers(testConfig.Scrapers)

	// Assert the error for unknown scraper type
	if err == nil {
		t.Errorf("Expected an error for unknown scraper type, but got nil")
	}
	expectedErrorMessage := "unknown scraper type 'UnknownScraper'"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMessage, err.Error())
	}
}
