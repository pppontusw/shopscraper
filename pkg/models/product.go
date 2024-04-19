package models

import (
	"database/sql"
	"time"
)

type Product struct {
	Name          string        `json:"name"`
	Shop          string        `json:"shop"`
	PreviousPrice sql.NullInt64 `json:"previousPrice"`
	Price         int           `json:"price"`
	Link          string        `json:"link"`
	FirstSeen     time.Time     `json:"firstSeen"`
	LastSeen      time.Time     `json:"lastSeen"`
	Notified      bool          `json:"notified"`
}
