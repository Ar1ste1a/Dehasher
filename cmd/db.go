package cmd

import (
	"Dehash/internal/export"
	"Dehash/internal/files"
	"Dehash/internal/sqlite"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strings"
)

var (
	// DB command flags
	dbPath string

	// DB query command flags
	usernameDBQuery              string
	emailDBQuery                 string
	ipDBQuery                    string
	passwordDBQuery              string
	hashDBQuery                  string
	nameDBQuery                  string
	vinDBQuery                   string
	licensePlateDBQuery          string
	addressDBQuery               string
	phoneDBQuery                 string
	socialDBQuery                string
	cryptoCurrencyAddressDBQuery string
	domainDBQuery                string
	limitResultsDB               int
	exactMatchDBQuery            bool
	outputFormatDB               string

	// DB command
	dbCmd = &cobra.Command{
		Use:   "db",
		Short: "Database operations for Dehasher",
		Long:  `Perform database operations like export, import, and query on the local Dehasher database.`,
	}
)

func init() {
	// Add subcommands to db command
	dbCmd.AddCommand(dbExportCmd)
	dbCmd.AddCommand(dbQueryCmd)

	// Add flags specific to db command
	dbCmd.PersistentFlags().StringVarP(&dbPath, "db-path", "D", "", "Path to database (default: ~/.local/share/Dehasher/dehashed.db)")

	// Add flags specific to db query command
	dbQueryCmd.Flags().StringVarP(&usernameDBQuery, "username", "u", "", "Filter by username")
	dbQueryCmd.Flags().StringVarP(&emailDBQuery, "email", "e", "", "Filter by email")
	dbQueryCmd.Flags().StringVarP(&ipDBQuery, "ip", "i", "", "Filter by IP address")
	dbQueryCmd.Flags().StringVarP(&passwordDBQuery, "password", "p", "", "Filter by password")
	dbQueryCmd.Flags().StringVarP(&hashDBQuery, "hash", "H", "", "Filter by hashed password")
	dbQueryCmd.Flags().StringVarP(&nameDBQuery, "name", "n", "", "Filter by name")
	dbQueryCmd.Flags().StringVarP(&vinDBQuery, "vin", "v", "", "Filter by VIN")
	dbQueryCmd.Flags().StringVarP(&licensePlateDBQuery, "license", "L", "", "Filter by license plate")
	dbQueryCmd.Flags().StringVarP(&addressDBQuery, "address", "a", "", "Filter by address")
	dbQueryCmd.Flags().StringVarP(&phoneDBQuery, "phone", "P", "", "Filter by phone number")
	dbQueryCmd.Flags().StringVarP(&socialDBQuery, "social", "s", "", "Filter by social media handle")
	dbQueryCmd.Flags().StringVarP(&cryptoCurrencyAddressDBQuery, "crypto", "c", "", "Filter by cryptocurrency address")
	dbQueryCmd.Flags().StringVarP(&domainDBQuery, "domain", "d", "", "Filter by domain/URL")
	dbQueryCmd.Flags().IntVarP(&limitResultsDB, "limit", "l", 100, "Limit number of results")
	dbQueryCmd.Flags().BoolVarP(&exactMatchDBQuery, "exact", "x", false, "Use exact matching instead of partial matching")
	dbQueryCmd.Flags().StringVarP(&outputFormatDB, "format", "f", "table", "Output format (json, table, simple)")
}

// DB export command
var dbExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export database to file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Exporting database...")
		// Create DBOptions with the provided parameters
		options := &sqlite.DBOptions{
			Username:       usernameDBQuery,
			Email:          emailDBQuery,
			IPAddress:      ipDBQuery,
			Password:       passwordDBQuery,
			HashedPassword: hashDBQuery,
			Name:           nameDBQuery,
			Limit:          limitResultsDB,
			ExactMatch:     exactMatchDBQuery,
		}

		// Check if at least one search parameter is provided
		if options.Username == "" && options.Email == "" && options.IPAddress == "" &&
			options.Password == "" && options.HashedPassword == "" && options.Name == "" {
			fmt.Println("Error: At least one search parameter is required.")
			cmd.Help()
			return
		}

		// Get the count of matching results
		count, err := sqlite.GetResultsCount(options)
		if err != nil {
			fmt.Printf("Error counting results: %v\n", err)
			return
		}

		// Query the database
		results, err := sqlite.QueryResults(options)
		if err != nil {
			fmt.Printf("Error querying database: %v\n", err)
			return
		}
		dhResults := sqlite.DehashedResults{Results: results}

		fmt.Printf("Found %d results (showing %d):\n", count, len(results))

		// Output results based on format
		ft := files.GetFileType(outputFormatDB)
		err = export.WriteToFile(dhResults, "dehasher_export", ft)
		if err != nil {
			zap.L().Error("write_to_file",
				zap.String("message", "failed to write to file"),
				zap.Error(err),
			)
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
		fmt.Printf("Exported successfully to file: dehasher_export%s\n", ft.Extension())
	},
}

