package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"shopscraper/pkg/config"
	"shopscraper/pkg/database"
	"shopscraper/pkg/models"
	"shopscraper/pkg/scraper"

	_ "github.com/lib/pq"
)

var debugMode bool
var db database.Database

type Scraper = scraper.Scraper

func main() {
	// Initialize the database connection pool
	connectionString := os.Getenv("SHOPSCRAPER_DB_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatalf("SHOPSCRAPER_DB_CONNECTION_STRING not provided")
	}
	db = database.NewPostgresDB()
	db.Initialize(connectionString, "products")
	defer db.Close()

	var configPath string
	var daemonMode bool
	var interval time.Duration
	var maxWorkers int
	flag.BoolVar(&daemonMode, "daemon", false, "enable daemon mode")
	flag.DurationVar(&interval, "interval", 1*time.Hour, "interval between scrapes (e.g., 30m, 1h, 2h45m)")
	flag.IntVar(&maxWorkers, "max-workers", 3, "maximum number of workers per scraper")
	flag.BoolVar(&debugMode, "debug", false, "enable debug mode")
	flag.StringVar(&configPath, "config-path", "./config/config.yaml", "path to configuration yaml file")
	flag.Parse()

	err := db.EnsureProductTableExists()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	programConfig, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Create a list of scrapers
	scrapers, err := scraper.CreateScrapers(programConfig.Scrapers)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if daemonMode {
		for {
			runScrapers(scrapers, maxWorkers)
			fmt.Printf("\nRun finished, waiting %s before next run..\n", interval.String())
			time.Sleep(interval)
		}
	} else {
		runScrapers(scrapers, maxWorkers)
	}
}

// runScrapers runs the scrapers and processes the results
func runScrapers(scrapers []scraper.Scraper, maxWorkers int) {
	// Create channels to receive scraped products from each scraper
	productChans := make([]chan []models.Product, len(scrapers))

	// Run each scraper concurrently
	var wg sync.WaitGroup

	for i, scraper := range scrapers {
		productChans[i] = make(chan []models.Product)
		wg.Add(1)
		go func(i int, scraper Scraper) {
			defer wg.Done()
			products, err := scraper.Scrape(maxWorkers)
			if err != nil {
				log.Printf("Failed to scrape using scraper %d: %v", i, err)
			}
			productChans[i] <- products

		}(i, scraper)
	}

	// Close the channels after all scrapers have finished
	go func() {
		wg.Wait()
		for _, ch := range productChans {
			close(ch)
		}
	}()

	// Aggregate the output from each scraper into a single list
	var allProducts []models.Product
	for _, ch := range productChans {
		products := <-ch
		allProducts = append(allProducts, products...)
	}

	if debugMode {
		log.Println("All products:")
		for _, p := range allProducts {
			fmt.Printf("%s - %s, Price: %d, Link: %s\n", p.Shop, p.Name, p.Price, p.Link)
		}
	}

	newProducts, err := db.SaveProducts(allProducts)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Remove items from the database that haven't been seen in 3 days or more
	db.RemoveOldProducts()

	if debugMode {
		log.Println("New products:")
		// Print all new products
		for _, p := range newProducts {
			fmt.Printf("%s - %s, Price: %d, Link: %s\n", p.Shop, p.Name, p.Price, p.Link)
		}
	}

}
