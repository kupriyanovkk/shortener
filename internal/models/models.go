package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type URL struct {
	UUID        int    `json:"uuid"`
	Short       string `json:"short_url"`
	Original    string `json:"original_url"`
	UserID      string `json:"user_id"`
	DeletedFlag bool   `json:"is_deleted"`
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserURL struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}
