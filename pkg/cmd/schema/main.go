package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/jawr/whois.bi/pkg/internal/cmdutil"
	"github.com/jawr/whois.bi/pkg/internal/domain"
	"github.com/jawr/whois.bi/pkg/internal/job"
	"github.com/jawr/whois.bi/pkg/internal/user"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := cmdutil.LoadDotEnv(); err != nil {
		return errors.WithMessage(err, "LoadDotEnv")
	}

	db, err := cmdutil.SetupDatabase()
	if err != nil {
		return errors.WithMessage(err, "SetupDatabase")
	}
	defer db.Close()

	setupSchema(db)

	return nil
}

func setupSchema(db *pg.DB) {
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
		}
	}
}
