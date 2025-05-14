package sqlite

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"os"
)

type DehashedResponse struct {
	Balance      int      `json:"balance"`
	Entries      []Result `json:"entries"`
	Success      bool     `json:"success"`
	Took         string   `json:"took"`
	TotalResults int      `json:"total"`
}

type Result struct {
	gorm.Model
	Id                    string `json:"id" xml:"id" yaml:"id"`
	Email                 string `json:"email" xml:"email" yaml:"email"`
	IpAddress             string `json:"ip_address" xml:"ip_address" yaml:"ip_address"`
	Username              string `json:"username" xml:"username" yaml:"username"`
	Password              string `json:"password" xml:"password" yaml:"password"`
	HashedPassword        string `json:"hashed_password" xml:"hashed_password" yaml:"hashed_password"`
	HashType              string `json:"hash_type" xml:"hash_type" yaml:"hash_type"`
	Name                  string `json:"name" xml:"name" yaml:"name"`
	Vin                   string `json:"vin" xml:"vin" yaml:"vin"`
	LicensePlate          string `json:"license_plate" xml:"license_plate" yaml:"license_plate"`
	Url                   string `json:"url" xml:"url" yaml:"url"`
	Social                string `json:"social" xml:"social" yaml:"social"`
	CryptoCurrencyAddress string `json:"cryptocurrency_address" xml:"cryptocurrency_address" yaml:"cryptocurrency_address"`
	Address               string `json:"address" xml:"address" yaml:"address"`
	Phone                 string `json:"phone" xml:"phone" yaml:"phone"`
	DatabaseName          string `json:"database_name" xml:"database_name" yaml:"database_name"`
	//RawRecord             RawRecord `json:"raw_record" xml:"raw_record" yaml:"raw_record"`
}

//type RawRecord struct {
//	LeOnly       bool   `json:"le_only" xml:"le_only" yaml:"le_only"`
//	Unstructured string `json:"unstructured" xml:"unstructured" yaml:"unstructured"`
//}

type DehashedResults struct {
	Results []Result `json:"results"`
}

func (dr *DehashedResults) ExtractCredentials() []Creds {
	var cred Creds
	var creds []Creds

	results := dr.Results

	for _, r := range results {
		if len(r.Password) > 0 {
			cred = Creds{Username: r.Username, Email: r.Email, Password: r.Password}
			creds = append(creds, cred)
		}
	}

	return creds
}

func NewDehashedResults(body io.Reader) ([]Result, int, int) {
	var response DehashedResponse

	err := json.NewDecoder(body).Decode(&response)
	if err != nil {
		fmt.Printf("Error Parsing Response Body: %v", err)
		os.Exit(-1)
	}

	return response.Entries, response.Balance, response.TotalResults
}
