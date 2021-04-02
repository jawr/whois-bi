package cmd

import (
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/jawr/whois-bi/pkg/internal/db"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/job"
	"github.com/jawr/whois-bi/pkg/internal/list"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	var schemaCmd = &cobra.Command{
		Use:   "schema",
		Short: "Add a domain to a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			loadDotEnv()

			dbConn, err := db.SetupDatabase()
			if err != nil {
				return errors.WithMessage(err, "SetupDatabase")
			}
			defer dbConn.Close()

			setupSchema(dbConn)

			return nil
		},
	}

	rootCmd.AddCommand(schemaCmd)
}

func setupSchema(db *pg.DB) {
	models := []interface{}{
		(*user.User)(nil),
		(*user.Recover)(nil),
		(*domain.Domain)(nil),
		(*domain.Record)(nil),
		(*domain.Whois)(nil),
		(*job.Job)(nil),
		(*list.List)(nil),
		(*job.Alert)(nil),
		(*job.ExpirationAlert)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:          false,
			FKConstraints: true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
			continue
		}
	}
}
