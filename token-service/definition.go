package token_service

import "github.com/devingen/api-core/dvnruntime"

type Scope string

type TokenPayload struct {
	Version   string `bson:"ver,omitempty"`
	Expires   string `bson:"exp,omitempty"`
	UserId    string `bson:"userId,omitempty"`
	SessionId string `bson:"sessionId,omitempty"`
	Scopes    string `bson:"scopes,omitempty"`
}

const (
	ContextKeyTokenPayload dvnruntime.ContextKey = "context-key-token-payload"
)

type ITokenService interface {
	GenerateToken(userId, sessionId string, scopes []Scope, duration int32) (string, error)
	ParseToken(accessToken string) (*TokenPayload, error)
}
