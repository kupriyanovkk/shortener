package models

// Request struct
type Request struct {
	URL string `json:"url"`
}

// Response struct
type Response struct {
	Result string `json:"result"`
}

// URL is a structure contains all URL data
type URL struct {
	UUID        int    `json:"uuid"`
	Short       string `json:"short_url"`
	Original    string `json:"original_url"`
	UserID      string `json:"user_id"`
	DeletedFlag bool   `json:"is_deleted"`
}

// BatchRequest is a structure for URL batching
type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse is a structure for URL batching
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserURL is a structure for user
type UserURL struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

// InternalStats is a structure for internal statistics
type InternalStats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}
