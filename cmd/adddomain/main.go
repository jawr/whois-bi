package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/jawr/monere/domain"
	"github.com/jawr/monere/user"
	"github.com/pkg/errors"
)

var (
	subdomain = flag.String("subdomain", "", "subdomain to add")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [email] [domain] -subdomain\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {

	args := flag.Args()

	if len(args) != 2 {
		return errors.New("Not enough args")
	}

	email := args[0]
	domainName := args[1]

	db, err := setupDatabase()
	if err != nil {
		return errors.Wrap(err, "setupDatabase")
	}
	defer db.Close()

	usr, err := user.GetUser(db, email)
	if err != nil {
		return errors.Wrap(err, "GetUser")
	}

	dom := domain.NewDomain(domainName, usr)

	if err := dom.Insert(db); err != nil {
		return errors.Wrap(err, "dom.Insert")
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
