package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	cfg     client.Config
	version = "0.1.0"
	commit  = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "projet-iac-cli",
	Short: "Projet IAC CLI",
	Long:  "CLI for the Projet IAC API (manage users, machines, and reservations).",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg.APIBase = strings.TrimRight(cfg.APIBase, "/")
		if cfg.APIPrefix != "" && !strings.HasPrefix(cfg.APIPrefix, "/") {
			cfg.APIPrefix = "/" + cfg.APIPrefix
		}
		return nil
	},
}

func init() {
	home, _ := os.UserHomeDir()
	defaultToken := filepath.Join(home, ".projet-iac", "token.json")

	cfg = client.Config{
		APIBase:               getenv("API_BASE", "https://localhost"),
		APIPrefix:             getenv("API_PREFIX", "/api"),
		VerifyTLS:             envBool("VERIFY_TLS", false),
		TokenFile:             getenv("TOKEN_FILE", defaultToken),
		RewriteLocalhost:      envBool("REWRITE_LOCALHOST", true),
		DockerHostGatewayName: getenv("DOCKER_HOST_GATEWAY_NAME", "host.docker.internal"),
		KeychainMode:          getenv("KEYCHAIN", "auto"), // auto|on|off
	}

	rootCmd.PersistentFlags().StringVar(&cfg.APIBase, "api-base", cfg.APIBase, "Base URL (e.g., https://localhost)")
	rootCmd.PersistentFlags().StringVar(&cfg.APIPrefix, "api-prefix", cfg.APIPrefix, "API prefix path (e.g., /api)")
	rootCmd.PersistentFlags().BoolVar(&cfg.VerifyTLS, "verify-tls", cfg.VerifyTLS, "Verify TLS certificates")
	rootCmd.PersistentFlags().StringVar(&cfg.TokenFile, "token-file", cfg.TokenFile, "Token cache file (~/.projet-iac/token.json) (used if keychain unavailable or disabled)")
	rootCmd.PersistentFlags().BoolVar(&cfg.RewriteLocalhost, "rewrite-localhost", cfg.RewriteLocalhost, "Rewrite localhost/127.0.0.1 to host.docker.internal")
	rootCmd.PersistentFlags().StringVar(&cfg.DockerHostGatewayName, "docker-host", cfg.DockerHostGatewayName, "Name used when rewriting localhost")
	rootCmd.PersistentFlags().StringVar(&cfg.KeychainMode, "keychain", cfg.KeychainMode, "Keychain usage: auto|on|off")

	rootCmd.Version = fmt.Sprintf("%s (%s)", version, commit)

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(machinesCmd)
	rootCmd.AddCommand(reservationsCmd)
	rootCmd.AddCommand(reserveCmd)
	rootCmd.AddCommand(releaseAllCmd)
	rootCmd.AddCommand(registerCmd)
}
