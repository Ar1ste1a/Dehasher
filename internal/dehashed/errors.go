package dehashed

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
		return DehashError{Code: 400, Message: "Too many requests were performed in a small amount of time. Please wait before querying the API."}
	case 401:
		return DehashError{Code: 401, Message: "Invalid API Credentials"}
	case 404:
		return DehashError{Code: 404, Message: "Method not permitted"}
	case 302:
		return DehashError{Code: 302, Message: "Invalid/Missing Query"}
	default:
		return DehashError{Code: -1, Message: "An unknown error has occurred"}
	}
}
