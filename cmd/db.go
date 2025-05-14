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
	nonEmptyFieldsDBQuery        string
	displayFieldsDBQuery         string

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
	dbQueryCmd.Flags().StringVar(&nonEmptyFieldsDBQuery, "non-empty", "", "Filter for non-empty fields (comma-separated list, e.g., 'password,email')")
	dbQueryCmd.Flags().StringVar(&displayFieldsDBQuery, "display", "", "Fields to display in output (comma-separated list, e.g., 'username,email,password')")
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

		// Parse non-empty fields if provided
		if nonEmptyFieldsDBQuery != "" {
			options.NonEmptyFields = strings.Split(nonEmptyFieldsDBQuery, ",")
		}

		// Parse display fields if provided
		if displayFieldsDBQuery != "" {
			options.DisplayFields = strings.Split(displayFieldsDBQuery, ",")
		}

		// Check if at least one search parameter is provided
		if options.Username == "" && options.Email == "" && options.IPAddress == "" &&
			options.Password == "" && options.HashedPassword == "" && options.Name == "" &&
			len(options.NonEmptyFields) == 0 {
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

		// Parse non-empty fields if provided
		if nonEmptyFieldsDBQuery != "" {
			options.NonEmptyFields = strings.Split(nonEmptyFieldsDBQuery, ",")
		}

		// Parse display fields if provided
		if displayFieldsDBQuery != "" {
			options.DisplayFields = strings.Split(displayFieldsDBQuery, ",")
		}

		// Check if at least one search parameter is provided
		if options.Username == "" && options.Email == "" && options.IPAddress == "" &&
			options.Password == "" && options.HashedPassword == "" && options.Name == "" &&
			options.Vin == "" && options.LicensePlate == "" && options.Address == "" &&
			options.Phone == "" && options.Social == "" && options.CryptoCurrencyAddress == "" && options.Domain == "" &&
			len(options.NonEmptyFields) == 0 {
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
			// Determine which fields to display
			type FieldInfo struct {
				Name   string
				Width  int
				Getter func(result sqlite.Result) string
			}

			// Define all available fields
			allFields := []FieldInfo{
				{"Username", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.Username), 20) }},
				{"Email", 30, func(r sqlite.Result) string { return truncate(arrayToString(r.Email), 30) }},
				{"IP Address", 15, func(r sqlite.Result) string { return truncate(arrayToString(r.IpAddress), 15) }},
				{"Password", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.Password), 20) }},
				{"Hashed Password", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.HashedPassword), 20) }},
				{"Name", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.Name), 20) }},
				{"VIN", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.Vin), 20) }},
				{"License Plate", 15, func(r sqlite.Result) string { return truncate(arrayToString(r.LicensePlate), 15) }},
				{"Address", 30, func(r sqlite.Result) string { return truncate(arrayToString(r.Address), 30) }},
				{"Phone", 15, func(r sqlite.Result) string { return truncate(arrayToString(r.Phone), 15) }},
				{"Social", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.Social), 20) }},
				{"Crypto Address", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.CryptoCurrencyAddress), 20) }},
				{"Domain/URL", 30, func(r sqlite.Result) string { return truncate(arrayToString(r.Url), 30) }},
			}

			// Select fields to display
			var fieldsToDisplay []FieldInfo
			if len(options.DisplayFields) > 0 {
				// Use specified fields
				for _, fieldName := range options.DisplayFields {
					fieldName = strings.ToLower(strings.TrimSpace(fieldName))
					for _, field := range allFields {
						if strings.ToLower(field.Name) == fieldName ||
							(fieldName == "ip" && strings.ToLower(field.Name) == "ip address") ||
							(fieldName == "hash" && strings.ToLower(field.Name) == "hashed password") ||
							(fieldName == "license" && strings.ToLower(field.Name) == "license plate") ||
							(fieldName == "crypto" && strings.ToLower(field.Name) == "crypto address") ||
							(fieldName == "url" && strings.ToLower(field.Name) == "domain/url") {
							fieldsToDisplay = append(fieldsToDisplay, field)
							break
						}
					}
				}
			} else {
				// Default fields (first 6)
				fieldsToDisplay = allFields[:6]
			}

			// Print table header
			formatStr := ""
			headerValues := []interface{}{}
			for _, field := range fieldsToDisplay {
				formatStr += "%-" + fmt.Sprintf("%d", field.Width) + "s "
				headerValues = append(headerValues, field.Name)
			}
			fmt.Printf(formatStr+"\n", headerValues...)

			// Print separator line
			separator := ""
			for _, field := range fieldsToDisplay {
				separator += strings.Repeat("-", field.Width) + " "
			}
			fmt.Println(separator)

			// Print each result
			for _, result := range results {
				rowValues := []interface{}{}
				for _, field := range fieldsToDisplay {
					rowValues = append(rowValues, field.Getter(result))
				}
				fmt.Printf(formatStr+"\n", rowValues...)
			}
		default:
			// Simple output
			for i, result := range results {
				fmt.Printf("Result %d:\n", i+1)

				// Determine which fields to display
				if len(options.DisplayFields) > 0 {
					// Display only specified fields
					for _, field := range options.DisplayFields {
						field = strings.ToLower(strings.TrimSpace(field))
						switch field {
						case "username":
							fmt.Printf("  Username: %s\n", result.Username)
						case "email":
							fmt.Printf("  Email: %s\n", result.Email)
						case "ip", "ipaddress", "ip_address":
							fmt.Printf("  IP Address: %s\n", result.IpAddress)
						case "password":
							fmt.Printf("  Password: %s\n", result.Password)
						case "hash", "hashed_password":
							fmt.Printf("  Hashed Password: %s\n", result.HashedPassword)
						case "name":
							fmt.Printf("  Name: %s\n", result.Name)
						case "vin":
							fmt.Printf("  VIN: %s\n", result.Vin)
						case "license", "license_plate":
							fmt.Printf("  License Plate: %s\n", result.LicensePlate)
						case "address":
							fmt.Printf("  Address: %s\n", result.Address)
						case "phone":
							fmt.Printf("  Phone: %s\n", result.Phone)
						case "social":
							fmt.Printf("  Social: %s\n", result.Social)
						case "crypto", "cryptocurrency_address":
							fmt.Printf("  Crypto Address: %s\n", result.CryptoCurrencyAddress)
						case "domain", "url":
							fmt.Printf("  Domain/URL: %s\n", result.Url)
						}
					}
				} else {
					// Display default fields
					fmt.Printf("  Username: %s\n", result.Username)
					fmt.Printf("  Email: %s\n", result.Email)
					fmt.Printf("  IP Address: %s\n", result.IpAddress)
					fmt.Printf("  Password: %s\n", result.Password)
					fmt.Printf("  Hashed Password: %s\n", result.HashedPassword)
					fmt.Printf("  Name: %s\n", result.Name)
				}
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
