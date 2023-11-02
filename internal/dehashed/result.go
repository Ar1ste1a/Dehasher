package dehashed

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type DehashResponse struct {
	Balance      int            `json:"balance"`
	Entries      []DehashResult `json:"entries"`
	Success      bool           `json:"success"`
	Took         string         `json:"took"`
	TotalResults int            `json:"total"`
}

//type DehashResult struct {
//	Id             string `json:"id"`
//	Email          string `json:"email"`
//	IpAddress      string `json:"ip_address"`
//	Username       string `json:"username"`
//	Password       string `json:"password"`
//	HashedPassword string `json:"hashed_password"`
//	HashType       string `json:"hash_type"`
//	Name           string `json:"name"`
//	Vin            string `json:"vin"`
//	Address        string `json:"address"`
//	Phone          string `json:"phone"`
//	DatabaseName   string `json:"database_name"`
//}

type DehashResult struct {
	Id             string `json:"id" xml:"id" yaml:"id"`
	Email          string `json:"email" xml:"email" yaml:"email"`
	IpAddress      string `json:"ip_address" xml:"ip_address" yaml:"ip_address"`
	Username       string `json:"username" xml:"username" yaml:"username"`
	Password       string `json:"password" xml:"password" yaml:"password"`
	HashedPassword string `json:"hashed_password" xml:"hashed_password" yaml:"hashed_password"`
	HashType       string `json:"hash_type" xml:"hash_type" yaml:"hash_type"`
	Name           string `json:"name" xml:"name" yaml:"name"`
	Vin            string `json:"vin" xml:"vin" yaml:"vin"`
	Address        string `json:"address" xml:"address" yaml:"address"`
	Phone          string `json:"phone" xml:"phone" yaml:"phone"`
	DatabaseName   string `json:"database_name" xml:"database_name" yaml:"database_name"`
}

func NewDehashResults(body io.Reader) ([]DehashResult, int, int) {
	var response DehashResponse

	err := json.NewDecoder(body).Decode(&response)
	if err != nil {
		fmt.Printf("Error Parsing Response Body: %v", err)
		os.Exit(-1)
	}

	return response.Entries, response.Balance, response.TotalResults
}
