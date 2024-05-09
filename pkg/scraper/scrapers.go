package scraper

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"shopscraper/pkg/config"
	"shopscraper/pkg/models"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const MaxRetryAttempts = 5

type HTMLGetter interface {
	GetHTML(currentURL string, attempts ...int) (string, error)
}

type Scraper interface {
	Scrape(int) ([]models.Product, error)
	ParseHTML(htmlContent, fetchedUrl string) ([]models.Product, string, error)
	GetPrice(s *goquery.Selection) (int, error)
	ParsePrice(itemPrice string) string
}

type BaseScraper struct {
	HTMLGetter
	Config config.ScraperConfig
}

func (bs *BaseScraper) Scrape(maxWorkers int) ([]models.Product, error) {
	var products []models.Product
	log.Println("Starting scraping of", bs.Config.ShopName)

	// Create a channel to receive scraped products
	productChan := make(chan []models.Product)

	// Create a channel to limit the number of concurrent workers
	semaphore := make(chan struct{}, maxWorkers)

	var wg sync.WaitGroup

	for _, url := range bs.Config.URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Acquire a worker from the semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			currentURL := url

			for {
				log.Println("Scraping", currentURL)

				htmlContent, err := bs.GetHTML(currentURL)
				if err != nil {
					log.Println("Error scraping", currentURL, ":", err)
					return
				}

				p, nextURL, err := bs.ParseHTML(htmlContent, url)
				if err != nil {
					log.Println("Error parsing HTML from", currentURL, ":", err)
					return
				}

				productChan <- p

				if nextURL == "" {
					break // Exit loop if there's no next URL
				}

				if nextURL == currentURL {
					break // Exit loop if the next URL is the same as the current URL
				}

				currentURL = nextURL
			}
		}(url)
	}

	// Close the product channel when all goroutines are done
	go func() {
		wg.Wait()
		close(productChan)
	}()

	// Collect scraped products from the channel
	for p := range productChan {
		products = append(products, p...)
	}

	return products, nil
}

type JavaScriptWebShopScraper struct {
	BaseScraper
}

func NewJavaScriptWebShopScraper(config config.ScraperConfig) *JavaScriptWebShopScraper {
	js := &JavaScriptWebShopScraper{
		BaseScraper: BaseScraper{
			Config: config,
		},
	}
	js.HTMLGetter = js
	return js
}

func (js *JavaScriptWebShopScraper) GetHTML(currentURL string, attempts ...int) (string, error) {
	attempt := 1
	if len(attempts) > 0 {
		attempt = attempts[0]
	}
	// create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new browser window
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var htmlContent string

	err := chromedp.Run(ctx,
		chromedp.Navigate(currentURL),
		chromedp.Sleep(5*time.Second), // Wait for JavaScript to render
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			if attempt <= MaxRetryAttempts {
				log.Println("Timed out, retrying", currentURL)
				return js.GetHTML(currentURL, attempt+1)
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	}

	// retry string being present means loading failed
	if js.Config.RetryString != "" {
		if strings.Contains(htmlContent, js.Config.RetryString) {
			if attempt <= MaxRetryAttempts {
				log.Println("Data did not load correctly, retrying", currentURL)
				return js.GetHTML(currentURL, attempt+1)
			} else {
				return "", fmt.Errorf("data did not load correctly on %s", currentURL)
			}
		}
	}

	return htmlContent, nil
}

type WebShopScraper struct {
	BaseScraper
}

// NewWebShopScraper creates a new instance of WebShopScraper
func NewWebShopScraper(config config.ScraperConfig) *WebShopScraper {
	ws := &WebShopScraper{
		BaseScraper: BaseScraper{
			Config: config,
		},
	}
	ws.HTMLGetter = ws
	return ws
}

func (ws *WebShopScraper) GetHTML(currentURL string, attempts ...int) (string, error) {
	resp, err := http.Get(currentURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL %s: status code %d", currentURL, resp.StatusCode)
	}

	htmlContent, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return "", err
	}

	return string(htmlContent), nil
}

func CreateScrapers(config []config.ScraperConfig) ([]Scraper, error) {
	var scrapers []Scraper
	for _, scraperConfig := range config {
		switch scraperConfig.Type {
		case "WebShopScraper":
			scraper := NewWebShopScraper(scraperConfig)
			scrapers = append(scrapers, scraper)
		case "JavaScriptWebShopScraper":
			scraper := NewJavaScriptWebShopScraper(scraperConfig)
			scrapers = append(scrapers, scraper)
		default:
			return nil, fmt.Errorf("unknown scraper type '%s'", scraperConfig.Type)
		}
	}
	return scrapers, nil

}
