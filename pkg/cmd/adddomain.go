package cmd

import (
	"github.com/jawr/whois-bi/pkg/internal/db"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	var email, domainName string

	var adddomainCmd = &cobra.Command{
		Use:   "adddomain",
		Short: "Add a domain to a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			loadDotEnv()

			dbConn, err := db.SetupDatabase()
			if err != nil {
				return errors.WithMessage(err, "SetupDatabase")
			}
			defer dbConn.Close()

			usr, err := user.GetUser(dbConn, email)
			if err != nil {
				return errors.WithMessagef(err, "GetUser '%s'", email)
			}

			dom := domain.NewDomain(domainName, usr)

			if err := dom.Insert(dbConn); err != nil {
				return errors.WithMessage(err, "Domain.Insert")
			}

			return nil
		},
	}

	adddomainCmd.Flags().StringVarP(&email, "email", "u", "", "email of user")
	adddomainCmd.Flags().StringVarP(&domainName, "domain", "p", "", "domain to add to user")

	adddomainCmd.MarkFlagRequired("email")
	adddomainCmd.MarkFlagRequired("domain")

	rootCmd.AddCommand(adddomainCmd)
}
