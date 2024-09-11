package service_controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type stringAsBool bool

func (sb *stringAsBool) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "true", `"true"`:
		*sb = true
	case "false", `"false"`:
		*sb = false
	default:
		return errors.New("invalid value for boolean")
	}
	return nil
}

type audience []string

func (a *audience) UnmarshalJSON(b []byte) error {
	var s string
	if json.Unmarshal(b, &s) == nil {
		*a = audience{s}
		return nil
	}
	var auds []string
	if err := json.Unmarshal(b, &auds); err != nil {
		return err
	}
	*a = auds
	return nil
}

// IDTokenClaims includes the fields returned in IDToken
// See https://openid.net/specs/openid-connect-core-1_0.html
type IDTokenClaims struct {
	Issuer     string   `json:"iss"`
	Audience   audience `json:"aud"`
	Subject    string   `json:"sub"`
	Email      string   `json:"email"`
	GivenName  string   `json:"given_name"`
	FamilyName string   `json:"family_name"`

	// Handle providers that return email_verified as a string
	// https://forums.aws.amazon.com/thread.jspa?messageID=949441&#949441 and
	// https://discuss.elastic.co/t/openid-error-after-authenticating-against-aws-cognito/206018/11
	EmailVerified stringAsBool `json:"email_verified"`
}

type IDToken struct {
	RawIDToken string
	ClientID   string
	Claims     IDTokenClaims
}

func (t *IDToken) Parse() error {
	payload, err := parseJWT(t.RawIDToken)
	if err != nil {
		return fmt.Errorf("oidc: malformed jwt: %v", err)
	}
	if err := json.Unmarshal(payload, &t.Claims); err != nil {
		return fmt.Errorf("oidc: failed to unmarshal claims: %v", err)
	}
	return nil
}

func (t *IDToken) Verify(ctx context.Context) error {

	provider, err := oidc.NewProvider(ctx, t.Claims.Issuer)
	if err != nil {
		return err
	}

	// Parse and verify ID Token payload.
	_, err = provider.Verifier(&oidc.Config{ClientID: t.Claims.Audience[0]}).Verify(ctx, t.RawIDToken)
	if err != nil {
		return err
	}
	return nil
}

func (tc *IDTokenClaims) IsEmailVerified() bool {
	return bool(tc.EmailVerified)
}

// ToMap returns the minimum required info for creating auth in database
func (tc *IDTokenClaims) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"iss":   tc.Issuer,
		"aud":   tc.Audience[0],
		"sub":   tc.Subject,
		"email": tc.Email,
	}
}

func parseJWT(p string) ([]byte, error) {
	parts := strings.Split(p, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("oidc: malformed jwt, expected 3 parts got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("oidc: malformed jwt payload: %v", err)
	}
	return payload, nil
}
