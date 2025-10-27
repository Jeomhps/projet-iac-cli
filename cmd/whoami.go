package cmd

import (
	"fmt"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
	"github.com/Jeomhps/projet-iac-cli/internal/output"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user info (/whoami)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		resp, err := cl.Get("/whoami", token)
		if err != nil {
			return err
		}
		fmt.Println(output.FormatJSON(resp.Body, colorMode))
		return nil
	},
}
