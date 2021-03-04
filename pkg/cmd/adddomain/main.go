package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jawr/whois.bi/pkg/internal/cmdutil"
	"github.com/jawr/whois.bi/pkg/internal/domain"
	"github.com/jawr/whois.bi/pkg/internal/user"
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

	if err := cmdutil.LoadDotEnv(); err != nil {
		return errors.WithMessage(err, "LoadDotEnv")
	}

	db, err := cmdutil.SetupDatabase()
	if err != nil {
		return errors.WithMessage(err, "SetupDatabase")
	}
	defer db.Close()

	usr, err := user.GetUser(db, email)
	if err != nil {
		return errors.WithMessagef(err, "GetUser '%s'", email)
	}

	dom := domain.NewDomain(domainName, usr)

	if err := dom.Insert(db); err != nil {
		return errors.WithMessage(err, "Domain.Insert")
	}

	return nil
}
