package whois

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type DehashedWHOISSearchRequest struct {
	Include     []string `json:"include,omitempty"`
	Exclude     []string `json:"exclude,omitempty"`
	IPAddress   string   `json:"ip_address,omitempty"`
	ReverseType string   `json:"reverse_type,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	MXAddress   string   `json:"mx_address,omitempty"`
	NSAddress   string   `json:"ns_address,omitempty"`
	SearchType  string   `json:"search_type,omitempty"`
}

func WhoisSearch(domain, apiKey string) (string, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		Domain:     domain,
		SearchType: "whois",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("response was nil")
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func WhoisHistory(domain, apiKey string) (string, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		Domain:     domain,
		SearchType: "whois-history",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		zap.L().Info("whois_history",
			zap.String("message", "response was not nil"),
		)
		defer res.Body.Close()
	}
	if err != nil {
		zap.L().Error("whois_history",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return "", err
	}
	if res == nil {
		zap.L().Error("whois_history",
			zap.String("message", "response was nil"),
		)
		return "", errors.New("response was nil")
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("whois_history",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return "", err
	}
	return string(b), nil
}

func ReverseWHOIS(include []string, exclude []string, reverseType, apiKey string) (string, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		Include:     include,
		Exclude:     exclude,
		ReverseType: reverseType,
		SearchType:  "reverse-whois",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("response was nil")
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func WhoisIP(ipAddress, apiKey string) ([]byte, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		IPAddress:  ipAddress,
		SearchType: "reverse-ip",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("response was nil")
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func WhoisMX(mxAddress, apiKey string) (string, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		MXAddress:  mxAddress,
		SearchType: "reverse-mx",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("response was nil")
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func WhoisNS(nsAddress, apiKey string) (string, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		NSAddress:  nsAddress,
		SearchType: "reverse-ns",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("response was nil")
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func WhoisSubdomainScan(domain, apiKey string) (string, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		Domain:     domain,
		SearchType: "subdomain-scan",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("response was nil")
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func GetWHOISCredits(apiKey string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.dehashed.com/v2/whois/credits", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("response was nil")
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
