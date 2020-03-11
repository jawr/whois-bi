package shared

import (
	"os"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
)

func SetupDatabase() (*pg.DB, error) {
	dbURL := os.Getenv("MONERE_DB_URL")
	dbAddr := os.Getenv("MONERE_DB_ADDR")
	dbNetwork := os.Getenv("MONERE_DB_NETWORK")

	if len(dbURL) == 0 {
		dbURL = "postgresql://jawr@/monere"
	}

	if len(dbAddr) == 0 {
		dbAddr = "/tmp/.s.PGSQL.5432"
	}

	options, err := pg.ParseURL(dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "ParseURL")
	}

	if dbNetwork == "unix" {
		options.Network = dbNetwork
		options.Addr = dbAddr
	}
	options.ApplicationName = "monere"
	options.TLSConfig = nil

	db := pg.Connect(options)

	return db, nil
}
