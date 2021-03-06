package cmd

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/jawr/whois-bi/pkg/internal/cmdutil"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/job"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	var schemaCmd = &cobra.Command{
		Use:   "schema",
		Short: "Add a domain to a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdutil.LoadDotEnv()

			db, err := cmdutil.SetupDatabase()
			if err != nil {
				return errors.WithMessage(err, "SetupDatabase")
			}
			defer db.Close()

			setupSchema(db)

			return nil
		},
	}

	rootCmd.AddCommand(schemaCmd)
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
			fmt.Fprintf(os.Stderr, "Error: %s", err)
			continue
		}
	}
}
