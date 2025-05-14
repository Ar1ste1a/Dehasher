package sqlite

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
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
	DehashedId            string   `json:"id" xml:"id" yaml:"id" gorm:"uniqueIndex"`
	Email                 []string `json:"email,omitempty" xml:"email,omitempty" yaml:"email,omitempty" gorm:"serializer:json"`
	IpAddress             []string `json:"ip_address,omitempty" xml:"ip_address,omitempty" yaml:"ip_address,omitempty" gorm:"serializer:json"`
	Username              []string `json:"username,omitempty" xml:"username,omitempty" yaml:"username,omitempty" gorm:"serializer:json"`
	Password              []string `json:"password,omitempty" xml:"password,omitempty" yaml:"password,omitempty" gorm:"serializer:json"`
	HashedPassword        []string `json:"hashed_password,omitempty" xml:"hashed_password,omitempty" yaml:"hashed_password,omitempty" gorm:"serializer:json"`
	HashType              string   `json:"hash_type,omitempty" xml:"hash_type,omitempty" yaml:"hash_type,omitempty"`
	Name                  []string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty" gorm:"serializer:json"`
	Vin                   []string `json:"vin,omitempty" xml:"vin,omitempty" yaml:"vin,omitempty" gorm:"serializer:json"`
	LicensePlate          []string `json:"license_plate,omitempty" xml:"license_plate,omitempty" yaml:"license_plate,omitempty" gorm:"serializer:json"`
	Url                   []string `json:"url,omitempty" xml:"url,omitempty" yaml:"url,omitempty" gorm:"serializer:json"`
	Social                []string `json:"social,omitempty" xml:"social,omitempty" yaml:"social,omitempty" gorm:"serializer:json"`
	CryptoCurrencyAddress []string `json:"cryptocurrency_address,omitempty" xml:"cryptocurrency_address,omitempty" yaml:"cryptocurrency_address,omitempty" gorm:"serializer:json"`
	Address               []string `json:"address,omitempty" xml:"address,omitempty" yaml:"address,omitempty" gorm:"serializer:json"`
	Phone                 []string `json:"phone,omitempty" xml:"phone,omitempty" yaml:"phone,omitempty" gorm:"serializer:json"`
	Company               []string `json:"company,omitempty" xml:"company,omitempty" yaml:"company,omitempty" gorm:"serializer:json"`
	DatabaseName          string   `json:"database_name,omitempty" xml:"database_name,omitempty" yaml:"database_name,omitempty"`
}

type DehashedResults struct {
	Results []Result `json:"results"`
}

func (dr *DehashedResults) ExtractCredentials() []Creds {
	var creds []Creds

	results := dr.Results

	for _, r := range results {
		if len(r.Password) > 0 {
			// Get first email if available
			email := ""
			if len(r.Email) > 0 {
				email = r.Email[0]
			}

			// Get first password
			password := r.Password[0]

			cred := Creds{Email: email, Password: password}
			creds = append(creds, cred)
		}
	}

	go func() {
		err := StoreCreds(creds)
		if err != nil {
			zap.L().Error("store_creds",
				zap.String("message", "failed to store creds"),
				zap.Error(err),
			)
			fmt.Printf("Error Storing Results: %v", err)
		}
	}()

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
