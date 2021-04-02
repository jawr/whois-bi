package cmd

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "whoisbi",
	Short: "Whois.bi - Tools",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

// LoadDotEnv reads in .env variables
func loadDotEnv() error {
	return godotenv.Load()
}

// SetupDatabase uses the POSTGRES_URI environment variable
// to create an pg.DB instance
func setupDatabase() (*pg.DB, error) {

	options, err := pg.ParseURL(os.Getenv("POSTGRES_URI"))
	if err != nil {
		return nil, err
	}

	options.ApplicationName = "whois.bi"
	options.TLSConfig = nil

	db := pg.Connect(options)

	return db, nil
}
