package query

import (
	"Dehash/internal/sqlite"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type DehashedParameter string

const (
	Username       DehashedParameter = "username"
	Email          DehashedParameter = "email"
	Password       DehashedParameter = "password"
	HashedPassword DehashedParameter = "hashed_password"
	Name           DehashedParameter = "name"
	IpAddress      DehashedParameter = "ip_address"
	Domain         DehashedParameter = "domain"
	Vin            DehashedParameter = "vin"
	LicensePlate   DehashedParameter = "license_plate"
	Address        DehashedParameter = "address"
	Phone          DehashedParameter = "phone"
	Social         DehashedParameter = "social"
	CryptoAddress  DehashedParameter = "cryptocurrency_address"
)

func (dp DehashedParameter) GetArgumentString(arg string) string {
	return fmt.Sprintf("%s:%s", string(dp), arg)
}

type DehashedSearchRequest struct {
	ForcePlaintext bool   `json:"-"`
	Page           int    `json:"page"`
	Query          string `json:"query"`
	Size           int    `json:"size"`
	Wildcard       bool   `json:"wildcard"`
	Regex          bool   `json:"regex"`
	DeDupe         bool   `json:"de_dupe"`
}

func NewDehashedSearchRequest(page, size int, wildcard, regex, forcePlaintext bool) *DehashedSearchRequest {
	return &DehashedSearchRequest{Page: page, Query: "", Size: size, Wildcard: wildcard, Regex: regex, DeDupe: true, ForcePlaintext: forcePlaintext}
}

func (dsr *DehashedSearchRequest) buildQuery(query string, param DehashedParameter) {
	if len(dsr.Query) > 0 {
		dsr.Query = fmt.Sprintf("%s&%s", strings.TrimSpace(dsr.Query), strings.TrimSpace(query))
	} else {
		dsr.Query = query
	}
}

func (dsr *DehashedSearchRequest) AddUsernameQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(Username.GetArgumentString(query), Username)
}

func (dsr *DehashedSearchRequest) AddEmailQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(Email.GetArgumentString(query), Email)
}

func (dsr *DehashedSearchRequest) AddIpAddressQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(IpAddress.GetArgumentString(query), IpAddress)
}

func (dsr *DehashedSearchRequest) AddDomainQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(Domain.GetArgumentString(query), Domain)
}

func (dsr *DehashedSearchRequest) AddPasswordQuery(query string) {
	if dsr.ForcePlaintext {
		dsr.buildQuery(Password.GetArgumentString(query), Password)
		return
	}
	hash := sha256.Sum256([]byte(query))
	query = hex.EncodeToString(hash[:])
	dsr.AddHashedPasswordQuery(query)
}

func (dsr *DehashedSearchRequest) AddVinQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(Vin.GetArgumentString(query), Vin)
}

func (dsr *DehashedSearchRequest) AddLicensePlateQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(LicensePlate.GetArgumentString(query), LicensePlate)
}

func (dsr *DehashedSearchRequest) AddAddressQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(Address.GetArgumentString(query), Address)
}

func (dsr *DehashedSearchRequest) AddPhoneQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(Phone.GetArgumentString(query), Phone)
}

func (dsr *DehashedSearchRequest) AddSocialQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(Social.GetArgumentString(query), Social)
}

func (dsr *DehashedSearchRequest) AddCryptoAddressQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(CryptoAddress.GetArgumentString(query), CryptoAddress)
}

func (dsr *DehashedSearchRequest) AddHashedPasswordQuery(query string) {
	dsr.buildQuery(HashedPassword.GetArgumentString(query), HashedPassword)
}

func (dsr *DehashedSearchRequest) AddNameQuery(query string) {
	query = strings.TrimSpace(query)
	dsr.buildQuery(Name.GetArgumentString(query), Name)
}

type DehashedClientV2 struct {
	apiKey  string
	results []sqlite.Result
}

func NewDehashedClientV2(apiKey string) *DehashedClientV2 {
	return &DehashedClientV2{apiKey: apiKey}
}

func (dcv2 *DehashedClientV2) Search(searchRequest DehashedSearchRequest) (int, error) {
	reqBody, _ := json.Marshal(searchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/search", bytes.NewReader(reqBody))
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", dcv2.apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		zap.L().Error("v2_search",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return -1, err
	}
	if res == nil {
		zap.L().Error("v2_search",
			zap.String("message", "response was nil"),
		)
		return -1, errors.New("response was nil")
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("v2_search",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return -1, err
	}

	var responseResults sqlite.DehashedResponse
	err = json.Unmarshal(b, &responseResults)
	if err != nil {
		zap.L().Error("v2_search",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		return -1, err
	}

	dcv2.results = append(dcv2.results, responseResults.Entries...)
	return responseResults.TotalResults, nil
}

func (dcv2 *DehashedClientV2) GetResults() sqlite.DehashedResults {
	return sqlite.DehashedResults{Results: dcv2.results}
}

func (dcv2 *DehashedClientV2) GetTotalResults() int {
	return len(dcv2.results)
}
