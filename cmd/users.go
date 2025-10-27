package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
	"github.com/Jeomhps/projet-iac-cli/internal/output"
	"github.com/Jeomhps/projet-iac-cli/internal/types"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users (admin)",
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users (admin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		resp, err := cl.Get("/users", token)
		if err != nil {
			return err
		}
		fmt.Println(output.FormatJSON(resp.Body, colorMode))
		return nil
	},
}

var (
	uCreateUsername string
	uCreateIsAdmin  bool
	uPasswordStdin  bool
	uNoConfirm      bool
)

var usersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a user (admin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(uCreateUsername) == "" {
			return fmt.Errorf("--username is required")
		}

		// Get password securely
		var password string
		if uPasswordStdin {
			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("reading password from stdin: %w", err)
			}
			password = strings.TrimRight(string(b), "\r\n")
			if password == "" {
				return fmt.Errorf("password from stdin is empty")
			}
		} else {
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				return fmt.Errorf("stdin is not a TTY. Use --password-stdin to provide the password via stdin")
			}
			fmt.Print("Password: ")
			p1, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				return fmt.Errorf("reading password: %w", err)
			}
			if len(p1) == 0 {
				return fmt.Errorf("password cannot be empty")
			}
			if !uNoConfirm {
				fmt.Print("Confirm password: ")
				p2, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()
				if err != nil {
					return fmt.Errorf("reading confirmation: %w", err)
				}
				if string(p1) != string(p2) {
					return fmt.Errorf("passwords do not match")
				}
			}
			password = string(p1)
		}

		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		payload := types.UserCreate{
			Username: uCreateUsername,
			Password: password,
			IsAdmin:  uCreateIsAdmin,
		}
		resp, err := cl.PostJSON("/users", token, payload)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 300 {
			return fmt.Errorf("create failed: %d %s", resp.StatusCode, string(resp.Body))
		}
		fmt.Println(output.FormatJSON(resp.Body, colorMode))
		return nil
	},
}

var (
	uDeleteUsername string
)

var usersDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user (admin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(uDeleteUsername) == "" {
			return fmt.Errorf("--username is required")
		}
		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		resp, err := cl.Delete("/users/"+uDeleteUsername, token)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 300 {
			return fmt.Errorf("delete failed: %d %s", resp.StatusCode, string(resp.Body))
		}
		fmt.Println("Deleted user", uDeleteUsername)
		return nil
	},
}

// Self-service signup (no token required)
var (
	signupUsername  string
	signupPassword  string
	signupAutoLogin bool
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Self-register a new user (POST /register)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(signupUsername) == "" {
			return fmt.Errorf("--username is required")
		}
		pass := strings.TrimSpace(signupPassword)
		if pass == "" {
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				return fmt.Errorf("stdin is not a TTY. Provide --password or add a --password-stdin flow if needed")
			}
			fmt.Print("Password: ")
			p1, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				return fmt.Errorf("reading password: %w", err)
			}
			if len(p1) == 0 {
				return fmt.Errorf("password cannot be empty")
			}
			fmt.Print("Confirm password: ")
			p2, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				return fmt.Errorf("reading confirmation: %w", err)
			}
			if string(p1) != string(p2) {
				return fmt.Errorf("passwords do not match")
			}
			pass = string(p1)
		}

		cl := client.New(cfg)
		payload := types.UserSignup{
			Username: signupUsername,
			Password: pass,
		}
		resp, err := cl.PostJSON("/register", "", payload)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 300 {
			return fmt.Errorf("signup failed: %d %s", resp.StatusCode, string(resp.Body))
		}
		fmt.Println("Signup successful.")
		if signupAutoLogin {
			token, exp, err := cl.Login(signupUsername, pass)
			if err != nil {
				return fmt.Errorf("auto-login failed: %w", err)
			}
			if err := cl.SaveToken(token, exp); err != nil {
				return fmt.Errorf("saving token failed: %w", err)
			}
			fmt.Println("Logged in and token cached.")
		}
		return nil
	},
}

func init() {
	// Wire commands into root
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(signupCmd)

	// users subcommands
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersCreateCmd)
	usersCmd.AddCommand(usersDeleteCmd)

	// Flags
	usersCreateCmd.Flags().StringVar(&uCreateUsername, "username", "", "Username")
	usersCreateCmd.Flags().BoolVar(&uCreateIsAdmin, "admin", false, "Set user as admin")
	usersCreateCmd.Flags().BoolVar(&uPasswordStdin, "password-stdin", false, "Read password from STDIN (for automation)")
	usersCreateCmd.Flags().BoolVar(&uNoConfirm, "no-confirm", false, "Skip password confirmation (useful with --password-stdin)")

	usersDeleteCmd.Flags().StringVar(&uDeleteUsername, "username", "", "Username to delete")

	signupCmd.Flags().StringVar(&signupUsername, "username", "", "Username")
	signupCmd.Flags().StringVar(&signupPassword, "password", "", "Password (omit to be securely prompted)")
	signupCmd.Flags().BoolVar(&signupAutoLogin, "login", false, "Automatically log in and cache token after signup")
}
