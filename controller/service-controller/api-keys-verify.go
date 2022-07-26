package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/model"
	"net/http"
)

func (c ServiceController) VerifyAPIKey(ctx context.Context, req core.Request) (*core.Response, error) {

	base, err := req.AssertPathParameter("base")
	if err != nil {
		return nil, err
	}

	apiKey, err := kimlik.AssertApiKey(ctx, base, c.DataService)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: model.APIKey{
			ID:        apiKey.ID,
			CreatedAt: apiKey.CreatedAt,
			UpdatedAt: apiKey.UpdatedAt,
			Revision:  apiKey.Revision,
			CreatedBy: apiKey.CreatedBy,
			Name:      apiKey.Name,
			KeyID:     apiKey.KeyID,
			Scopes:    apiKey.Scopes,
		},
	}, nil
}
