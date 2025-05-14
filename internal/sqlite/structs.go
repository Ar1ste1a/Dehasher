package sqlite

import (
	"Dehash/internal/files"
	"fmt"
	"gorm.io/gorm"
)

type DBOptions struct {
	Username       string
	Email          string
	IPAddress      string
	Password       string
	HashedPassword string
	Name           string
	Limit          int
	ExactMatch     bool
}

func NewDBOptions() *DBOptions {
	return &DBOptions{
		Limit:      100, // Default limit
		ExactMatch: false,
	}
}

type QueryOptions struct {
	gorm.Model
	MaxRecords         int            `json:"max_records"`
	MaxRequests        int            `json:"max_requests"`
	OutputFormat       files.FileType `json:"output_format"`
	OutputFile         string         `json:"output_file"`
	RegexMatch         bool           `json:"regex_match"`
	WildcardMatch      bool           `json:"wildcard_match"`
	UsernameQuery      string         `json:"username_query"`
	EmailQuery         string         `json:"email_query"`
	IpQuery            string         `json:"ip_query"`
	PassQuery          string         `json:"pass_query"`
	HashQuery          string         `json:"hash_query"`
	NameQuery          string         `json:"name_query"`
	DomainQuery        string         `json:"domain_query"`
	VinQuery           string         `json:"vin_query"`
	LicensePlateQuery  string         `json:"license_plate_query"`
	AddressQuery       string         `json:"address_query"`
	PhoneQuery         string         `json:"phone_query"`
	SocialQuery        string         `json:"social_query"`
	CryptoAddressQuery string         `json:"crypto_address_query"`
	PrintBalance       bool           `json:"print_balance"`
	CredsOnly          bool           `json:"creds_only"`
}

func NewQueryOptions(maxRecords, maxRequests int, outputFormat, outputFile, usernameQuery, emailQuery, ipQuery, passQuery, hashQuery, nameQuery, domainQuery, vinQuery, licensePlateQuery, addressQuery, phoneQuery, socialQuery, cryptoAddressQuery string, regexMatch, wildcardMatch, printBalance, credsOnly bool) *QueryOptions {
	return &QueryOptions{
		MaxRecords:         maxRecords,
		MaxRequests:        maxRequests,
		OutputFormat:       files.GetFileType(outputFormat),
		OutputFile:         outputFile,
		PrintBalance:       printBalance,
		CredsOnly:          credsOnly,
		RegexMatch:         regexMatch,
		WildcardMatch:      wildcardMatch,
		UsernameQuery:      usernameQuery,
		EmailQuery:         emailQuery,
		IpQuery:            ipQuery,
		PassQuery:          passQuery,
		HashQuery:          hashQuery,
		NameQuery:          nameQuery,
		DomainQuery:        domainQuery,
		VinQuery:           vinQuery,
		LicensePlateQuery:  licensePlateQuery,
		AddressQuery:       addressQuery,
		PhoneQuery:         phoneQuery,
		SocialQuery:        socialQuery,
		CryptoAddressQuery: cryptoAddressQuery,
	}
}

type Creds struct {
	gorm.Model
	Email    string `json:"email" yaml:"email" xml:"email"`
	Username string `json:"username" yaml:"username" xml:"username"`
	Password string `json:"password" yaml:"password" xml:"password"`
}

func (c Creds) ToString() string {
	return fmt.Sprintf("%s%s%s", c.Username, "%", c.Password)
}
