package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ScraperConfig struct {
	Type             string   `yaml:"type"`
	ShopName         string   `yaml:"shopName"`
	URLs             []string `yaml:"urls"`
	ItemSelector     string   `yaml:"itemSelector"`
	NameSelector     string   `yaml:"nameSelector"`
	PriceSelector    []string `yaml:"priceSelector"`
	LinkSelector     string   `yaml:"linkSelector"`
	NextPageSelector string   `yaml:"nextPageSelector"`
	PriceFormat      string   `yaml:"priceFormat"`
	RetryString      string   `yaml:"retryString"`
}

type EmailConfig struct {
	Sender    string `yaml:"sender"`
	Recipient string `yaml:"recipient"`
	Subject   string `yaml:"subject"`
	Server    string `yaml:"server"`
	Port      string `yaml:"port"`
}

type ProgramConfig struct {
	Scrapers []ScraperConfig `yaml:"scrapers"`
	Email    EmailConfig     `yaml:"email"`
}

// readConfig reads the YAML configuration file and returns the ScraperConfig struct
func ReadConfig(filepath string) (*ProgramConfig, error) {
	configData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var config ProgramConfig
	err = yaml.Unmarshal(configData, &config)

	if err != nil {
		return nil, err
	}
	return &config, nil
}
