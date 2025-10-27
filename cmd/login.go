package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	loginUsername string
	loginPassword string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login and cache token (OS keychain when available)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := client.New(cfg)

		u := strings.TrimSpace(loginUsername)
		p := loginPassword

		reader := bufio.NewReader(os.Stdin)
		if u == "" {
			fmt.Print("Username: ")
			uu, _ := reader.ReadString('\n')
			u = strings.TrimSpace(uu)
		}
		if p == "" {
			fmt.Print("Password: ")
			b, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				return fmt.Errorf("reading password: %w", err)
			}
			p = string(b)
		}
		if u == "" || p == "" {
			return fmt.Errorf("username and password are required")
		}

		token, exp, err := cl.Login(u, p)
		if err != nil {
			return err
		}
		if err := cl.SaveToken(token, exp); err != nil {
			return err
		}

		if cl.UsingKeychain() {
			fmt.Println("Logged in. Token stored in OS keychain.")
		} else {
			fmt.Println("Logged in. Token cached at:", cfg.TokenFile)
		}
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Username")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password (omit to prompt securely)")
}
