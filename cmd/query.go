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
	outputFormat               string
	outputFile                 string
	regexMatch                 bool
	wildcardMatch              bool
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
	credsOnly                  bool
	printBalance               bool

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
				apiKey,
				apiEmail,
				printBalance,
			)

			// Start querying
			fmt.Println("[*] Querying Dehashed API...")
			dehasher.Start()
			fmt.Println("\n[*] Completing Process")
		},
	}
)

func init() {
	// Add flags specific to query command
	queryCmd.Flags().IntVarP(&maxRecords, "max-records", "m", 30000, "Maximum amount of records to return")
	queryCmd.Flags().IntVarP(&maxRequests, "max-requests", "r", -1, "Maximum number of requests to make")
	queryCmd.Flags().BoolVarP(&printBalance, "print-balance", "B", false, "Print remaining balance after requests")
	queryCmd.Flags().BoolVarP(&regexMatch, "regex-match", "R", false, "Use regex matching on fields (u=username, e=email, i=ip, p=password, q=hash, n=name)")
	queryCmd.Flags().BoolVarP(&wildcardMatch, "wildcard-match", "W", false, "Use wildcard matching on fields (u=username, e=email, i=ip, p=password, q=hash, n=name)")
	queryCmd.Flags().BoolVarP(&credsOnly, "creds-only", "C", false, "Return credentials only")
	queryCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format (json, yaml, xml, txt)")
	queryCmd.Flags().StringVarP(&outputFile, "output", "o", "query", "File to output results to including extension")
	queryCmd.Flags().StringVarP(&usernameQuery, "username", "U", "", "Username query")
	queryCmd.Flags().StringVarP(&emailQuery, "email-query", "E", "", "Email query")
	queryCmd.Flags().StringVarP(&ipQuery, "ip", "I", "", "IP address query")
	queryCmd.Flags().StringVarP(&domainQuery, "domain", "D", "", "Domain query")
	queryCmd.Flags().StringVarP(&passwordQuery, "password", "P", "", "Password query")
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
