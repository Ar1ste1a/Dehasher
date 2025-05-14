package query

import (
	"Dehash/internal/export"
	"Dehash/internal/sqlite"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
)

// Dehasher is a struct for querying the Dehashed API
type Dehasher struct {
	options        sqlite.QueryOptions
	username       string
	email          string
	ipAddress      string
	password       string
	hashedPassword string
	name           string
	query          string
	page           int
	params         map[string]string
	client         *DehashedClient
}

// NewDehasher creates a new Dehasher
func NewDehasher(options *sqlite.QueryOptions) *Dehasher {
	dh := &Dehasher{
		options: *options,
	}
	dh.escapeReservedCharacters()
	dh.constructMap()
	dh.setQueries()
	if len(dh.params) == 0 {
		fmt.Println("At least one return type is required")
		os.Exit(-1)
	}

	return dh
}

// escapeReservedCharacters escapes reserved characters in the query
func (dh *Dehasher) escapeReservedCharacters() {
	reserved := strings.Split("+ - = && || > < ! ( ) { } [ ] ^ \" ~ * ? : \\", " ")

	dh.username = escapeString(dh.options.UsernameQuery, reserved)
	dh.email = escapeString(dh.options.EmailQuery, reserved)
	dh.ipAddress = escapeString(dh.options.IpQuery, reserved)
	dh.password = escapeString(dh.options.PassQuery, reserved)
	dh.hashedPassword = escapeString(dh.options.HashQuery, reserved)
	dh.name = escapeString(dh.options.NameQuery, reserved)
}

// SetClientCredentials sets the client credentials for the dehasher
func (dh *Dehasher) SetClientCredentials(key, email string, printBal bool) {
	dh.client = NewDehashedClient(key, email, printBal)
}

// setQueries sets the number of queries to make based on the number of records and requests
func (dh *Dehasher) setQueries() {
	var numQueries int

	switch {
	case dh.options.MaxRequests == 0:
		zap.L().Error("max requests cannot be zero")
		fmt.Println("[!] Max Requests cannot be zero")
		os.Exit(1)
	case dh.options.MaxRecords <= 10000 || dh.options.MaxRequests == 1:
		numQueries = 1
		if dh.options.MaxRecords > 10000 {
			dh.options.MaxRecords = 10000
		}
		zap.L().Info("max requests set to 1", zap.Int("max_records", dh.options.MaxRecords))
	case dh.options.MaxRequests < 0 && dh.options.MaxRecords > 20000:
		numQueries = 3
		dh.options.MaxRecords = 10000
		zap.L().Info("max requests set to 3", zap.Int("max_records", dh.options.MaxRecords))
	case dh.options.MaxRequests < 0 && dh.options.MaxRecords > 10000:
		numQueries = 2
		dh.options.MaxRecords = 10000
		zap.L().Info("max requests set to 2", zap.Int("max_records", dh.options.MaxRecords))
	case dh.options.MaxRecords < 0 && dh.options.MaxRecords < 10000:
		numQueries = 1
		zap.L().Info("max requests set to 1", zap.Int("max_records", dh.options.MaxRecords))
	case dh.options.MaxRequests == 2 && dh.options.MaxRecords > 20000:
		numQueries = 2
		dh.options.MaxRecords = 10000
		zap.L().Info("max requests set to 2", zap.Int("max_records", dh.options.MaxRecords))
	case dh.options.MaxRequests == 2 && dh.options.MaxRecords <= 10000:
		numQueries = 1
		zap.L().Info("max requests set to 1", zap.Int("max_records", dh.options.MaxRecords))
	default:
		numQueries = 3
		dh.options.MaxRecords = 10000
		zap.L().Info("max requests set to 3", zap.Int("max_records", dh.options.MaxRecords))
	}

	dh.options.MaxRequests = numQueries
	fmt.Printf("Making %d Requests for %d Records (%d Total)", dh.options.MaxRequests, dh.options.MaxRecords, dh.options.MaxRequests*dh.options.MaxRecords)
}

// Start starts the querying process
func (dh *Dehasher) Start() {
	dh.client.buildQuery(dh.params)
	page := 1
	offset := 0
	foundNum := 0
	for i := 0; i < dh.options.MaxRequests; i++ {
		dh.client.setResults(dh.options.MaxRecords)
		dh.client.setPage(page)
		found := dh.client.Do()

		if found-offset < dh.options.MaxRecords {
			fmt.Printf("\n\t\t[+] Retrieved %d Records", found-offset)
			fmt.Printf("\n[-] Not Enough Entries, ending queries")
			break
		} else {
			if found-offset > dh.options.MaxRecords {
				foundNum = dh.options.MaxRecords
			} else {
				foundNum = found - offset
			}
			fmt.Printf("\n\t\t[*] Retrieved %d Records", foundNum)
			offset += dh.options.MaxRecords
			page += 1
		}
	}

	dh.parseResults()
}

// buildQuery builds the query string
func (dh *Dehasher) buildQuery() {
	for key, value := range dh.params {
		println(key, value)
	}
}

// constructMap constructs the query map
func (dh *Dehasher) constructMap() {
	urlParams := map[string]string{}

	if len(dh.username) > 0 {
		urlParams["username"] = dh.username
	}
	if len(dh.email) > 0 {
		urlParams["email"] = dh.email
	}
	if len(dh.ipAddress) > 0 {
		urlParams["ip_address"] = dh.ipAddress
	}
	if len(dh.hashedPassword) > 0 {
		urlParams["hashed_password"] = dh.hashedPassword
	}
	if len(dh.password) > 0 {
		urlParams["password"] = dh.password
	}
	if len(dh.name) > 0 {
		urlParams["name"] = dh.name
	}

	dh.params = urlParams
}

// escapeString escapes reserved characters in the query
func escapeString(input string, reserved []string) string {
	for _, char := range reserved {
		input = strings.Replace(input, char, "\\"+char, -1)
	}
	return input
}

// parseResults parses the results and writes them to a file
func (dh *Dehasher) parseResults() {
	var data []byte
	results := dh.client.GetResults()

	defer func() {
		zap.L().Info("storing_results")
		err := sqlite.StoreResults(results)
		if err != nil {
			zap.L().Error("store_results",
				zap.String("message", "failed to store results"),
				zap.Error(err),
			)
		}
	}()

	if len(results.Results) > 0 {
		fmt.Printf("\n\t[*] Writing entries to file: %s.%s", dh.options.OutputFile, dh.options.OutputFormat.String())
		if !dh.options.CredsOnly {
			err := export.WriteToFile(results, dh.options.OutputFile, dh.options.OutputFormat)
			if err != nil {
				fmt.Printf("\n[!] Error Writing to file: %v\n\tOutputting to terminal.", err)
				data, err = json.MarshalIndent(results, "", "  ")
				fmt.Println(string(data))
				os.Exit(0)
			} else {
				fmt.Println("\n\t\t[*] Success\n")
				os.Exit(1)
			}
		} else {
			creds := results.ExtractCredentials()
			err := export.WriteCredsToFile(creds, dh.options.OutputFile, dh.options.OutputFormat)
			if err != nil {
				fmt.Printf("\n[!] Error Writing to file: %v\n\tOutputting to terminal.", err)
				data, err = json.MarshalIndent(creds, "", "  ")
				fmt.Println(string(data))
				os.Exit(0)
			} else {
				fmt.Println("\n\t\t[*] Success\n")
				os.Exit(1)
			}
		}
	}
}
