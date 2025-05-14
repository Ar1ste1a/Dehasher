package cmd

import (
	"Dehash/internal/whois"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strings"
)

var (
	// WHOIS command flags
	whoisDomain       string
	whoisIPAddress    string
	whoisMXAddress    string
	whoisNSAddress    string
	whoisInclude      string
	whoisExclude      string
	whoisReverseType  string
	whoisOutputFormat string
	whoisShowCredits  bool

	// WHOIS command
	whoisCmd = &cobra.Command{
		Use:   "whois",
		Short: "Dehashed WHOIS lookups and reverse WHOIS searches",
		Long:  `Perform WHOIS lookups, history searches, reverse WHOIS searches, IP lookups, MX lookups, NS lookups, and subdomain scans.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if API key is provided
			key := apiKey

			// If not provided as flag, try to get from stored value
			if key == "" {
				key = getStoredApiKey()
			}

			// Validate credentials
			if key == "" {
				fmt.Println("API key is required. Use --key flag or set it with set-key command.")
				return
			}

			// Show credits if requested
			if whoisShowCredits {
				credits, err := whois.GetWHOISCredits(key)
				if err != nil {
					zap.L().Error("get_whois_credits",
						zap.String("message", "failed to get whois credits"),
						zap.Error(err),
					)
					fmt.Printf("Error getting WHOIS credits: %v\n", err)
					return
				}
				fmt.Println("WHOIS Credits:")
				fmt.Println(credits)
				return
			}

			// Determine which operation to perform based on flags
			if whoisDomain != "" {
				// Domain lookup
				result, err := whois.WhoisSearch(whoisDomain, key)
				if err != nil {
					zap.L().Error("whois_search",
						zap.String("message", "failed to perform whois search"),
						zap.Error(err),
					)
					fmt.Printf("Error performing WHOIS lookup: %v\n", err)
					return
				}
				fmt.Println("WHOIS Lookup Result:")
				fmt.Println(result)

				// Also perform history search
				history, err := whois.WhoisHistory(whoisDomain, key)
				if err != nil {
					zap.L().Error("whois_history",
						zap.String("message", "failed to perform whois history lookup"),
						zap.Error(err),
					)
					fmt.Printf("Error performing WHOIS history lookup: %v\n", err)
				} else {
					fmt.Println("\nWHOIS History:")
					fmt.Println(history)
				}

				// Also perform subdomain scan
				subdomains, err := whois.WhoisSubdomainScan(whoisDomain, key)
				if err != nil {
					zap.L().Error("whois_subdomain_scan",
						zap.String("message", "failed to perform subdomain scan"),
						zap.Error(err),
					)
					fmt.Printf("Error performing subdomain scan: %v\n", err)
				} else {
					fmt.Println("\nSubdomain Scan:")
					fmt.Println(subdomains)
				}
				return
			}

			if whoisIPAddress != "" {
				// IP lookup
				result, err := whois.WhoisIP(whoisIPAddress, key)
				if err != nil {
					zap.L().Error("whois_ip",
						zap.String("message", "failed to perform ip lookup"),
						zap.Error(err),
					)
					fmt.Printf("Error performing IP lookup: %v\n", err)
					return
				}
				fmt.Println("IP Lookup Result:")
				fmt.Println(string(result))
				return
			}

			if whoisMXAddress != "" {
				// MX lookup
				result, err := whois.WhoisMX(whoisMXAddress, key)
				if err != nil {
					zap.L().Error("whois_mx",
						zap.String("message", "failed to perform mx lookup"),
						zap.Error(err),
					)
					fmt.Printf("Error performing MX lookup: %v\n", err)
					return
				}
				fmt.Println("MX Lookup Result:")
				fmt.Println(result)
				return
			}

			if whoisNSAddress != "" {
				// NS lookup
				result, err := whois.WhoisNS(whoisNSAddress, key)
				if err != nil {
					zap.L().Error("whois_ns",
						zap.String("message", "failed to perform ns lookup"),
						zap.Error(err),
					)
					fmt.Printf("Error performing NS lookup: %v\n", err)
					return
				}
				fmt.Println("NS Lookup Result:")
				fmt.Println(result)
				return
			}

			if whoisInclude != "" || whoisExclude != "" {
				// Reverse WHOIS
				includeTerms := []string{}
				if whoisInclude != "" {
					includeTerms = strings.Split(whoisInclude, ",")
				}

				excludeTerms := []string{}
				if whoisExclude != "" {
					excludeTerms = strings.Split(whoisExclude, ",")
				}

				if whoisReverseType == "" {
					whoisReverseType = "registrant"
				}

				result, err := whois.ReverseWHOIS(includeTerms, excludeTerms, whoisReverseType, key)
				if err != nil {
					fmt.Printf("Error performing reverse WHOIS: %v\n", err)
					return
				}
				fmt.Println("Reverse WHOIS Result:")
				fmt.Println(result)
				return
			}

			// If no specific operation was requested
			cmd.Help()
		},
	}
)

func init() {
	// Add whois command to root command
	rootCmd.AddCommand(whoisCmd)

	// Add flags specific to whois command
	whoisCmd.Flags().StringVarP(&whoisDomain, "domain", "d", "", "Domain for WHOIS lookup, history search, and subdomain scan")
	whoisCmd.Flags().StringVarP(&whoisIPAddress, "ip", "i", "", "IP address for reverse IP lookup")
	whoisCmd.Flags().StringVarP(&whoisMXAddress, "mx", "m", "", "MX address for reverse MX lookup")
	whoisCmd.Flags().StringVarP(&whoisNSAddress, "ns", "n", "", "NS address for reverse NS lookup")
	whoisCmd.Flags().StringVarP(&whoisInclude, "include", "I", "", "Terms to include in reverse WHOIS search (comma-separated)")
	whoisCmd.Flags().StringVarP(&whoisExclude, "exclude", "E", "", "Terms to exclude in reverse WHOIS search (comma-separated)")
	whoisCmd.Flags().StringVarP(&whoisReverseType, "type", "t", "registrant", "Type of reverse WHOIS search (registrant, email, organization, address, phone)")
	whoisCmd.Flags().StringVarP(&whoisOutputFormat, "format", "f", "text", "Output format (text, json)")
	whoisCmd.Flags().BoolVarP(&whoisShowCredits, "credits", "c", false, "Show remaining WHOIS credits")

	// Add API key flag
	whoisCmd.Flags().StringVarP(&apiKey, "key", "k", "", "Dehashed API key")
}
