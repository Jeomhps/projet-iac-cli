package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
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
		var out any
		_ = json.Unmarshal(resp.Body, &out)
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return nil
	},
}
