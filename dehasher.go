package main

import (
	"Dehash/internal/argparse"
	"Dehash/internal/dehashed"
	"strings"
)

func getParser() *argparse.Argparse {
	parser := argparse.NewParser("Dehasher", "Dehashed Tool", "v1.0")
	parser.AddArg("-k", "--key", "API Key", "string", "", true)
	parser.AddArg("-a", "--authorized-email", "Email to pair with key for authentication", "string", "", true)
	parser.AddArg("-m", "--max-records", "Maximum amount of records to return", "integer", "30000", false)
	parser.AddArg("-r", "--max-requests", "Maximum number of requests to make", "integer", "-1", false)
	parser.AddArg("-x", "--exact-match", "Use Exact Matching on fields", "string", "", false)
	parser.AddArg("-t", "--list-tokens", "List the number of tokens remaining", "bool", "false", false)
	parser.AddArg("-o", "--output-file", "File to output results to", "string", "dehashed", false)
	parser.AddArg("-T", "--output-txt", "Output to text file", "bool", "false", false)
	parser.AddArg("-J", "--output-json", "Output to JSON file", "bool", "true", false)
	parser.AddArg("-Y", "--output-yaml", "Output to YAML file", "bool", "false", false)
	parser.AddArg("-A", "--query-all-fields", "Return All Fields", "bool", "false", false)
	parser.AddArg("-u", "--query-username", "Return Usernames", "bool", "false", false)
	parser.AddArg("-U", "--username-query", "Username Query", "string", "", false)
	parser.AddArg("-e", "--query-email", "Return Emails", "bool", "false", false)
	parser.AddArg("-E", "--email-query", "Email Query", "string", "", false)
	parser.AddArg("-i", "--query-ip-address", "Return IP Addresses", "bool", "false", false)
	parser.AddArg("-I", "--ip-address-query", "IP Address Query", "string", "", false)
	parser.AddArg("-p", "--query-password", "Return Passwords", "bool", "false", false)
	parser.AddArg("-P", "--password-query", "Password Query", "string", "", false)
	parser.AddArg("-q", "--query-hashed-password", "Return Hashed Passwords", "bool", "false", false)
	parser.AddArg("-Q", "--hashed-password-query", "Hashed Password Query", "string", "", false)
	parser.AddArg("-n", "--query-name", "Return Names", "bool", "false", false)
	parser.AddArg("-N", "--name-query", "Name Query", "string", "", false)

	return parser
}

func startDehash(p *argparse.Argparse) {
	var (
		eUsername, eEmail, eIP, ePassword, eHashedPassword, eName string
		qUsername                                                 = false
		qEmail                                                    = false
		qIP                                                       = false
		qPassword                                                 = false
		qHashedPassword                                           = false
		qName                                                     = false
	)

	if p.Get("query-all-fields").(bool) {
		qUsername = true
		qEmail = true
		qIP = true
		qPassword = true
		qHashedPassword = true
		qName = true
	} else {
		qUsername = p.Get("query-username").(bool)
		qEmail = p.Get("query-email").(bool)
		qIP = p.Get("query-ip-address").(bool)
		qPassword = p.Get("query-password").(bool)
		qHashedPassword = p.Get("query-hashed-password").(bool)
		qName = p.Get("query-name").(bool)
	}

	if qUsername {
		eUsername = p.Get("username-query").(string)
	}

	if qEmail {
		eEmail = p.Get("email-query").(string)
	}

	if qIP {
		eIP = p.Get("ip-address-query").(string)
	}

	if qPassword {
		ePassword = p.Get("password-query").(string)
	}

	if qHashedPassword {
		eHashedPassword = p.Get("hashed-password-query").(string)
	}

	if qName {
		eName = p.Get("name-query").(string)
	}

	exactMatch := strings.ToLower(p.Get("exact-match").(string))
	if len(exactMatch) > 0 {
		// parse the exact match flags, encompass in double quotes
	}

	dehash := dehashed.NewDehasher(
		qUsername,
		qEmail,
		qIP,
		qPassword,
		qHashedPassword,
		qName,
		eUsername,
		eEmail,
		eIP,
		ePassword,
		eHashedPassword,
		eName,
		p.Get("max-records").(int),
		p.Get("max-requests").(int))

	println(dehash)
}

func main() {
	parser := getParser()
	parser.Parse()

	// Create new Dehasher Object
	startDehash(parser)
}
