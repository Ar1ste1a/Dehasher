package dehashed

import (
	"fmt"
	"os"
)

type Dehasher struct {
	queryUsername       bool
	queryEmail          bool
	queryIpAddress      bool
	queryPassword       bool
	queryHashedPassword bool
	queryName           bool
	username            string
	email               string
	ipAddress           string
	password            string
	hashedPassword      string
	name                string
	query               string
	page                int
	records             int
	requests            int
	params              map[string]string
}

func NewDehasher(qUsername, qEmail, qIP, qPass, qHPass, qName bool, eUsername, eEmail, eIP, ePass, eHPass, eName string, maxRecords, maxRequests int) *Dehasher {
	dh := &Dehasher{
		queryUsername:       qUsername,
		queryEmail:          qEmail,
		queryIpAddress:      qIP,
		queryPassword:       qPass,
		queryHashedPassword: qHPass,
		queryName:           qName,
		username:            eUsername,
		email:               eEmail,
		ipAddress:           eIP,
		password:            ePass,
		hashedPassword:      eHPass,
		name:                eName,
		records:             maxRecords,
		requests:            maxRequests,
		page:                1,
	}
	dh.constructMap()
	if len(dh.params) == 0 {
		fmt.Println("At least one return type is required")
		os.Exit(-1)
	}
	return dh
}

func (dh *Dehasher) buildQuery() {
	for key, value := range dh.params {
		println(key, value)
	}

	//dh.query = fmt.Sprintf()
}

func (dh *Dehasher) constructMap() {
	urlParams := map[string]string{}

	if dh.queryUsername {
		urlParams["username"] = dh.username
	}
	if dh.queryEmail {
		urlParams["email"] = dh.email
	}
	if dh.queryIpAddress {
		urlParams["ip_address"] = dh.ipAddress
	}
	if dh.queryHashedPassword {
		urlParams["hashed_password"] = dh.hashedPassword
	}
	if dh.queryPassword {
		urlParams["password"] = dh.password
	}
	if dh.queryName {
		urlParams["name"] = dh.name
	}

	dh.params = urlParams
}

type DehashResult struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	IpAddress      string `json:"ip_address"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	HashedPassword string `json:"hashed_password"`
	HashType       string `json:"hash_type"`
	Name           string `json:"name"`
	Vin            string `json:"vin"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	DatabaseName   string `json:"database_name"`
}
