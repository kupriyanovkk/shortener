package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type URL struct {
	UUID     int    `json:"uuid"`
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}
