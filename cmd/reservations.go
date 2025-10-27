package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
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
		var out any
		_ = json.Unmarshal(resp.Body, &out)
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
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
		q := url.Values{}
		q.Set("count", strconv.Itoa(reserveCount))
		q.Set("duration", strconv.Itoa(reserveDuration))
		q.Set("reservation_password", reservePassword)
		if reserveAsUser != "" {
			q.Set("username", reserveAsUser)
		}
		resp, err := cl.Get("/reserve?"+q.Encode(), token)
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

var releaseAllCmd = &cobra.Command{
	Use:   "release-all",
	Short: "Release all reservations (admin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		_, err = cl.Get("/release_all", token)
		if err != nil {
			return err
		}
		fmt.Println("All reservations released.")
		return nil
	},
}

func init() {
	reserveCmd.Flags().IntVar(&reserveCount, "count", 1, "Number of machines")
	reserveCmd.Flags().IntVar(&reserveDuration, "duration", 60, "Duration in minutes")
	reserveCmd.Flags().StringVar(&reservePassword, "password", "", "Reservation password to set on machines")
	reserveCmd.Flags().StringVar(&reserveAsUser, "as-user", "", "Logical username to reserve for (defaults to API user)")
}
