package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	core_dto "github.com/devingen/api-core/dto"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/model"
	"net/http"
)

func (c ServiceController) GetTenantInfo(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	item, err := c.DataService.GetTenantInfo(ctx, base)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       item,
	}, nil
}

func (c ServiceController) UpdateTenantInfo(ctx context.Context, req core.Request) (*core.Response, error) {

	// TODO enable webhook to check if the user has permissions to do this
	//_, interceptorStatusCode, interceptorError := c.InterceptorService.Pre(ctx, req)
	//if interceptorError != nil {
	//	return &core.Response{
	//		StatusCode: interceptorStatusCode,
	//		Body:       interceptorError,
	//	}, nil
	//}

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	item, err := c.DataService.GetTenantInfo(ctx, base)
	if err != nil {
		return nil, err
	}

	var body model.TenantInfo
	_, err = kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, err
	}

	updatedAt, revision, err := c.DataService.UpdateTenantInfo(ctx, base, &model.TenantInfo{
		ID:               item.ID,
		Name:             body.Name,
		LogoURL:          body.LogoURL,
		TermsOfUseURL:    body.TermsOfUseURL,
		PrivacyPolicyURL: body.PrivacyPolicyURL,
		SupportURL:       body.SupportURL,
		SupportEmail:     body.SupportEmail,
	})
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: core_dto.UpdateEntryResponse{
			ID:        item.ID,
			UpdatedAt: *updatedAt,
			Revision:  revision,
		},
	}, nil
}
