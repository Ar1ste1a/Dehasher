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
	dc.query = fmt.Sprintf("%s&page=%d", dc.query, page)
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

/*
Default results per call 100
Max per call 10,000
Default MaxRecords accessible via pagination, 30,000,
rate limit: 5 MaxRequests per second (per ip + credential combo. More creds/ more queries)
	rate limit response:
		HTTP Response Code: 400 { "Error 400": "Too many MaxRequests were performed in a small amount of time. Please wait a bit before querying the API."}
	Unauthorized response:
		HTTP Response Code: 401 { "message": "Invalid API credentials.", "success": false }
	Method not allowed:
		HTTP Response Code: 404
	Invalid Query/Missing Query
		HTTP Response Code: 302
*/

/* Authentication


 */

/*
Queries: Must be a GET
	to search for exact record, wrap with double quotes
		https://api.dehashed.com/search?query=username:"test"

		where the username field contains dave123
                     username:dave123
		where the email field contains dave or david. If you omit the OR operator the default operator will be used:
                    email:(dave OR david)
                    email:(dave david)
		where the name field contains the exact phrase "john smith":
                    name:"John Smith"

		REGEX
		Regular expression patterns can be embedded in the query string by wrapping them in forward-slashes "/":
                    name:/joh?n(ath[oa]n)/

		Reserved Chars
		If you need to use any of the characters which function as operators in your query itself (and not as operators), then you should escape them with a leading backslash. For instance, to search for (1+1)=2, you would need to write your query as \(1\+1\)\=2.
		The reserved characters are: + - = && || > < ! ( ) { } [ ] ^ " ~ * ? : \
		Failing to escape these special characters correctly could lead to a syntax error which prevents your query from running.

		Sizing and Pagination
		DeHashed allows users to query up to 30,000 results and up to 10,000 results per call.
		Sizing was introduced to help users achieve two things: Either speed up response times, or save on credits..
		If you don't care for speed, and want to reduce your credit cost (instead of paginating 10x and paying 10 credits), you could increase the &MaxRecords= parameter, the limit is 10,000. This will significantly slow down your search, however return more results in one call. If you care about speed, you could leave the parameter to default (100) or reduce it further to increase speed.
		Pagination hasn't changed. Simply add the &page= parameter to your search, and indicate the next set of results you wish to access
		Pagination and Sizing can, and should be used together. The current limit on pagination is 30,000 results. If you set MaxRecords to 1 (&MaxRecords=1) you could paginate to the 30,000th page (&page=30000). If you set the MaxRecords to 10,000(&MaxRecords=10000) then you can only paginate to the 3rd page (&page=3)

		Getting Next set of results
		Response default to 100 results per query, to get the next 100 results, simply add &page=2 to the end of the url. (Note: results are limited to 30,000).

		Calling API, and decreasing/increasing the result amount:
		Response defaults to 100 results per query, to increase or decrease the amount of results/MaxRecords per response, simply add &MaxRecords={amount} to the end of the url. For Example: limiting to just 1 Result per API Response:&MaxRecords=1 (Fast, however it can cost quite a bit of credits querying a lot). If you don't care for speed, and want to save on credits you can simply append &MaxRecords=10000 to the end of your request URL. (Note: results are limited to 10,000. You cannot paginate past the 2nd page of results assuming you keep MaxRecords parameter default (100)).
*/

/*
Possible Queries
email, ip_address, username, password, hashed_password, name, and any other data points.
*/

/*
Possible results

id, email, ip_address, username, password, hashed_password, hash_type, name, vin, address, phone, database_name
*/
