package query

import (
	"Dehash/internal/export"
	"Dehash/internal/sqlite"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
)

// Dehasher is a struct for querying the Dehashed API
type Dehasher struct {
	options  sqlite.QueryOptions
	nextPage int
	request  *DehashedSearchRequest
	client   *DehashedClientV2
}

// NewDehasher creates a new Dehasher
func NewDehasher(options *sqlite.QueryOptions) *Dehasher {
	dh := &Dehasher{
		options:  *options,
		nextPage: options.StartingPage + 1,
	}
	dh.setQueries()
	dh.request = NewDehashedSearchRequest(dh.options.StartingPage, dh.options.MaxRecords, dh.options.WildcardMatch, dh.options.RegexMatch, false)
	dh.buildRequest()
	return dh
}

// SetClientCredentials sets the client credentials for the dehasher
func (dh *Dehasher) SetClientCredentials(key string) {
	dh.client = NewDehashedClientV2(key)
}

func (dh *Dehasher) getNextPage() int {
	nextPage := dh.nextPage
	dh.nextPage += 1
	return nextPage
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
	fmt.Printf("Making %d Requests for %d Records (%d Total)\n", dh.options.MaxRequests, dh.options.MaxRecords, dh.options.MaxRequests*dh.options.MaxRecords)
}

// Start starts the querying process
func (dh *Dehasher) Start() {
	fmt.Println("[*] Querying Dehashed API...")
	for i := 0; i < dh.options.MaxRequests; i++ {
		fmt.Printf("\n\t[*] Performing Request...")
		count, err := dh.client.Search(*dh.request)
		if err != nil {
			fmt.Printf("[!] Error performing request: %v", err)
			os.Exit(-1)
		}

		if count < dh.options.MaxRecords {
			fmt.Printf("\n\t\t[+] Retrieved %d Records", count)
			fmt.Printf("\n[-] Not Enough Entries, ending queries")
			break
		} else {
			fmt.Printf("\n\t\t[+] Retrieved %d Records", dh.options.MaxRecords)
		}

		dh.request.Page = dh.getNextPage()
	}

	dh.parseResults()
}

// buildRequest constructs the query map
func (dh *Dehasher) buildRequest() {
	if len(dh.options.UsernameQuery) > 0 {
		dh.request.AddUsernameQuery(dh.options.UsernameQuery)
	}
	if len(dh.options.EmailQuery) > 0 {
		dh.request.AddEmailQuery(dh.options.EmailQuery)
	}
	if len(dh.options.IpQuery) > 0 {
		dh.request.AddIpAddressQuery(dh.options.IpQuery)
	}
	if len(dh.options.HashQuery) > 0 {
		dh.request.AddHashedPasswordQuery(dh.options.HashQuery)
	}
	if len(dh.options.PassQuery) > 0 {
		dh.request.AddPasswordQuery(dh.options.PassQuery)
	}
	if len(dh.options.NameQuery) > 0 {
		dh.request.AddNameQuery(dh.options.NameQuery)
	}
	if len(dh.options.DomainQuery) > 0 {
		dh.request.AddDomainQuery(dh.options.DomainQuery)
	}
	if len(dh.options.VinQuery) > 0 {
		dh.request.AddVinQuery(dh.options.VinQuery)
	}
	if len(dh.options.LicensePlateQuery) > 0 {
		dh.request.AddLicensePlateQuery(dh.options.LicensePlateQuery)
	}
	if len(dh.options.AddressQuery) > 0 {
		dh.request.AddAddressQuery(dh.options.AddressQuery)
	}
	if len(dh.options.PhoneQuery) > 0 {
		dh.request.AddPhoneQuery(dh.options.PhoneQuery)
	}
	if len(dh.options.SocialQuery) > 0 {
		dh.request.AddSocialQuery(dh.options.SocialQuery)
	}
	if len(dh.options.CryptoAddressQuery) > 0 {
		dh.request.AddCryptoAddressQuery(dh.options.CryptoAddressQuery)
	}
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
