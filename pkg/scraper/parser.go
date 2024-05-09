package scraper

import (
	"fmt"
	"log"
	"shopscraper/pkg/models"
	"shopscraper/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (bs *BaseScraper) ParsePrice(itemPrice string) string {
	switch bs.Config.PriceFormat {
	case "reverse":
		// example: 1.499,00€
		itemPrice = strings.Split(itemPrice, ",")[0]
		itemPrice = strings.ReplaceAll(itemPrice, ".", "")
	case "double_eur":
		// example: 1 499,00EUR 2 500,00EUR
		// this is the case in shops that use one text paragraph
		// with strikethrough text to indicate a sale
		itemPrice = strings.Split(itemPrice, "EUR")[0]
	default:
		itemPrice = strings.Split(itemPrice, ".")[0]
		itemPrice = strings.Split(itemPrice, ",")[0]
	}
	itemPrice = strings.ReplaceAll(itemPrice, " ", "")
	itemPrice = strings.ReplaceAll(itemPrice, "\u00a0", "")
	itemPrice = strings.ReplaceAll(itemPrice, "\t", "")
	itemPrice = strings.ReplaceAll(itemPrice, "\n", "")
	itemPrice = strings.ReplaceAll(itemPrice, "EUR", "")
	itemPrice = strings.ReplaceAll(itemPrice, "€", "")
	return itemPrice
}

func (bs *BaseScraper) GetPrice(s *goquery.Selection) (int, error) {
	var itemPrice string
	var lowestPrice int = 999999

	for _, selector := range bs.Config.PriceSelector {
		findPrice := s.Find(selector)
		if len(findPrice.Nodes) == 1 {
			itemPrice = bs.ParsePrice(findPrice.Text())
			break
		} else {
			findPrice.Each(func(j int, p *goquery.Selection) {
				currentPrice := bs.ParsePrice(p.Text())
				currentPriceInt, err := strconv.Atoi(currentPrice)
				if err != nil {
					log.Println("Error converting price to integer", currentPrice)
					return // continue
				}
				if currentPriceInt < lowestPrice {
					lowestPrice = currentPriceInt
					itemPrice = currentPrice
				}
			})
		}
	}

	if itemPrice == "" {
		return 0, fmt.Errorf("unable to extract price from %s", s.Text())
	}

	itemPriceInt, err := strconv.Atoi(itemPrice)
	if err != nil {
		log.Println("Error converting final price to integer:", itemPrice)
		return 0, err
	}

	return itemPriceInt, nil
}

// ParseHTML parses the HTML content and extracts product information
func (bs *BaseScraper) ParseHTML(htmlContent, fetchedUrl string) ([]models.Product, string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, "", err
	}

	var products []models.Product
	doc.Find(bs.Config.ItemSelector).Each(func(i int, s *goquery.Selection) {
		itemName := s.Find(bs.Config.NameSelector).Text()
		itemName = strings.TrimLeft(itemName, "-. \t\n")
		itemName = strings.TrimRight(itemName, "-. \t\n")

		var itemPrice int
		if len(bs.Config.PriceSelector) != 0 {
			itemPrice, err = bs.GetPrice(s)
			if err != nil {
				log.Printf("Failed to get price %v", err)
			}
		} else {
			itemPrice = 0
		}

		itemLink, _ := s.Find(bs.Config.LinkSelector).Attr("href")
		itemLink, err = utils.EnsureFullUrl(itemLink, fetchedUrl, bs.Config.UniqueParameters)
		if err != nil {
			log.Printf("Failed to get full URL %v", err)
		}

		if itemName != "" && itemLink != "" {
			product := models.Product{
				Name:     itemName,
				Shop:     bs.Config.ShopName,
				Price:    itemPrice,
				Link:     itemLink,
				LastSeen: time.Now().UTC(),
				Notified: false,
			}

			// Check if the product already exists in the products slice
			exists := false
			for _, p := range products {
				if p.Name == product.Name && p.Price == product.Price && p.Link == product.Link {
					exists = true
					break
				}
			}

			// Append the product to the products slice only if it doesn't already exist
			if !exists {
				products = append(products, product)

			}

		}
	})

	// Check for pagination if nextPageSelector is provided
	nextURL := ""
	if bs.Config.NextPageSelector != "" {
		doc.Find(bs.Config.NextPageSelector).Each(func(i int, s *goquery.Selection) {
			if nextURL != "" {
				return
			}
			if href, exists := s.Attr("href"); exists {
				nextURL, err = utils.EnsureFullUrl(href, fetchedUrl, []string{})
				if err != nil {
					log.Printf("Failed to get full URL %v", err)
				}
			}
		})
	}

	return products, nextURL, nil
}
