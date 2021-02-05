package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/jawr/whois.bi/internal/cmdutil"
	"github.com/jawr/whois.bi/internal/domain"
	"github.com/jawr/whois.bi/internal/job"
	"github.com/jawr/whois.bi/internal/user"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := cmdutil.SetupDatabase()
	if err != nil {
		return errors.Wrap(err, "SetupDatabase")
	}

	setupSchema(db)

	if len(os.Getenv("ADMIN_EMAIL")) == 0 {
		return nil
	}

	if _, err := user.GetUser(db, os.Getenv("ADMIN_EMAIL")); err == nil {
		return nil
	}

	user, err := user.NewUser(
		os.Getenv("ADMIN_EMAIL"),
		os.Getenv("ADMIN_PASSWORD"),
	)
	if err != nil {
		return err
	}

	user.VerifiedAt = time.Now()

	if err := user.Insert(db); err != nil {
		return err
	}

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
