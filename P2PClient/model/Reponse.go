package model

// Response - Reponse struct
type Response struct {
	Length string
	Code   string
	Count  string
	Ips    []string
}

// SearchResponse - Response for search
type SearchResponse struct {
	Length string
	Code   string
	Count  int
	IP     string
	Port   string
	Hops   string
	Files  []string
}
