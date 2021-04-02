package cmd

import (
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	var email, password string

	var adduserCmd = &cobra.Command{
		Use:   "adduser",
		Short: "Add a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			loadDotEnv()

			db, err := setupDatabase()
			if err != nil {
				return errors.WithMessage(err, "SetupDatabase")
			}
			defer db.Close()

			u, err := user.NewUser(email, password)
			if err != nil {
				return errors.WithMessage(err, "NewUser")
			}

			if err := u.Insert(db); err != nil {
				return errors.WithMessage(err, "User.Insert")
			}

			return nil
		},
	}

	adduserCmd.Flags().StringVarP(&email, "email", "u", "", "email to add")
	adduserCmd.Flags().StringVarP(&password, "password", "p", "", "password for new user")

	adduserCmd.MarkFlagRequired("email")
	adduserCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(adduserCmd)
}
