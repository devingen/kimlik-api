package webhookis

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
)

func (service WebhookInterceptorService) Final(ctx context.Context, req core.Request, responseBody interface{}) {
	if service.Client == nil {
		return
	}

	service.Client.R().EnableTrace().
		SetBody(dto.WebhookFinalRequest{
			Method:         req.HTTPMethod,
			Path:           req.Path,
			PathParameters: req.PathParameters,
			QueryParams:    req.QueryStringParameters,
			Header:         req.Headers,
			ResponseBody:   responseBody,
		}).
		Post("/final")
	return
}
