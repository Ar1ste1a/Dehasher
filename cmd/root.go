package cmd

import (
	"Dehash/internal/badger"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var (
	// Global Flags
	apiKey   string
	apiEmail string

	// rootCmd is the base command for the CLI.
	rootCmd = &cobra.Command{
		Use:   "dehasher",
		Short: `Dehasher is a cli tool for querying query.`,
		Long: fmt.Sprintf(
			"%s\n%s",
			`
 ______   _______           _______  _______           _______  _______ 
(  __  \ (  ____ \|\     /|(  ___  )(  ____ \|\     /|(  ____ \(  ____ )
| (  \  )| (    \/| )   ( || (   ) || (    \/| )   ( || (    \/| (    )|
| |   ) || (__    | (___) || (___) || (_____ | (___) || (__    | (____)|
| |   | ||  __)   |  ___  ||  ___  |(_____  )|  ___  ||  __)   |     __)
| |   ) || (      | (   ) || (   ) |      ) || (   ) || (      | (\ (   
| (__/  )| (____/\| )   ( || )   ( |/\____) || )   ( || (____/\| ) \ \__
(______/ (_______/|/     \||/     \|\_______)|/     \|(_______/|/   \__/
An Ar1ste1a Project                                                                        
`,
			`––•–√\/––√\/––•––––•–√\/––√\/––•––––•–√\/––√\/––•––
  Dehasher can query the query API for:
  - Emails
  - Usernames
  - Password
  - Hashes
  - IP Addresses
  - Names
  - VINs
  - License Plates
  - Addresses
  - Phones
  - Social Media
  - Crypto Currency Addresses
  Dehasher supports:
  - Regex Matching
  - Exact Matching
––•–√\/––√\/––•––––•–√\/––√\/––•––––•–√\/––√\/––•––
`,
		),
		Version: "v1.0",
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal("execute_root_command",
			zap.String("message", "failed to execute root command"),
			zap.Error(err),
		)
		fmt.Printf("[!] %v", err)
		os.Exit(1)
	}
}

func init() {
	// Hide the default help command
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// Add global flags for API key and email
	rootCmd.PersistentFlags().StringVarP(&apiKey, "key", "k", "", "API Key for authentication")
	rootCmd.PersistentFlags().StringVarP(&apiEmail, "email", "e", "", "Email to pair with API key for authentication")

	// Add subcommands
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(setKeyCmd)
	rootCmd.AddCommand(setEmailCmd)
}

// Command to set API key
var setKeyCmd = &cobra.Command{
	Use:   "set-key [key]",
	Short: "Set and store API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		// Store key in badger DB
		err := storeApiKey(key)
		if err != nil {
			fmt.Printf("Error storing API key: %v\n", err)
			return
		}
		fmt.Println("API key stored successfully")
	},
}

// Command to set API email
var setEmailCmd = &cobra.Command{
	Use:   "set-email [email]",
	Short: "Set and store API email",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		// Store email in badger DB
		err := storeApiEmail(email)
		if err != nil {
			fmt.Printf("Error storing API email: %v\n", err)
			return
		}
		fmt.Println("API email stored successfully")
	},
}

// Helper functions to store API credentials
func storeApiKey(key string) error {
	err := badger.StoreKey(key)
	if err != nil {
		fmt.Printf("Error storing API key: %v\n", err)
		return err
	}
	return nil
}

func storeApiEmail(email string) error {
	err := badger.StoreEmail(email)
	if err != nil {
		fmt.Printf("Error storing API email: %v\n", err)
		return err
	}
	return nil
}
