package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	core_dto "github.com/devingen/api-core/dto"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/model"
)

func (c ServiceController) GetIntegrationSettings(ctx context.Context, req core.Request) (*core.Response, error) {

	base := req.PathParameters["base"]

	integrationSettings, err := c.DataService.GetIntegrationSettings(ctx, base)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       integrationSettings,
	}, nil
}

func (c ServiceController) UpdateIntegrationSettings(ctx context.Context, req core.Request) (*core.Response, error) {

	base := req.PathParameters["base"]

	var body model.IntegrationSettings
	_, err := kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, err
	}

	// Get existing settings to preserve ID
	existingSettings, err := c.DataService.GetIntegrationSettings(ctx, base)
	if err != nil {
		return nil, err
	}

	body.ID = existingSettings.ID
	updatedAt, revision, err := c.DataService.UpdateIntegrationSettings(ctx, base, &body)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: core_dto.UpdateEntryResponse{
			ID:        body.ID,
			UpdatedAt: *updatedAt,
			Revision:  revision,
		},
	}, nil
}
