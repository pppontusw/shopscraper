package database

import (
	"database/sql"
	"log"
	"shopscraper/pkg/models"
	"shopscraper/pkg/utils"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db               *sql.DB
	productTableName string
}

func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

func (p *PostgresDB) Initialize(connStr string, tableName string) error {
	var err error
	p.db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	p.productTableName = tableName
	p.db.SetMaxOpenConns(25)
	p.db.SetMaxIdleConns(10)
	p.db.SetConnMaxLifetime(5 * time.Minute)
	return p.db.Ping()
}

func (p *PostgresDB) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresDB) EnsureProductTableExists() error {
	_, err := p.db.Exec(`
        CREATE TABLE IF NOT EXISTS ` + p.productTableName + ` (
            name TEXT,
            shop TEXT,
			previous_price INT,
            price INT,
            link TEXT,
            first_seen TIMESTAMP,
            last_seen TIMESTAMP,
            notified BOOLEAN,
            UNIQUE (name, shop, link)
        )
    `)
	return err
}

func (p *PostgresDB) GetNonNotifiedProducts() ([]models.Product, error) {
	rows, err := p.db.Query("SELECT name, shop, previous_price, price, link, first_seen, last_seen, notified FROM " + p.productTableName + " WHERE notified = false")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.Name, &product.Shop, &product.PreviousPrice, &product.Price, &product.Link, &product.FirstSeen, &product.LastSeen, &product.Notified)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (p *PostgresDB) GetAllProducts() ([]models.Product, error) {
	rows, err := p.db.Query("SELECT name, shop, previous_price, price, link, first_seen, last_seen, notified FROM " + p.productTableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.Name, &product.Shop, &product.PreviousPrice, &product.Price, &product.Link, &product.FirstSeen, &product.LastSeen, &product.Notified)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (p *PostgresDB) SaveProducts(products []models.Product) ([]models.Product, error) {
	var newProducts []models.Product

	// Prepare the upsert statement outside the loop to avoid re-preparing it for every product
	stmt, err := p.db.Prepare(`INSERT INTO ` + p.productTableName + ` (name, shop, price, link, first_seen, last_seen, notified)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (name, shop, link) DO UPDATE 
        SET price = EXCLUDED.price,
            previous_price = CASE WHEN ` + p.productTableName + `.price != EXCLUDED.price THEN ` + p.productTableName + `.price ELSE ` + p.productTableName + `.previous_price END,
            last_seen = EXCLUDED.last_seen,
            notified = (CASE WHEN ` + p.productTableName + `.price != EXCLUDED.price THEN false ELSE ` + p.productTableName + `.notified END)
        RETURNING (xmax = 0) AS is_inserted;`) // xmax = 0 will return true if it was an insert operation

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, product := range products {
		var isInserted bool
		err := stmt.QueryRow(product.Name, product.Shop, product.Price, product.Link, product.LastSeen, product.LastSeen, product.Notified).Scan(&isInserted)
		if err != nil {
			return nil, err
		}

		// If it was an insert, add to newProducts
		if isInserted {
			newProducts = append(newProducts, product)
		}
	}

	return newProducts, nil
}

func (p *PostgresDB) SetNotifiedProducts(products []models.Product) error {
	// Update notified status for products in the database
	for _, product := range products {
		_, err := p.db.Exec("UPDATE "+p.productTableName+" SET notified = true WHERE name = $1 AND shop = $2 AND link = $3", product.Name, product.Shop, product.Link)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PostgresDB) RemoveOldProducts(timeBack time.Duration) error {
	threshold := utils.GetPastTimeThreshold(timeBack)

	// Delete items from the database where last_seen is older than the threshold
	_, err := p.db.Exec("DELETE FROM "+p.productTableName+" WHERE last_seen < $1", threshold.UTC())
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (p *PostgresDB) DropProductTable() error {
	_, err := p.db.Exec("DROP TABLE IF EXISTS " + p.productTableName)
	return err
}
