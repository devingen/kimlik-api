package service_controller

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"

	core "github.com/devingen/api-core"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
)

const JWTSigningKeyID = "89ce3598c473af1bda4bff95e6c8736450206fba"
const AccessTokenExpirationTime = 240 * time.Hour
const IDTokenExpirationTime = 1 * time.Hour

// OAuth2Token authenticates user and returns access token with given grant type.
func (c ServiceController) OAuth2Token(ctx context.Context, req core.Request) (*core.Response, error) {
	// TODO check client ID and secrets
	// TODO get these values as "Content-Type: application/x-www-form-urlencoded" ?????

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var params dto.OAuth2TokenRequest
	err := req.AssertBody(&params)
	if err != nil {
		return nil, err
	}

	switch params.GrantType {
	case dto.OAuth2GrantTypeAuthorizationCode:
		return c.handleGrantTypeAuthorizationCode(ctx, req, base, params)
	case dto.OAuth2GrantTypeOIDC:
		return c.handleGrantTypeKimlikOIDC(ctx, req, base, params)
	case dto.OAuth2GrantTypePassword:
		return c.handleGrantTypePassword(ctx, req, base, params)
	case dto.OAuth2GrantTypeRefreshToken:
		return c.handleGrantTypeRefreshToken(ctx, base, params)
	}

	return nil, core.NewError(http.StatusBadRequest, "invalid-grant-type")
}

func (c ServiceController) handleGrantTypePassword(ctx context.Context, req core.Request, base string, params dto.OAuth2TokenRequest) (*core.Response, error) {
	if params.Username == nil {
		return nil, core.NewError(http.StatusBadRequest, "username-missing")
	}
	auth, user, err := c.validateSessionWithPassword(ctx, base, *params.Username, *params.Password)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := c.createSuccessfulSessionAndGenerateToken(ctx, req, base, auth, user)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: dto.OAuth2TokenResponse{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    AccessTokenExpirationTime.Seconds(),
			RefreshToken: refreshToken,
		},
	}, nil
}

func (c ServiceController) handleGrantTypeAuthorizationCode(ctx context.Context, req core.Request, base string, params dto.OAuth2TokenRequest) (*core.Response, error) {

	accessCodes, err := c.DataService.FindOAuth2AccessCodes(ctx, base, bson.M{
		"code": params.Code,
	})
	if err != nil {
		return nil, err
	}

	if len(accessCodes) == 0 {
		return nil, core.NewError(http.StatusUnauthorized, "invalid-code")
	}
	accessCode := accessCodes[0]

	// Validate the code verifier against the code challenge
	if !verifyCodeChallenge(*params.CodeVerifier, accessCode.CodeChallenge) {
		return nil, core.NewError(http.StatusUnauthorized, "pkce-validation-failed")
	}

	tenantInfo, err := c.DataService.GetTenantInfo(ctx, base)
	if err != nil {
		return nil, err
	}

	if tenantInfo == nil {
		return nil, core.NewError(http.StatusNotFound, "tenant-info-not-initialised")
	}

	if tenantInfo.OAuth2IssuerIdentifier == nil {
		return nil, core.NewError(http.StatusNotAcceptable, "tenant-info-oauth2-issuer-identifier-missing")
	}

	user, err := c.DataService.FindUserWithId(ctx, base, accessCode.CreatedBy.ID.Hex())
	if err != nil {
		return nil, err
	}

	idToken, err := generateIDToken(*tenantInfo.OAuth2IssuerIdentifier, *tenantInfo.OAuth2SigningPrivateKey, *accessCode.ClientID, user)
	if err != nil {
		return nil, err
	}

	// TODO create session and generate access token for GET /userinfo
	return &core.Response{
		StatusCode: http.StatusOK,
		Body: dto.OAuth2TokenResponse{
			IDToken: idToken,
		},
	}, nil
}

func (c ServiceController) handleGrantTypeKimlikOIDC(ctx context.Context, req core.Request, base string, params dto.OAuth2TokenRequest) (*core.Response, error) {
	givenName := ""
	if params.GivenName != nil {
		givenName = *params.GivenName
	}
	familyName := ""
	if params.FamilyName != nil {
		familyName = *params.FamilyName
	}
	auth, user, isNewUser, err := c.validateSessionWithIDToken(ctx, base, *params.IDToken, givenName, familyName)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := c.createSuccessfulSessionAndGenerateToken(ctx, req, base, auth, user)
	if err != nil {
		return nil, err
	}

	responseStatusCode := http.StatusOK
	if isNewUser {
		responseStatusCode = http.StatusCreated
	}

	return &core.Response{
		StatusCode: responseStatusCode,
		Body: dto.OAuth2TokenResponse{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    AccessTokenExpirationTime.Seconds(),
			RefreshToken: refreshToken,
		},
	}, nil
}

func (c ServiceController) handleGrantTypeRefreshToken(ctx context.Context, base string, params dto.OAuth2TokenRequest) (*core.Response, error) {
	sessions, err := c.DataService.FindSessions(ctx, base, bson.M{
		"refreshToken": c.TokenService.HashRefreshToken(*params.RefreshToken),
	})
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, core.NewStatusError(http.StatusUnauthorized)
	}
	session := sessions[0]

	refreshToken, err := c.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	accessToken, err := c.generateAccessToken(session.User.ID.Hex(), session.ID.Hex())
	if err != nil {
		return nil, err
	}

	_, _, err = c.DataService.UpdateSession(ctx, base, &model.Session{
		ID:           session.ID,
		RefreshToken: &refreshToken.HashedToken,
	})
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: dto.OAuth2TokenResponse{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    AccessTokenExpirationTime.Seconds(),
			RefreshToken: refreshToken.RawToken,
		},
	}, nil
}

func verifyCodeChallenge(codeVerifier, codeChallenge string) bool {
	// Step 1: Recalculate the code challenge from the codeVerifier
	sha2 := sha256.New()
	io.WriteString(sha2, codeVerifier)
	calculatedCodeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))

	// Step 2: Compare the calculated code challenge with the one received
	// Remove padding from the received codeChallenge for correct comparison
	return calculatedCodeChallenge == codeChallenge
}

// Parse PEM-encoded private key
func ParsePrivateKeyFromPEM(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func generateIDToken(issuer, signingPrivateKey, audience string, user *model.User) (string, error) {
	pkcsPrivateKey, err := ParsePrivateKeyFromPEM([]byte(signingPrivateKey))
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := IDTokenClaims{
		Issuer:     issuer,
		Audience:   []string{audience},
		Subject:    user.ID.Hex(),
		Email:      *user.Email,
		Name:       *user.FirstName + " " + *user.LastName,
		GivenName:  *user.FirstName,
		FamilyName: *user.LastName,
		Expiration: now.Add(IDTokenExpirationTime).Unix(),
		IssuedAt:   now.Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = JWTSigningKeyID

	// Sign the token with the private RSA key
	tokenString, err := token.SignedString(pkcsPrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