// DB query command
var dbQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query local database",
	Long:  `Query the local database for previously run dehasher queries based on various parameters.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create DBOptions with the provided parameters
		options := &sqlite.DBOptions{
			Username:              usernameDBQuery,
			Email:                 emailDBQuery,
			IPAddress:             ipDBQuery,
			Password:              passwordDBQuery,
			HashedPassword:        hashDBQuery,
			Name:                  nameDBQuery,
			Vin:                   vinDBQuery,
			LicensePlate:          licensePlateDBQuery,
			Address:               addressDBQuery,
			Phone:                 phoneDBQuery,
			Social:                socialDBQuery,
			CryptoCurrencyAddress: cryptoCurrencyAddressDBQuery,
			Domain:                domainDBQuery,
			Limit:                 limitResultsDB,
			ExactMatch:            exactMatchDBQuery,
		}

		// Check if at least one search parameter is provided
		if options.Username == "" && options.Email == "" && options.IPAddress == "" &&
			options.Password == "" && options.HashedPassword == "" && options.Name == "" &&
			options.Vin == "" && options.LicensePlate == "" && options.Address == "" &&
			options.Phone == "" && options.Social == "" && options.CryptoCurrencyAddress == "" && options.Domain == "" {
			fmt.Println("Error: At least one search parameter is required.")
			cmd.Help()
			return
		}

		// Get the count of matching results
		count, err := sqlite.GetResultsCount(options)
		if err != nil {
			fmt.Printf("Error counting results: %v\n", err)
			return
		}

		// Query the database
		results, err := sqlite.QueryResults(options)
		if err != nil {
			fmt.Printf("Error querying database: %v\n", err)
			return
		}

		// Display the results
		fmt.Printf("Found %d results (showing %d):\n", count, len(results))

		if len(results) == 0 {
			fmt.Println("No results found.")
			return
		}

		// Output results based on format
		switch outputFormatDB {
		case "json":
			data, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting results: %v\n", err)
				return
			}
			fmt.Println(string(data))
		case "table":
			// Print table header
			fmt.Printf("%-20s %-30s %-15s %-20s %-20s %-20s\n", "Username", "Email", "IP Address", "Password", "Hashed Password", "Name")
			fmt.Println("----------------------------------------------------------------------------------------------------")

			// Print each result
			for _, result := range results {
				fmt.Printf("%-20s %-30s %-15s %-20s %-20s %-20s\n",
					truncate(arrayToString(result.Username), 20),
					truncate(arrayToString(result.Email), 30),
					truncate(arrayToString(result.IpAddress), 15),
					truncate(arrayToString(result.Password), 20),
					truncate(arrayToString(result.HashedPassword), 20),
					truncate(arrayToString(result.Name), 20))
			}
		default:
			// Simple output
			for i, result := range results {
				fmt.Printf("Result %d:\n", i+1)
				fmt.Printf("  Username: %s\n", result.Username)
				fmt.Printf("  Email: %s\n", result.Email)
				fmt.Printf("  IP Address: %s\n", result.IpAddress)
				fmt.Printf("  Password: %s\n", result.Password)
				fmt.Printf("  Hashed Password: %s\n", result.HashedPassword)
				fmt.Printf("  Name: %s\n", result.Name)
				fmt.Println()
			}
		}
	},
}

// truncate truncates a string to the specified length and adds ellipsis if needed
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func arrayToString(a []string) string {
	return strings.Join(a, ", ")
}
