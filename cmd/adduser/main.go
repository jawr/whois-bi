package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jawr/monere/shared"
	"github.com/jawr/monere/user"
	"github.com/pkg/errors"
)

var (
	email    = flag.String("email", "", "email to add")
	password = flag.String("password", "", "password to add")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -email -password\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := shared.SetupDatabase()
	if err != nil {
		return errors.Wrap(err, "setupDatabase")
	}
	defer db.Close()

	u, err := user.NewUser(*email, *password)
	if err != nil {
		return errors.Wrap(err, "NewUser")
	}

	if err := u.Insert(db); err != nil {
		return errors.Wrap(err, "Insert")
	}

	return nil
}
