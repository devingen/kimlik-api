package token_service

type Scope string

type TokenPayload struct {
	Version   string `bson:"ver,omitempty"`
	Expires   string `bson:"exp,omitempty"`
	UserId    string `bson:"userId,omitempty"`
	SessionId string `bson:"sessionId,omitempty"`
	Scopes    string `bson:"scopes,omitempty"`
}

type ITokenService interface {
	GenerateToken(userId, sessionId string, scopes []Scope, duration int32) (string, error)
	ParseToken(accessToken string) (*TokenPayload, error)
}
