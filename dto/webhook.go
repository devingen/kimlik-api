package dto

type WebhookPreRequest struct {
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	QueryParams map[string][]string    `json:"queryParams"`
	Header      map[string]string      `json:"header"`
	Body        map[string]interface{} `json:"body"`
}

type WebhookPreResponse struct {
}

type WebhookConsumeSAMLAuthResponseRequest struct {
	User        WebhookConsumeSAMLAuthResponseRequestUser `json:"user"`
	QueryParams map[string][]string                       `json:"queryParams"`
}

type WebhookConsumeSAMLAuthResponseRequestUser struct {
	Email     string                 `json:"email"`
	FirstName string                 `json:"firstName"`
	LastName  string                 `json:"lastName"`
	Meta      map[string]interface{} `json:"meta"`
}

type WebhookConsumeSAMLAuthResponseResponse map[string]interface{}
