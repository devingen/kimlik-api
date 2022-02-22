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

	keyMappings := defaultValueKeyMapping
	if samlConfig.AttributeKeyMappingTemplate == model.AttributeKeyMappingTemplateAzureAD {
		keyMappings = azureADValueKeyMapping
	}

	body := dto.WebhookConsumeSAMLAuthResponseRequest{
		User: dto.WebhookConsumeSAMLAuthResponseRequestUser{
			Email:     assertionInfo.Values.Get(keyMappings.Email),
			FirstName: assertionInfo.Values.Get(keyMappings.FirstName),
			LastName:  assertionInfo.Values.Get(keyMappings.LastName),
			Meta:      map[string]interface{}{},
		},
		QueryParams: req.QueryStringParameters,
	}

	if samlConfig.MetaAttributeKeyMapping != nil {
		// pass the extra fields if defined in the config
		for keyInMeta, keyInAssertionInfo := range samlConfig.MetaAttributeKeyMapping {
			if _, has := assertionInfo.Values[keyInAssertionInfo]; has {
				body.User.Meta[keyInMeta] = assertionInfo.Values.Get(keyInAssertionInfo)
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

// ValueKeyMapping is used for parsing the SAML data from different identity providers.
type ValueKeyMapping struct {
	Email     string
	FirstName string
	LastName  string
}

var defaultValueKeyMapping = &ValueKeyMapping{
	Email:     "Email",
	FirstName: "FirstName",
	LastName:  "LastName",
}

// See https://docs.microsoft.com/en-us/azure/active-directory/develop/reference-saml-tokens
var azureADValueKeyMapping = &ValueKeyMapping{
	Email:     "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress",
	FirstName: "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname",
	LastName:  "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname",
}
