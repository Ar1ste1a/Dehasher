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
	Took         string   `json:"took"`
	TotalResults int      `json:"total"`
}

type Result struct {
	gorm.Model
	DehashedId            string   `json:"id" xml:"id" yaml:"id" gorm:"uniqueIndex"`
	Email                 string   `json:"email,omitempty" xml:"email,omitempty" yaml:"email,omitempty"`
	EmailArray            []string `json:"-" xml:"-" yaml:"-" gorm:"-"`
	IpAddress             string   `json:"ip_address,omitempty" xml:"ip_address,omitempty" yaml:"ip_address,omitempty"`
	Username              string   `json:"username,omitempty" xml:"username,omitempty" yaml:"username,omitempty"`
	Password              string   `json:"password,omitempty" xml:"password,omitempty" yaml:"password,omitempty"`
	HashedPassword        string   `json:"hashed_password,omitempty" xml:"hashed_password,omitempty" yaml:"hashed_password,omitempty"`
	HashType              string   `json:"hash_type,omitempty" xml:"hash_type,omitempty" yaml:"hash_type,omitempty"`
	Name                  string   `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty"`
	NameArray             []string `json:"-" xml:"-" yaml:"-" gorm:"-"`
	Vin                   string   `json:"vin,omitempty" xml:"vin,omitempty" yaml:"vin,omitempty"`
	LicensePlate          string   `json:"license_plate,omitempty" xml:"license_plate,omitempty" yaml:"license_plate,omitempty"`
	Url                   string   `json:"url,omitempty" xml:"url,omitempty" yaml:"url,omitempty"`
	UrlArray              []string `json:"-" xml:"-" yaml:"-" gorm:"-"`
	Social                string   `json:"social,omitempty" xml:"social,omitempty" yaml:"social,omitempty"`
	CryptoCurrencyAddress string   `json:"cryptocurrency_address,omitempty" xml:"cryptocurrency_address,omitempty" yaml:"cryptocurrency_address,omitempty"`
	Address               string   `json:"address,omitempty" xml:"address,omitempty" yaml:"address,omitempty"`
	AddressArray          []string `json:"-" xml:"-" yaml:"-" gorm:"-"`
	Phone                 string   `json:"phone,omitempty" xml:"phone,omitempty" yaml:"phone,omitempty"`
	PhoneArray            []string `json:"-" xml:"-" yaml:"-" gorm:"-"`
	Company               string   `json:"company,omitempty" xml:"company,omitempty" yaml:"company,omitempty"`
	CompanyArray          []string `json:"-" xml:"-" yaml:"-" gorm:"-"`
	DatabaseName          string   `json:"database_name,omitempty" xml:"database_name,omitempty" yaml:"database_name,omitempty"`
}

type DehashedResults struct {
	Results []Result `json:"results"`
}

func (dr *DehashedResults) ExtractCredentials() []Creds {
	var creds []Creds

	results := dr.Results

	for _, r := range results {
		if r.Password != "" {
			// Get first email if available
			email := ""
			if len(r.EmailArray) > 0 {
				email = r.EmailArray[0]
			}

			cred := Creds{Email: email, Password: r.Password}
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

// UnmarshalJSON custom unmarshaler to handle array fields
func (r *Result) UnmarshalJSON(data []byte) error {
	// Create a temporary struct with the same fields
	type ResultAlias Result

	// Create a map for raw decoding
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// Initialize the result with default values
	tmp := &ResultAlias{}

	// Unmarshal the basic fields
	if err := json.Unmarshal(data, tmp); err != nil {
		return err
	}

	// Copy the basic fields
	*r = Result(*tmp)

	// Handle array fields
	if email, ok := rawMap["email"]; ok {
		json.Unmarshal(email, &r.EmailArray)
	}

	if name, ok := rawMap["name"]; ok {
		json.Unmarshal(name, &r.NameArray)
	}

	if url, ok := rawMap["url"]; ok {
		json.Unmarshal(url, &r.UrlArray)
	}

	if address, ok := rawMap["address"]; ok {
		json.Unmarshal(address, &r.AddressArray)
	}

	if phone, ok := rawMap["phone"]; ok {
		json.Unmarshal(phone, &r.PhoneArray)
	}

	if company, ok := rawMap["company"]; ok {
		json.Unmarshal(company, &r.CompanyArray)
	}

	return nil
}

// MarshalJSON custom marshaler to handle array fields
func (r *Result) MarshalJSON() ([]byte, error) {
	type ResultAlias Result

	// Create a map to build our custom JSON
	m := make(map[string]interface{})

	// Convert to a map first
	tmp := ResultAlias(*r)
	tmpData, err := json.Marshal(tmp)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(tmpData, &m)

	// Add array fields to the map
	if len(r.EmailArray) > 0 {
		emailJSON, err := json.Marshal(r.EmailArray)
		if err == nil {
			r.Email = string(emailJSON)
		}
	}

	if len(r.NameArray) > 0 {
		nameJSON, err := json.Marshal(r.NameArray)
		if err == nil {
			r.Name = string(nameJSON)
		}
	}

	if len(r.UrlArray) > 0 {
		urlJSON, err := json.Marshal(r.UrlArray)
		if err == nil {
			r.Url = string(urlJSON)
		}
	}

	if len(r.AddressArray) > 0 {
		addressJSON, err := json.Marshal(r.AddressArray)
		if err == nil {
			r.Address = string(addressJSON)
		}
	}

	if len(r.PhoneArray) > 0 {
		phoneJSON, err := json.Marshal(r.PhoneArray)
		if err == nil {
			r.Phone = string(phoneJSON)
		}
	}

	if len(r.CompanyArray) > 0 {
		companyJSON, err := json.Marshal(r.CompanyArray)
		if err == nil {
			r.Company = string(companyJSON)
		}
	}

	return json.Marshal(m)
}

// AfterFind deserializes JSON strings to arrays after fetching from database
func (r *Result) AfterFind(tx *gorm.DB) error {
	// Convert JSON strings back to arrays
	if r.Email != "" {
		json.Unmarshal([]byte(r.Email), &r.EmailArray)
	}

	if r.Name != "" {
		json.Unmarshal([]byte(r.Name), &r.NameArray)
	}

	if r.Url != "" {
		json.Unmarshal([]byte(r.Url), &r.UrlArray)
	}

	if r.Address != "" {
		json.Unmarshal([]byte(r.Address), &r.AddressArray)
	}

	if r.Phone != "" {
		json.Unmarshal([]byte(r.Phone), &r.PhoneArray)
	}

	if r.Company != "" {
		json.Unmarshal([]byte(r.Company), &r.CompanyArray)
	}

	return nil
}
