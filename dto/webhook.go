package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type WebhookPreRequest struct {
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	QueryParams map[string][]string    `json:"queryParams"`
	Header      map[string]string      `json:"header"`
	Body        map[string]interface{} `json:"body"`
}

type WebhookPreResponse struct {

	// QueryEnhance is used to add extra fields to the query for GET list requests.
	QueryEnhance *QueryEnhance `json:"queryEnhance"`
}

type QueryEnhance struct {
	// IDsIn filters the returned items to contain only the given IDs. All items are returned otherwise.
	IDsIn []primitive.ObjectID `json:"idsIn"`
}

type WebhookFinalRequest struct {
	Method         string              `json:"method"`
	Path           string              `json:"path"`
	PathParameters map[string]string   `json:"pathParameters"`
	QueryParams    map[string][]string `json:"queryParams"`
	Header         map[string]string   `json:"header"`
	ResponseBody   interface{}         `json:"responseBody"`
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
