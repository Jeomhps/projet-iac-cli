package cmd

import (
	"fmt"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
	"github.com/Jeomhps/projet-iac-cli/internal/output"
	"github.com/spf13/cobra"
)

var reservationsCmd = &cobra.Command{
	Use:   "reservations",
	Short: "List active reservations",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		resp, err := cl.Get("/reservations", token)
		if err != nil {
			return err
		}
		fmt.Println(output.FormatJSON(resp.Body, colorMode))
		return nil
	},
}

var (
	reserveCount    int
	reserveDuration int
	reservePassword string
	reserveAsUser   string
)

var reserveCmd = &cobra.Command{
	Use:   "reserve",
	Short: "Reserve N machines",
	RunE: func(cmd *cobra.Command, args []string) error {
		if reserveCount <= 0 || reserveDuration <= 0 || reservePassword == "" {
			return fmt.Errorf("--count, --duration and --password are required and must be > 0")
		}
		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		payload := map[string]any{
			"count":                reserveCount,
			"duration_minutes":     reserveDuration,
			"reservation_password": reservePassword,
		}
		if reserveAsUser != "" {
			payload["username"] = reserveAsUser
		}
		resp, err := cl.PostJSON("/reservations", token, payload)
		if err != nil {
			return err
		}
		fmt.Println(output.FormatJSON(resp.Body, colorMode))
		return nil
	},
}

func init() {
	reserveCmd.Flags().IntVar(&reserveCount, "count", 1, "Number of machines")
	reserveCmd.Flags().IntVar(&reserveDuration, "duration", 60, "Duration in minutes")
	reserveCmd.Flags().StringVar(&reservePassword, "password", "", "Reservation password to set on machines")
	reserveCmd.Flags().StringVar(&reserveAsUser, "as-user", "", "Logical username to reserve for (defaults to API user)")
}
