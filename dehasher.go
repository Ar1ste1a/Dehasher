package main

import (
	"Dehash/internal/argparse"
	"Dehash/internal/dehashed"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func isValidRegex(s string) bool {
	if _, err := regexp.Compile(s); err != nil {
		return false
	}
	return true
}

func getParser() *argparse.Argparse {
	parser := argparse.NewParser("Dehasher", "Dehashed Tool", "v1.0")
	parser.AddArg("-k", "--key", "API Key", "string", "", true)
	parser.AddArg("-a", "--authorized-email", "Email to pair with key for authentication", "string", "", true)
	parser.AddArg("-m", "--max-records", "Maximum amount of records to return", "integer", "30000", false)
	parser.AddArg("-r", "--max-requests", "Maximum number of requests to make", "integer", "-1", false)
	parser.AddArg("-B", "--print-balance", "Print remaining balance after requests", "bool", "false", false)
	parser.AddArg("-X", "--exact-match", "Use Exact Matching on fields", "string", "", false)
	parser.AddArg("-R", "--regex-match", "Use Regex Matching on fields", "string", "", false)
	parser.AddArg("-t", "--list-tokens", "List the number of tokens remaining", "bool", "false", false)
	parser.AddArg("-o", "--output-file-name", "File to output results to", "string", "dehashed", false)
	parser.AddArg("-T", "--output-txt", "Output to text file", "bool", "false", false)
	parser.AddArg("-J", "--output-json", "Output to JSON file", "bool", "true", false)
	parser.AddArg("-Y", "--output-yaml", "Output to YAML file", "bool", "false", false)
	parser.AddArg("-x", "--output-xml", "Output to XML file", "bool", "false", false)
	parser.AddArg("-U", "--username-query", "Username Query", "string", "", false)
	parser.AddArg("-E", "--email-query", "Email Query", "string", "", false)
	parser.AddArg("-I", "--ip-address-query", "IP Address Query", "string", "", false)
	parser.AddArg("-P", "--password-query", "Password Query", "string", "", false)
	parser.AddArg("-Q", "--hashed-password-query", "Hashed Password Query", "string", "", false)
	parser.AddArg("-N", "--name-query", "Name Query", "string", "", false)
	parser.AddArg("-C", "--creds-only", "Return Credentials Only", "bool", "false", false)

	return parser
}

func startDehash(p *argparse.Argparse) {
	var (
		eUsername, eEmail, eIP, ePassword, eHashedPassword, eName string
	)

	eUsername = p.Get("username-query").(string)
	eEmail = p.Get("email-query").(string)
	eIP = p.Get("ip-address-query").(string)
	ePassword = p.Get("password-query").(string)
	eHashedPassword = p.Get("hashed-password-query").(string)
	eName = p.Get("name-query").(string)

	exactMatch := strings.ToLower(p.Get("exact-match").(string))
	if len(exactMatch) > 0 {
		// Determine which params are Exact Matches, encompass in double quotes
		for _, param := range exactMatch {
			alias := strings.ToLower(string(rune(param)))

			switch alias {
			case "u":
				if isValidRegex(eUsername) {
					tmp := fmt.Sprintf("\"%s\"", eUsername)
					eUsername = tmp
				} else {
					fmt.Printf("Invalid Regex: %s", eUsername)
					os.Exit(-1)
				}
			case "e":
				if isValidRegex(eEmail) {
					tmp := fmt.Sprintf("\"%s\"", eEmail)
					eEmail = tmp
				} else {
					fmt.Printf("Invalid Regex: %s", eEmail)
					os.Exit(-1)
				}
			case "i":
				if isValidRegex(eIP) {
					tmp := fmt.Sprintf("\"%s\"", eIP)
					eIP = tmp
				} else {
					fmt.Printf("Invalid Regex: %s", eIP)
					os.Exit(-1)
				}
			case "p":
				if isValidRegex(ePassword) {
					tmp := fmt.Sprintf("\"%s\"", ePassword)
					ePassword = tmp
				} else {
					fmt.Printf("Invalid Regex: %s", ePassword)
					os.Exit(-1)
				}
			case "q":
				if isValidRegex(eHashedPassword) {
					tmp := fmt.Sprintf("\"%s\"", eHashedPassword)
					eHashedPassword = tmp
				} else {
					fmt.Printf("Invalid Regex: %s", eHashedPassword)
					os.Exit(-1)
				}
			case "n":
				if isValidRegex(eName) {
					tmp := fmt.Sprintf("\"%s\"", eName)
					eName = tmp
				} else {
					fmt.Printf("Invalid Regex: %s", eName)
					os.Exit(-1)
				}
			default:
				fmt.Printf("\nUnknown parameter set for 'Regular Expression match': %s\n", alias)
				os.Exit(-1)
			}
		}
	}

	regexMatch := strings.ToLower(p.Get("regex-match").(string))
	if len(regexMatch) > 0 {
		// parse the regex match flags, encompass in forward slashes
		for _, param := range regexMatch {
			alias := string(rune(param))

			switch alias {
			case "u":
				tmp := fmt.Sprintf("/%s/", eUsername)
				eUsername = tmp
			case "e":
				tmp := fmt.Sprintf("/%s/", eEmail)
				eEmail = tmp
			case "i":
				tmp := fmt.Sprintf("/%s/", eIP)
				eIP = tmp
			case "p":
				tmp := fmt.Sprintf("/%s/", ePassword)
				ePassword = tmp
			case "q":
				tmp := fmt.Sprintf("/%s/", eHashedPassword)
				eHashedPassword = tmp
			case "n":
				tmp := fmt.Sprintf("/%s/", eName)
				eName = tmp
			default:
				fmt.Printf("\nUnknown parameter set for 'Regex Match': %s\n", alias)
				os.Exit(-1)
			}
		}
	}

	dehash := dehashed.NewDehasher(
		eUsername,
		eEmail,
		eIP,
		ePassword,
		eHashedPassword,
		eName,
		p.Get("max-records").(int),
		p.Get("max-requests").(int),
		p.Get("creds-only").(bool))

	dehash.SetClientCredentials(
		p.Get("key").(string),
		p.Get("authorized-email").(string),
		p.Get("print-balance").(bool))

	filetype := "json"
	if p.Get("output-yaml").(bool) {
		filetype = "yaml"
	}
	if p.Get("output-xml").(bool) {
		filetype = "xml"
	}
	if p.Get("output-txt").(bool) {
		filetype = "txt"
	}

	filename := "dehash"
	tmp := p.Get("output-file-name").(string)
	if len(tmp) > 0 {
		filename = tmp
	}

	dehash.SetOutputFile(
		filetype,
		filename)
	dehash.Start()
	fmt.Println("\n[*] Completing Process")
}

func main() {
	parser := getParser()
	parser.Parse()

	// Create new Dehasher Object
	startDehash(parser)
}
