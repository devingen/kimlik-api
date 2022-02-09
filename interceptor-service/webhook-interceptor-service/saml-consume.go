package webhookis

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	saml2 "github.com/russellhaering/gosaml2"
	"net/http"
	"net/url"
)

func (service WebhookInterceptorService) SAMLConsume(ctx context.Context, req core.Request, samlConfig *model.SAMLConfig, assertionInfo *saml2.AssertionInfo) (*dto.WebhookConsumeSAMLAuthResponseResponse, int, interface{}) {
	if service.Client == nil {
		return nil, 0, nil
	}

	body := dto.WebhookConsumeSAMLAuthResponseRequest{
		User: dto.WebhookConsumeSAMLAuthResponseRequestUser{
			Email:     assertionInfo.Values.Get("Email"),
			FirstName: assertionInfo.Values.Get("FirstName"),
			LastName:  assertionInfo.Values.Get("LastName"),
			Meta:      map[string]interface{}{},
		},
		QueryParams: req.QueryStringParameters,
	}

	if samlConfig.SAMLResponseValues != nil {
		// pass the extra fields if defined in the config
		for _, key := range samlConfig.SAMLResponseValues {
			if _, has := assertionInfo.Values[key]; has {
				body.User.Meta[key] = assertionInfo.Values.Get(key)
			}
		}
	}

	resp, err := service.Client.R().EnableTrace().
		SetBody(body).
		SetResult(&dto.WebhookConsumeSAMLAuthResponseResponse{}).
		SetError(&map[string]interface{}{}).
		Post("/saml/consume")
	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, http.StatusInternalServerError, core.NewError(http.StatusInternalServerError, "webhook-api-is-unreachable")
		}
		return nil, resp.StatusCode(), err
	}
	if resp.StatusCode() > 399 {
		return nil, resp.StatusCode(), resp.Error()
	}

	webhookResponse := resp.Result().(*dto.WebhookConsumeSAMLAuthResponseResponse)
	return webhookResponse, resp.StatusCode(), nil
}
