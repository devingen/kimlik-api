package is

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	saml2 "github.com/russellhaering/gosaml2"
)

// IKimlikInterceptorService defines the functionality of the interceptors
type IKimlikInterceptorService interface {
	Pre(ctx context.Context, req core.Request) (*dto.WebhookPreResponse, int, interface{})
	Final(ctx context.Context, req core.Request, responseBody interface{})
	SAMLConsume(ctx context.Context, req core.Request, samlConfig *model.SAMLConfig, assertionInfo *saml2.AssertionInfo) (*dto.WebhookConsumeSAMLAuthResponseResponse, int, interface{})
}
