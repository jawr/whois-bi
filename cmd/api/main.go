package main

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/jawr/monere/api"
	"github.com/jawr/monere/sender"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := setupDatabase()
	if err != nil {
		return errors.Wrap(err, "setupDatabase")
	}

	emailer := sender.NewSender()

	server := api.NewServer(db, emailer)

	if err := server.Run("localhost:8444"); err != nil {
		return errors.Wrap(err, "Run")
	}

	return nil
}

func setupDatabase() (*pg.DB, error) {
	options, err := pg.ParseURL("postgresql://jawr@/monere")
	if err != nil {
		return nil, errors.Wrap(err, "ParseURL")
	}

	options.Network = "unix"
	options.Addr = "/tmp/.s.PGSQL.5432"
	options.ApplicationName = "monere-schema"
	options.TLSConfig = nil

	db := pg.Connect(options)

	return db, nil
}
