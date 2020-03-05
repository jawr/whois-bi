package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/jawr/monere/domain"
	"github.com/jawr/monere/job"
	"github.com/jawr/monere/user"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := setupDatabase()
	if err != nil {
		return errors.Wrap(err, "setupDatabase")
	}

	if err := setupSchema(db); err != nil {
		return errors.Wrap(err, "setupSchema")
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

func setupSchema(db *pg.DB) error {
	models := []interface{}{
		(*user.User)(nil),
		(*domain.Domain)(nil),
		(*domain.Record)(nil),
		(*domain.Whois)(nil),
		(*job.Job)(nil),
	}

	for _, model := range models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:          false,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error: %s", err)
			continue
			return err
		}
	}
	return nil
}
