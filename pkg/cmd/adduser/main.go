package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jawr/whois.bi/pkg/internal/cmdutil"
	"github.com/jawr/whois.bi/pkg/internal/user"
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
	if err := cmdutil.LoadDotEnv(); err != nil {
		return errors.WithMessage(err, "LoadDotEnv")
	}

	db, err := cmdutil.SetupDatabase()
	if err != nil {
		return errors.WithMessage(err, "SetupDatabase")
	}
	defer db.Close()

	u, err := user.NewUser(*email, *password)
	if err != nil {
		return errors.WithMessage(err, "NewUser")
	}

	if err := u.Insert(db); err != nil {
		return errors.WithMessage(err, "User.Insert")
	}

	return nil
}
