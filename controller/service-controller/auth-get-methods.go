package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
)

type OAuth2Provider struct {
	Issuer   string `json:"issuer"`
	Name     string `json:"name"`
	ClientID string `json:"clientId"`
}

type GetAuthMethodsResponse struct {
	Methods         []string         `json:"methods"`
	OAuth2Providers []OAuth2Provider `json:"oauth2Providers,omitempty"`
}

func (c ServiceController) GetAuthMethods(ctx context.Context, req core.Request) (*core.Response, error) {
	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	email, err := req.AssertQueryStringParameter("email")
	if err != nil {
		return nil, err
	}

	// Find user by email
	user, err := c.DataService.FindUserWithEmail(ctx, base, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, core.NewError(http.StatusNotFound, "user-not-found")
	}

	// Find all auth methods for this user
	auths, err := c.DataService.FindAuthsOfUser(ctx, base, user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// Get all OAuth2 configs for this base to match issuers
	oauth2Configs, err := c.DataService.FindOAuth2Configs(ctx, base, nil)
	if err != nil {
		return nil, err
	}

	// Create a map of issuer -> OAuth2Config for quick lookup
	oauth2ConfigMap := make(map[string]*struct {
		Name     string
		ClientID string
	})
	for _, config := range oauth2Configs {
		if config.Issuer != nil && config.Name != nil && config.ClientID != nil {
			oauth2ConfigMap[*config.Issuer] = &struct {
				Name     string
				ClientID string
			}{
				Name:     *config.Name,
				ClientID: *config.ClientID,
			}
		}
	}

	// Extract unique auth types and OAuth2 providers
	methodsMap := make(map[string]bool)
	oauth2Providers := []OAuth2Provider{}
	oauth2ProvidersMap := make(map[string]bool) // To track unique providers by issuer

	for _, auth := range auths {
		methodsMap[string(auth.Type)] = true

		// If this is an OpenID auth, collect provider information
		if auth.Type == "openid" && auth.OpenID != nil && auth.OpenID.Iss != "" {
			issuer := auth.OpenID.Iss

			// Only add unique providers and only if we have a matching OAuth2 config
			if !oauth2ProvidersMap[issuer] {
				if configInfo, exists := oauth2ConfigMap[issuer]; exists {
					oauth2ProvidersMap[issuer] = true
					oauth2Providers = append(oauth2Providers, OAuth2Provider{
						Issuer:   issuer,
						Name:     configInfo.Name,
						ClientID: configInfo.ClientID,
					})
				}
			}
		}
	}

	// Convert map to slice
	methods := make([]string, 0, len(methodsMap))
	for method := range methodsMap {
		methods = append(methods, method)
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: GetAuthMethodsResponse{
			Methods:         methods,
			OAuth2Providers: oauth2Providers,
		},
	}, nil
}
