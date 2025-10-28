package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
	"github.com/Jeomhps/projet-iac-cli/internal/configloader"
	"github.com/spf13/cobra"
)

var (
	cfg     client.Config
	version = "dev"
	commit  = "dev"

	// flag vars (separate from cfg so we can control precedence)
	flagConfigPath        string
	flagAPIBase           string
	flagAPIPrefix         string
	flagVerifyTLS         bool
	flagTokenFile         string
	flagRewriteLocalhost  bool
	flagDockerHostGateway string
	flagKeychainMode      string
	flagColorMode         string

	// final output color mode resolved from config/env/flags
	colorMode string
)

var rootCmd = &cobra.Command{
	Use:   "projet-iac-cli",
	Short: "Projet IAC CLI",
	Long:  "CLI for the Projet IAC API (manage users, machines, and reservations).",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Build final cfg from: defaults <- config file <- env <- explicit flags
		if err := buildConfig(cmd); err != nil {
			return err
		}

		// Normalize a couple of fields
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

	// Built-in defaults for the NEW API (no prefix)
	cfg = client.Config{
		APIBase:               "https://localhost/api",
		APIPrefix:             "", // no /api prefix in the new API
		VerifyTLS:             false,
		TokenFile:             defaultToken,
		RewriteLocalhost:      true,
		DockerHostGatewayName: "host.docker.internal",
		KeychainMode:          "auto", // auto|on|off
	}
	colorMode = "auto" // auto|always|never

	// Flags (bind to separate vars so we can decide precedence)
	rootCmd.PersistentFlags().StringVar(&flagConfigPath, "config", getenv("CONFIG_FILE", configloader.DefaultPath()), "Path to config file (YAML)")

	rootCmd.PersistentFlags().StringVar(&flagAPIBase, "api-base", cfg.APIBase, "Base URL (e.g., https://localhost)")
	rootCmd.PersistentFlags().StringVar(&flagAPIPrefix, "api-prefix", cfg.APIPrefix, "API prefix path (leave empty for new API)")
	rootCmd.PersistentFlags().BoolVar(&flagVerifyTLS, "verify-tls", cfg.VerifyTLS, "Verify TLS certificates")
	rootCmd.PersistentFlags().StringVar(&flagTokenFile, "token-file", cfg.TokenFile, "Token cache file (~/.projet-iac/token.json) (used if keychain unavailable/disabled)")
	rootCmd.PersistentFlags().BoolVar(&flagRewriteLocalhost, "rewrite-localhost", cfg.RewriteLocalhost, "Rewrite localhost/127.0.0.1 to host.docker.internal")
	rootCmd.PersistentFlags().StringVar(&flagDockerHostGateway, "docker-host", cfg.DockerHostGatewayName, "Name used when rewriting localhost")
	rootCmd.PersistentFlags().StringVar(&flagKeychainMode, "keychain", cfg.KeychainMode, "Keychain usage: auto|on|off")
	rootCmd.PersistentFlags().StringVar(&flagColorMode, "color", colorMode, "Colorize JSON output: auto|always|never")

	rootCmd.Version = fmt.Sprintf("%s (%s)", version, commit)

	// Commands
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(machinesCmd)
	rootCmd.AddCommand(reservationsCmd)
	rootCmd.AddCommand(reserveCmd)
	rootCmd.AddCommand(registerCmd)
	// Removed: release-all (legacy) and signup (no endpoint in new API)
}

// Build final cfg from defaults <- config file <- env <- flags
func buildConfig(cmd *cobra.Command) error {
	// 1) Config file (if present)
	confPath := flagConfigPath
	if confPath == "" {
		confPath = configloader.DefaultPath()
	}
	if fc, exists, err := configloader.LoadFile(confPath); err != nil {
		return fmt.Errorf("load config file %s: %w", confPath, err)
	} else if exists {
		if fc.APIBase != nil {
			cfg.APIBase = *fc.APIBase
		}
		if fc.APIPrefix != nil {
			cfg.APIPrefix = *fc.APIPrefix
		}
		if fc.VerifyTLS != nil {
			cfg.VerifyTLS = *fc.VerifyTLS
		}
		if fc.TokenFile != nil {
			cfg.TokenFile = *fc.TokenFile
		}
		if fc.RewriteLocalhost != nil {
			cfg.RewriteLocalhost = *fc.RewriteLocalhost
		}
		if fc.DockerHostGatewayName != nil {
			cfg.DockerHostGatewayName = *fc.DockerHostGatewayName
		}
		if fc.KeychainMode != nil && *fc.KeychainMode != "" {
			cfg.KeychainMode = *fc.KeychainMode
		}
		if fc.ColorMode != nil && *fc.ColorMode != "" {
 			colorMode = *fc.ColorMode
		}
	}

	// Helper to check if a flag was explicitly set
	flagChanged := func(name string) bool { return cmd.Flags().Changed(name) }

	// 2) Environment variables (only apply when corresponding flag not explicitly set)
	if !flagChanged("api-base") {
		if v, ok := getenvOpt("API_BASE"); ok {
			cfg.APIBase = v
		}
	}
	if !flagChanged("api-prefix") {
		if v, ok := getenvOpt("API_PREFIX"); ok {
			cfg.APIPrefix = v
		}
	}
	if !flagChanged("verify-tls") {
		if v, ok := envBoolOpt("VERIFY_TLS"); ok {
			cfg.VerifyTLS = v
		}
	}
	if !flagChanged("token-file") {
		if v, ok := getenvOpt("TOKEN_FILE"); ok {
			cfg.TokenFile = v
		}
	}
	if !flagChanged("rewrite-localhost") {
		if v, ok := envBoolOpt("REWRITE_LOCALHOST"); ok {
			cfg.RewriteLocalhost = v
		}
	}
	if !flagChanged("docker-host") {
		if v, ok := getenvOpt("DOCKER_HOST_GATEWAY_NAME"); ok {
			cfg.DockerHostGatewayName = v
		}
	}
	if !flagChanged("keychain") {
		if v, ok := getenvOpt("KEYCHAIN"); ok {
			cfg.KeychainMode = v
		}
	}
	if !flagChanged("color") {
		if v, ok := getenvOpt("COLOR"); ok {
			colorMode = strings.ToLower(strings.TrimSpace(v))
		}
	}

	// 3) Explicit flags override everything
	if flagChanged("api-base") {
		cfg.APIBase = flagAPIBase
	}
	if flagChanged("api-prefix") {
		cfg.APIPrefix = flagAPIPrefix
	}
	if flagChanged("verify-tls") {
		cfg.VerifyTLS = flagVerifyTLS
	}
	if flagChanged("token-file") {
		cfg.TokenFile = flagTokenFile
	}
	if flagChanged("rewrite-localhost") {
		cfg.RewriteLocalhost = flagRewriteLocalhost
	}
	if flagChanged("docker-host") {
		cfg.DockerHostGatewayName = flagDockerHostGateway
	}
	if flagChanged("keychain") {
		cfg.KeychainMode = strings.ToLower(strings.TrimSpace(flagKeychainMode))
	}
	if flagChanged("color") {
		colorMode = strings.ToLower(strings.TrimSpace(flagColorMode))
	}

	return nil
}
