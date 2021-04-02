package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/jawr/whois-bi/pkg/internal/db"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/job"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	dbConn, err := db.SetupDatabase()
	if err != nil {
		return errors.Wrap(err, "SetupDatabase")
	}
	defer dbConn.Close()

	setupSchema(dbConn)

	if len(os.Getenv("ADMIN_EMAIL")) == 0 {
		return nil
	}

	if _, err := user.GetUser(dbConn, os.Getenv("ADMIN_EMAIL")); err == nil {
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

	if err := user.Insert(dbConn); err != nil {
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
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:          false,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error: %s", err)
			continue
		}
	}
}
