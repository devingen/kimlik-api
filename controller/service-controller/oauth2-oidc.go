package service_controller

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"

	core "github.com/devingen/api-core"
)

type OpenIDConfiguration struct {
	Issuer                 string   `json:"issuer"`
	AuthorizationEndpoint  string   `json:"authorization_endpoint"`
	TokenEndpoint          string   `json:"token_endpoint"`
	UserinfoEndpoint       string   `json:"userinfo_endpoint"`
	JwksUri                string   `json:"jwks_uri"`
	RegistrationEndpoint   string   `json:"registration_endpoint"`
	ScopesSupported        []string `json:"scopes_supported"`
	ResponseTypesSupported []string `json:"response_types_supported"`
	GrantTypesSupported    []string `json:"grant_types_supported"`
}

func (c ServiceController) OAuth2GetOIDCConfiguration(ctx context.Context, req core.Request) (*core.Response, error) {
	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	tenantInfo, err := c.DataService.GetTenantInfo(ctx, base)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: OpenIDConfiguration{
			Issuer:                 *tenantInfo.OAuth2IssuerIdentifier,
			AuthorizationEndpoint:  *tenantInfo.OAuth2AuthorizationURL,
			JwksUri:                "https://api.kimlik.das.devingen.io/" + base + "/oauth2/certs",
			TokenEndpoint:          "https://api.kimlik.das.devingen.io/" + base + "/oauth2/token",
			UserinfoEndpoint:       "https://api.kimlik.das.devingen.io/" + base + "/oauth2/userinfo",
			RegistrationEndpoint:   "https://api.kimlik.das.devingen.io/" + base + "/oauth2/register",
			ScopesSupported:        []string{"openid", "profile", "email"},
			ResponseTypesSupported: []string{"code", "id_token", "token"},
			GrantTypesSupported:    []string{"refresh_token", "authorization_code", "client_credentials"},
		},
	}, nil
}

func (c ServiceController) OAuth2GetJWKS(ctx context.Context, req core.Request) (*core.Response, error) {
	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	tenantInfo, err := c.DataService.GetTenantInfo(ctx, base)
	if err != nil {
		return nil, err
	}

	key, err := ConvertPEMPublicKeyToJWK(*tenantInfo.OAuth2SigningPublicKey)
	if err != nil {
		return nil, core.NewError(http.StatusInternalServerError, err.Error())
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: map[string]interface{}{
			"keys": []map[string]interface{}{key},
		},
	}, nil
}

func ConvertPEMPublicKeyToJWK(pemPublicKey string) (map[string]interface{}, error) {
	// Decode the PEM encoded public key
	block, _ := pem.Decode([]byte(pemPublicKey))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	// Parse the public key
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Ensure the key is of type RSA
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of RSA type")
	}

	// Prepare the JWK representation
	jwk := map[string]interface{}{
		"kty": "RSA",
		"alg": "RS256",
		"e":   "AQAB",
		"kid": JWTSigningKeyID,
		"n":   base64.RawURLEncoding.EncodeToString(rsaPubKey.N.Bytes()),
	}

	return jwk, nil
}

type SigningKey struct {
	Key   string `json:"kty"`
	KeyID string `json:"kid"`
	Use   string `json:"use"`
	Alg   string `json:"alg"`
	N     string `json:"n"`
	E     string `json:"e"`
}
