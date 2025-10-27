package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
	"github.com/Jeomhps/projet-iac-cli/internal/types"
	"github.com/spf13/cobra"
)

var machinesCmd = &cobra.Command{
	Use:   "machines",
	Short: "Manage machines",
}

var machinesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List machines",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		resp, err := cl.Get("/machines", token)
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
	mAddName     string
	mAddHost     string
	mAddPort     int
	mAddUser     string
	mAddPassword string
)

var machinesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a machine (admin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if mAddName == "" || mAddHost == "" || mAddPort <= 0 || mAddUser == "" || mAddPassword == "" {
			return fmt.Errorf("all fields required: --name --host --port --user --password")
		}
		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		m := types.MachineCreate{
			Name:     mAddName,
			Host:     mAddHost,
			Port:     mAddPort,
			User:     mAddUser,
			Password: mAddPassword,
		}
		if cl.ShouldRewrite(m.Host) {
			m.Host = cfg.DockerHostGatewayName
		}
		resp, err := cl.PostJSON("/machines", token, m)
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

var mDelName string

var machinesDelCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a machine (admin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if mDelName == "" {
			return fmt.Errorf("--name is required")
		}
		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		_, err = cl.Delete("/machines/"+mDelName, token)
		if err != nil {
			return err
		}
		fmt.Println("Deleted", mDelName)
		return nil
	},
}

func init() {
	machinesCmd.AddCommand(machinesListCmd)
	machinesCmd.AddCommand(machinesAddCmd)
	machinesCmd.AddCommand(machinesDelCmd)

	machinesAddCmd.Flags().StringVar(&mAddName, "name", "", "Machine name")
	machinesAddCmd.Flags().StringVar(&mAddHost, "host", "", "Machine host (rewritten if localhost/127.0.0.1)")
	machinesAddCmd.Flags().IntVar(&mAddPort, "port", 22, "SSH port")
	machinesAddCmd.Flags().StringVar(&mAddUser, "user", "root", "SSH user")
	machinesAddCmd.Flags().StringVar(&mAddPassword, "password", "", "SSH password")

	machinesDelCmd.Flags().StringVar(&mDelName, "name", "", "Machine name")
}
