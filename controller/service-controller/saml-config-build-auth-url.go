package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	"net/http"
)

func (controller ServiceController) BuildSAMLAuthURL(ctx context.Context, req core.Request) (interface{}, int, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	samlConfigID, hasSamlConfigID := req.PathParameters["id"]
	if !hasSamlConfigID {
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-path-param-saml-config-id")
	}

	samlConfig, err := controller.DataService.GetSAMLConfig(ctx, base, samlConfigID)
	if err != nil {
		return nil, 0, err
	}

	authURL, err := buildSAMLAuthURL(samlConfig)
	if err != nil {
		return nil, 0, err
	}

	return dto.BuildSAMLAuthURLResponse{AuthURL: authURL}, http.StatusOK, nil
}

func buildSAMLAuthURL(config *model.SAMLConfig) (*string, error) {
	sp, err := getSp(config)
	if err != nil {
		return nil, err
	}

	authURL, err := sp.BuildAuthURL("")
	return &authURL, err
}
