package query

import (
	"Dehash/internal/sqlite"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type DehashedClient struct {
	key      string
	email    string
	results  []sqlite.Result
	client   *http.Client
	query    string
	params   string
	printBal bool
	total    int
	balance  int
}

var baseUrl = "https://api.dehashed.com/v2/search"

func NewDehashedClient(key, email string, printBal bool) *DehashedClient {
	return &DehashedClient{key: key, email: email, results: make([]sqlite.Result, 0), client: &http.Client{}, printBal: printBal}
}

func (dc *DehashedClient) getKey() string {
	return dc.key
}

func (dc *DehashedClient) getEmail() string {
	return dc.email
}

func (dc *DehashedClient) GetResults() sqlite.DehashedResults {
	return sqlite.DehashedResults{Results: dc.results}
}

func (dc *DehashedClient) buildQuery(params map[string]string) {
	urlParams := url.Values{}
	urlString := baseUrl

	if len(params) > 0 {
		urlString += "?query="

		for k, v := range params {
			if len(v) > 0 {
				urlParams.Add(k, v)
			}
		}
	}

	tmp, _ := url.QueryUnescape(urlParams.Encode())
	tmp2 := strings.Replace(tmp, "=", ":", -1)
	dc.params = tmp2
	urlString += dc.params
	dc.query = urlString
}

func (dc *DehashedClient) setResults(results int) {
	dc.query = fmt.Sprintf("%s?query=%s&size=%d", baseUrl, dc.params, results)
}

func (dc *DehashedClient) setPage(page int) {
	dc.query = fmt.Sprintf("%s&nextPage=%d", dc.query, page)
}

func (dc *DehashedClient) Do() int {
	fmt.Printf("\n\t[*] Performing Request...")
	req, err := http.NewRequest("GET", dc.query, nil)
	if err != nil {
		fmt.Printf("[!] Error constructing request: %v", err)
		os.Exit(-1)
	}

	dc.setAuth(req)
	req.Header.Add("Dehashed-Api-Key", dc.getKey())
	req.Header.Add("Accept", "application/json")
	resp, err := dc.client.Do(req)
	if err != nil {
		fmt.Printf("[!] Error performing request: %s\n%v", dc.query, err)
		os.Exit(-1)
	}

	if resp.StatusCode != 200 {
		dhErr := GetDehashedError(resp.StatusCode)
		fmt.Println()
		log.Fatal(dhErr.Error())
	}

	entries, balance, total := sqlite.NewDehashedResults(resp.Body)
	dc.results = append(dc.results, entries...)
	dc.balance = balance
	dc.total += total
	if dc.printBal {
		fmt.Printf("\n\t\t[*] Balance Remaining: %d", balance)
	}
	return total
}

func (dc *DehashedClient) setAuth(r *http.Request) {
	r.SetBasicAuth(dc.email, dc.key)
}

func (dc *DehashedClient) GetDomains() int {
	return dc.balance
}
