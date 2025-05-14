package sqlite

type DehashedSearchRequest struct {
	Page     int    `json:"page"`
	Query    string `json:"query"`
	Size     int    `json:"size"`
	Wildcard bool   `json:"wildcard"`
	Regex    bool   `json:"regex"`
	DeDupe   bool   `json:"de_dupe"`
}

func NewDehashedSearchRequest(size int, wildcard, regex bool) *DehashedSearchRequest {
	return &DehashedSearchRequest{
		Page:     0,
		Size:     size,
		Wildcard: false,
		Regex:    false,
		DeDupe:   true,
	}
}
