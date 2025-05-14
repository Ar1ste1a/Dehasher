package query

type DehashError struct {
	Message string
	Code    int
}

type DehashResponseError struct {
	HttpResponse int `json:"HTTP Response Code"`
}

func (de *DehashError) Error() string {
	return de.Message
}

func GetDehashedError(c int) DehashError {
	switch c {
	case 400:
		return DehashError{Code: 400, Message: "There is an issue with authentication. Please check your API key and email. If you haven't, refresh your API Key "}
	case 401:
		return DehashError{Code: 401, Message: "You need a search subscription and API credits to use the API, please purchase a search subscription and add credits to your account."}
	case 403:
		return DehashError{Code: 403, Message: "Insufficient Credits"}
	case 404:
		return DehashError{Code: 404, Message: "Method not permitted"}
	case 429:
		return DehashError{Code: 420, Message: "Rate Limited"}
	case 302:
		return DehashError{Code: 302, Message: "Invalid/Missing Query"}
	default:
		return DehashError{Code: -1, Message: "An unknown error has occurred"}
	}
}
