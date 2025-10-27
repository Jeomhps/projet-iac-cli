package cmd

import (
	"fmt"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Delete cached token",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := client.New(cfg)
		if err := cl.DeleteToken(); err != nil {
			return err
		}
		fmt.Println("Logged out; cached token removed.")
		return nil
	},
}
