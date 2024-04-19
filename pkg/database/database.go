package database

import (
	"shopscraper/pkg/models"
)

type Database interface {
	Initialize(connStr string, tableName string) error
	Close() error
	EnsureProductTableExists() error
	GetNonNotifiedProducts() ([]models.Product, error)
	GetAllProducts() ([]models.Product, error)
	SaveProducts(products []models.Product) ([]models.Product, error)
	SetNotifiedProducts(products []models.Product) error
	RemoveOldProducts() error
	DropProductTable() error
}
