package cmd

import (
	"Dehash/internal/badger"
	"Dehash/internal/query"
	"Dehash/internal/sqlite"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	// Query command flags
	maxRecords                 int
	maxRequests                int
	startingPage               int
	credsOnly                  bool
	printBalance               bool
	regexMatch                 bool
	wildcardMatch              bool
	outputFormat               string
	outputFile                 string
	usernameQuery              string
	emailQuery                 string
	ipQuery                    string
	passwordQuery              string
	hashQuery                  string
	nameQuery                  string
	domainQuery                string
	vinQuery                   string
	licensePlateQuery          string
	addressQuery               string
	phoneQuery                 string
	socialQuery                string
	cryptoCurrencyAddressQuery string

	// Query command
	queryCmd = &cobra.Command{
		Use:   "query",
		Short: "Query the Dehashed API",
		Long:  `Query the Dehashed API for emails, usernames, passwords, hashes, IP addresses, and names.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if API key and email are provided
			key := apiKey
			email := apiEmail

			// If not provided as flags, try to get from stored values
			if key == "" {
				key = getStoredApiKey()
			}
			if email == "" {
				email = getStoredApiEmail()
			}

			// Validate credentials
			if key == "" || email == "" {
				fmt.Println("API key and email are required. Use --key and --email flags or set them with set-key and set-email commands.")
				return
			}

			// Create new QueryOptions
			queryOptions := sqlite.NewQueryOptions(
				maxRecords,
				maxRequests,
				startingPage,
				outputFormat,
				outputFile,
				usernameQuery,
				emailQuery,
				ipQuery,
				passwordQuery,
				hashQuery,
				nameQuery,
				domainQuery,
				vinQuery,
				licensePlateQuery,
				addressQuery,
				phoneQuery,
				socialQuery,
				cryptoCurrencyAddressQuery,
				regexMatch,
				wildcardMatch,
				printBalance,
				credsOnly,
			)

			// Create new Dehasher
			dehasher := query.NewDehasher(queryOptions)
			dehasher.SetClientCredentials(
				key,
			)

			// Start querying
			dehasher.Start()
			fmt.Println("\n[*] Completing Process")
		},
	}
)

func init() {
	// Add flags specific to query command
	queryCmd.Flags().IntVarP(&maxRecords, "max-records", "m", 30000, "Maximum amount of records to return")
	queryCmd.Flags().IntVarP(&maxRequests, "max-requests", "r", -1, "Maximum number of requests to make")
	queryCmd.Flags().IntVarP(&startingPage, "starting-page", "s", 1, "Starting page for requests")
	queryCmd.Flags().BoolVarP(&printBalance, "print-balance", "b", false, "Print remaining balance after requests")
	queryCmd.Flags().BoolVarP(&regexMatch, "regex-match", "R", false, "Use regex matching on query fields")
	queryCmd.Flags().BoolVarP(&wildcardMatch, "wildcard-match", "W", false, "Use wildcard matching on query fields (Use ? to replace a single character, and * for multiple characters)")
	queryCmd.Flags().BoolVarP(&credsOnly, "creds-only", "C", false, "Return credentials only")
	queryCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format (json, yaml, xml, txt)")
	queryCmd.Flags().StringVarP(&outputFile, "output", "o", "query", "File to output results to including extension")
	queryCmd.Flags().StringVarP(&usernameQuery, "username", "U", "", "Username query")
	queryCmd.Flags().StringVarP(&emailQuery, "email-query", "E", "", "Email query")
	queryCmd.Flags().StringVarP(&ipQuery, "ip", "I", "", "IP address query")
	queryCmd.Flags().StringVarP(&domainQuery, "domain", "D", "", "Domain query")
	queryCmd.Flags().StringVarP(&passwordQuery, "password", "P", "", "Password query")
	queryCmd.Flags().StringVarP(&vinQuery, "vin", "V", "", "VIN query")
	queryCmd.Flags().StringVarP(&licensePlateQuery, "license", "L", "", "License plate query")
	queryCmd.Flags().StringVarP(&addressQuery, "address", "A", "", "Address query")
	queryCmd.Flags().StringVarP(&phoneQuery, "phone", "M", "", "Phone query")
	queryCmd.Flags().StringVarP(&socialQuery, "social", "S", "", "Social query")
	queryCmd.Flags().StringVarP(&cryptoCurrencyAddressQuery, "crypto", "B", "", "Crypto currency address query")
	queryCmd.Flags().StringVarP(&hashQuery, "hash", "Q", "", "Hashed password query")
	queryCmd.Flags().StringVarP(&nameQuery, "name", "N", "", "Name query")

	// Add mutually exclusive flags to exact match and regex match
	queryCmd.MarkFlagsMutuallyExclusive("regex-match", "wildcard-match")
}

// Helper functions to get stored API credentials
func getStoredApiKey() string {
	return badger.GetKey()
}

func getStoredApiEmail() string {
	return badger.GetEmail()
}
