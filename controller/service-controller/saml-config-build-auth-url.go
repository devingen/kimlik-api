package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	"net/http"
)

func (c ServiceController) BuildSAMLAuthURL(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	samlConfigID, hasSamlConfigID := req.PathParameters["id"]
	if !hasSamlConfigID {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-saml-config-id")
	}

	samlConfig, err := c.DataService.GetSAMLConfig(ctx, base, samlConfigID)
	if err != nil {
		return nil, err
	}

	authURL, err := buildSAMLAuthURL(samlConfig)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       dto.BuildSAMLAuthURLResponse{AuthURL: authURL},
	}, nil
}

func buildSAMLAuthURL(config *model.SAMLConfig) (*string, error) {
	sp, err := getSp(config)
	if err != nil {
		return nil, err
	}

	authURL, err := sp.BuildAuthURL("")
	return &authURL, err
}
