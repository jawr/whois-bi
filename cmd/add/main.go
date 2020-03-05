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

const (
	DefaultEmail    string = "jess@lawrence.pm"
	DefaultPassword string = "jess@lawrence.pm"
)

var (
	subdomain = flag.String("subdomain", "", "subdomain to add")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [domain] -subdomain\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var domainName string

	for _, arg := range flag.Args() {
		domainName = arg
		break
	}

	if len(domainName) == 0 {
		return errors.New("no domain passed")
	}

	db, err := setupDatabase()
	if err != nil {
		return errors.Wrap(err, "setupDatabase")
	}
	defer db.Close()

	usr, err := ensureUser(db)
	if err != nil {
		return errors.Wrap(err, "ensureUser")
	}

	dom := domain.NewDomain(domainName, usr)

	if err := dom.Insert(db); err != nil {
		return errors.Wrap(err, "dom.Insert")
	}

	return nil
}

func ensureUser(db *pg.DB) (user.User, error) {
	u, err := user.GetUser(db, DefaultEmail)
	if err != nil {
		u, err = user.NewUser(DefaultEmail, DefaultPassword)
		if err != nil {
			return user.User{}, errors.Wrap(err, "NewUser")
		}

		if err := u.Insert(db); err != nil {
			return user.User{}, errors.Wrap(err, "Insert")
		}
	}

	return u, nil
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
