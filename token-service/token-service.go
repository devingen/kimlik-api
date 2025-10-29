package token_service

type TokenPayload struct {
	Version   string   `bson:"ver,omitempty"`
	Expires   string   `bson:"exp,omitempty"`
	UserID    string   `bson:"userId,omitempty"`
	SessionID string   `bson:"sessionId,omitempty"`
	Scopes    []string `bson:"scopes,omitempty"`
}

func (token *TokenPayload) ContainsScope(scope string) bool {
	for _, s := range token.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

type RefreshToken struct {
	// HashedToken should be saved in database.
	HashedToken string

	// RawToken should be returned to client.
	RawToken string
}

type ITokenService interface {
	// GenerateRefreshToken generates an opaque token and returns the raw and hashed tokens. The raw token
	// should be returned to the client and the hashed token must be stored in a secure place.
	GenerateRefreshToken() (*RefreshToken, error)

	// HashRefreshToken returns the hashed token to be able to find it in database.
	HashRefreshToken(token string) string

	// GenerateAccessToken returns access token with given user and session details.
	GenerateAccessToken(userId, sessionId string, scopes []string, exp int64) (string, error)

	// ParseAccessToken validates the token and returns the token payload.
	ParseAccessToken(accessToken string) (*TokenPayload, error)

	// GenerateAuthorizationCode generates random authorization code.
	GenerateAuthorizationCode() (*string, error)
}
