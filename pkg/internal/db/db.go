package db

import (
	"os"

	"github.com/go-pg/pg/v10"
)

// SetupDatabase uses the POSTGRES_URI environment variable
// to create an pg.DB instance
func SetupDatabase() (*pg.DB, error) {

	options, err := pg.ParseURL(os.Getenv("POSTGRES_URI"))
	if err != nil {
		return nil, err
	}

	options.ApplicationName = "whois.bi"
	options.TLSConfig = nil

	db := pg.Connect(options)

	return db, nil
}
