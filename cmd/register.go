package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
	"github.com/Jeomhps/projet-iac-cli/internal/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var regFile string

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register machines from a YAML file (admin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if regFile == "" {
			return fmt.Errorf("--file is required")
		}
		data, err := os.ReadFile(regFile)
		if err != nil {
			return err
		}
		var parsed any
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return fmt.Errorf("parse YAML: %w", err)
		}

		var machines []types.MachineCreate
		switch t := parsed.(type) {
		case map[string]any:
			if m, ok := t["machines"]; ok {
				b, _ := json.Marshal(m)
				_ = json.Unmarshal(b, &machines)
			}
		case []any:
			b, _ := json.Marshal(t)
			_ = json.Unmarshal(b, &machines)
		default:
			return fmt.Errorf("YAML must be a list or contain a top-level 'machines' list")
		}

		if len(machines) == 0 {
			fmt.Println("No machines to register.")
			return nil
		}

		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}

		var anyFailed bool
		for _, m := range machines {
			if m.Name == "" || m.Host == "" || m.Port <= 0 || m.User == "" || m.Password == "" {
				fmt.Println("Skipping incomplete entry:", m)
				continue
			}
			if cl.ShouldRewrite(m.Host) {
				m.Host = cfg.DockerHostGatewayName
			}
			resp, err := cl.PostJSON("/machines", token, m)
			if err != nil {
				anyFailed = true
				fmt.Printf("Failed to add %s: %v\n", m.Name, err)
				continue
			}
			if resp.StatusCode >= 300 {
				anyFailed = true
				fmt.Printf("Failed to add %s: %d %s\n", m.Name, resp.StatusCode, string(resp.Body))
				continue
			}
			fmt.Printf("Added %s (%s:%d)\n", m.Name, m.Host, m.Port)
		}

		if anyFailed {
			return fmt.Errorf("one or more machines failed to register")
		}
		return nil
	},
}

func init() {
	registerCmd.Flags().StringVarP(&regFile, "file", "f", "", "Path to machines YAML (e.g., provision/machines.yml)")
}
