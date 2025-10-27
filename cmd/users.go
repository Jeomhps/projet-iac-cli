package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Jeomhps/projet-iac-cli/internal/client"
	"github.com/Jeomhps/projet-iac-cli/internal/types"
	"github.com/spf13/cobra"
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
		var out any
		_ = json.Unmarshal(resp.Body, &out)
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return nil
	},
}

var (
	uCreateUsername string
	uCreatePassword string
	uCreateIsAdmin  bool
)

var usersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a user (admin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if uCreateUsername == "" || uCreatePassword == "" {
			return fmt.Errorf("--username and --password are required")
		}
		cl := client.New(cfg)
		token, err := cl.GetToken()
		if err != nil {
			return err
		}
		payload := types.UserCreate{
			Username: uCreateUsername,
			Password: uCreatePassword,
			IsAdmin:  uCreateIsAdmin,
		}
		resp, err := cl.PostJSON("/users", token, payload)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 300 {
			return fmt.Errorf("create failed: %d %s", resp.StatusCode, string(resp.Body))
		}
		var out any
		_ = json.Unmarshal(resp.Body, &out)
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return nil
	},
}

var uDeleteUsername string

var usersDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user (admin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if uDeleteUsername == "" {
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

// Self-service signup (no token)
var (
	signupUsername  string
	signupPassword  string
	signupAutoLogin bool
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Self-register a new user (POST /register)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if signupUsername == "" || signupPassword == "" {
			return fmt.Errorf("--username and --password are required")
		}
		cl := client.New(cfg)
		payload := types.UserSignup{
			Username: signupUsername,
			Password: signupPassword,
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
			token, exp, err := cl.Login(signupUsername, signupPassword)
			if err != nil {
				return fmt.Errorf("auto-login failed: %w", err)
			}
			if err := cl.SaveToken(token, exp); err != nil {
				return fmt.Errorf("saving token failed: %w", err)
			}
			fmt.Println("Logged in and token cached at:", cfg.TokenFile)
		}
		return nil
	},
}

func init() {
	// users
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersCreateCmd)
	usersCmd.AddCommand(usersDeleteCmd)

	usersCreateCmd.Flags().StringVar(&uCreateUsername, "username", "", "Username")
	usersCreateCmd.Flags().StringVar(&uCreatePassword, "password", "", "Password")
	usersCreateCmd.Flags().BoolVar(&uCreateIsAdmin, "admin", false, "Set user as admin")

	usersDeleteCmd.Flags().StringVar(&uDeleteUsername, "username", "", "Username to delete")

	// signup
	signupCmd.Flags().StringVar(&signupUsername, "username", "", "Username")
	signupCmd.Flags().StringVar(&signupPassword, "password", "", "Password")
	signupCmd.Flags().BoolVar(&signupAutoLogin, "login", false, "Automatically log in and cache token after signup")
}
