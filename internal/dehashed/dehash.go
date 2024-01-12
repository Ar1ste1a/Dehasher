package dehashed

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

type Dehasher struct {
	username       string
	email          string
	ipAddress      string
	password       string
	hashedPassword string
	name           string
	query          string
	page           int
	records        int
	requests       int
	params         map[string]string
	client         *DehashClient
	fileType       string
	fileName       string
	credsOnly      bool
}

func NewDehasher(eUsername, eEmail, eIP, ePass, eHPass, eName string, maxRecords, maxRequests int, credsOnly bool) *Dehasher {
	dh := &Dehasher{
		username:       eUsername,
		email:          eEmail,
		ipAddress:      eIP,
		password:       ePass,
		hashedPassword: eHPass,
		name:           eName,
		records:        maxRecords,
		requests:       maxRequests,
		page:           1,
		credsOnly:      credsOnly,
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

func (dh *Dehasher) escapeReservedCharacters() {
	reserved := strings.Split("+ - = && || > < ! ( ) { } [ ] ^ \" ~ * ? : \\", " ")

	dh.username = escapeString(dh.username, reserved)
	dh.email = escapeString(dh.email, reserved)
	dh.ipAddress = escapeString(dh.ipAddress, reserved)
	dh.password = escapeString(dh.password, reserved)
	dh.hashedPassword = escapeString(dh.hashedPassword, reserved)
	dh.name = escapeString(dh.name, reserved)
}

func (dh *Dehasher) SetClientCredentials(key, email string, printBal bool) {
	dh.client = NewDehashClient(key, email, printBal)
}

func (dh *Dehasher) SetOutputFile(filetype, filename string) {
	dh.fileType = filetype
	dh.fileName = filename
}

func (dh *Dehasher) setQueries() {
	var numQueries int

	switch {
	case dh.requests == 0:
		fmt.Println("Max Requests cannot be zero")
		os.Exit(-1)
	case dh.records <= 10000 || dh.requests == 1:
		numQueries = 1
		if dh.records > 10000 {
			dh.records = 10000
		}
	case dh.requests < 0 && dh.records > 20000:
		numQueries = 3
		dh.records = 10000
	case dh.requests < 0 && dh.records > 10000:
		numQueries = 2
		dh.records = 10000
	case dh.records < 0 && dh.records < 10000:
		numQueries = 1
	case dh.requests == 2 && dh.records > 20000:
		numQueries = 2
		dh.records = 10000
	case dh.requests == 2 && dh.records <= 10000:
		numQueries = 1
	default:
		numQueries = 3
		dh.records = 10000
	}

	dh.requests = numQueries
	fmt.Printf("Making %d Requests for %d Records (%d Total)", dh.requests, dh.records, dh.requests*dh.records)
}

func (dh *Dehasher) Start() {
	dh.client.buildQuery(dh.params)
	page := 1
	offset := 0
	foundNum := 0
	for i := 0; i < dh.requests; i++ {
		dh.client.setResults(dh.records)
		dh.client.setPage(page)
		found := dh.client.Do()

		if found-offset < dh.records {
			fmt.Printf("\n\t\t[+] Retrieved %d Records", found-offset)
			fmt.Printf("\n[-] Not Enough Entries, ending queries")
			break
		} else {
			if found-offset > dh.records {
				foundNum = dh.records
			} else {
				foundNum = found - offset
			}
			fmt.Printf("\n\t\t[*] Retrieved %d Records", foundNum)
			offset += dh.records
			page += 1
		}
	}

	dh.parseResults()
}

func (dh *Dehasher) buildQuery() {
	for key, value := range dh.params {
		println(key, value)
	}
}

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

func escapeString(input string, reserved []string) string {
	for _, char := range reserved {
		input = strings.Replace(input, char, "\\"+char, -1)
	}
	return input
}

func (dh *Dehasher) parseResults() {
	var data []byte
	results := dh.client.GetResults()

	if len(results) > 0 {
		fmt.Printf("\n\t[*] Writing entries to file: %s.%s", dh.fileName, dh.fileType)
		if !dh.credsOnly {
			err := dh.writeToFile(results)
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
			creds := dh.extractCreds(results)
			err := dh.writeCredsToFile(creds)
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

func (dh *Dehasher) writeCredsToFile(creds []Creds) error {
	var data []byte
	var err error

	switch dh.fileType {
	case "json":
		data, err = json.MarshalIndent(creds, "", "  ")
	case "xml":
		data, err = xml.MarshalIndent(creds, "", "  ")
	case "yaml":
		data, err = yaml.Marshal(creds)
	case "txt":
		var outStrings []string
		for _, c := range creds {
			outStrings = append(outStrings, c.ToString()+"\n")
		}
		data = []byte(strings.Join(outStrings, ""))
	default:
		return errors.New("unsupported file type")
	}

	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s.%s", dh.fileName, dh.fileType)
	return ioutil.WriteFile(filePath, data, 0644)
}

func (dh *Dehasher) writeToFile(result []DehashResult) error {
	var data []byte
	var err error

	switch dh.fileType {
	case "json":
		data, err = json.MarshalIndent(result, "", "  ")
	case "xml":
		data, err = xml.MarshalIndent(result, "", "  ")
	case "yaml":
		data, err = yaml.Marshal(result)
	case "txt":
		var outStrings []string
		for _, r := range result {
			out := fmt.Sprintf(
				"Id: %s\nEmail: %s\nIpAddress: %s\nUsername: %s\nPassword: %s\nHashedPassword: %s\nHashType: %s\nName: %s\nVin: %s\nAddress: %s\nPhone: %s\nDatabaseName: %s\n\n",
				r.Id, r.Email, r.IpAddress, r.Username, r.Password, r.HashedPassword, r.HashType, r.Name, r.Vin, r.Address, r.Phone, r.DatabaseName)
			outStrings = append(outStrings, out)
		}
		data = []byte(strings.Join(outStrings, ""))
	default:
		return errors.New("unsupported file type")
	}

	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s.%s", dh.fileName, dh.fileType)
	return ioutil.WriteFile(filePath, data, 0644)
}

type Creds struct {
	Username string `json:"username" yaml:"username" xml:"username"`
	Password string `json:"password" yaml:"password" xml:"password"`
}

func (c Creds) ToString() string {
	return fmt.Sprintf("%s%s%s", c.Username, "%", c.Password)
}

func (dh *Dehasher) extractCreds(results []DehashResult) []Creds {
	var cred Creds
	var creds []Creds

	for _, r := range results {
		if len(r.Password) > 0 {
			var username string
			if len(r.Username) > 0 {
				username = r.Username
			}
			if len(r.Email) > 0 {
				if len(username) > 0 {
					username = fmt.Sprintf("%s/%s", r.Username, r.Email)
				} else {
					username = r.Email
				}
			}
			cred = Creds{Username: username, Password: r.Password}
			creds = append(creds, cred)
		}
	}

	return creds
}
