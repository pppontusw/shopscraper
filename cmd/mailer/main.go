package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"shopscraper/pkg/config"
	"shopscraper/pkg/database"
	"shopscraper/pkg/mailer"
	"time"

	_ "github.com/lib/pq"
)

var db database.Database

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
	flag.BoolVar(&daemonMode, "daemon", false, "enable daemon mode")
	flag.DurationVar(&interval, "interval", 5*time.Minute, "minimum interval between emails (e.g., 30m, 1h, 2h45m)")
	flag.StringVar(&configPath, "config-path", "./config/config.yaml", "path to configuration yaml file")
	flag.Parse()

	programConfig, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if daemonMode {
		for {
			getAndNotify(*programConfig)
			fmt.Printf("Mailer run finished, waiting %s before next run..\n", interval.String())
			time.Sleep(interval)
		}
	} else {
		getAndNotify(*programConfig)
	}
}

func getAndNotify(programConfig config.ProgramConfig) {
	nonNotifiedProducts, err := db.GetNonNotifiedProducts()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if len(nonNotifiedProducts) > 0 {
		err = mailer.SendEmail(&mailer.RealSmtpSender{}, nonNotifiedProducts, programConfig)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		err = db.SetNotifiedProducts(nonNotifiedProducts)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	} else {
		log.Println("No products found to notify")
	}
}
