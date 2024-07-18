package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
)

func (c ServiceController) LoginWithSAML(ctx context.Context, req core.Request) (*core.Response, error) {

	buildSAMLAuthURLResponse, err := c.BuildSAMLAuthURL(ctx, req)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusFound,
		Headers: map[string]string{
			"Location": core.StringValue(buildSAMLAuthURLResponse.Body.(dto.BuildSAMLAuthURLResponse).AuthURL),
		},
	}, nil
}
