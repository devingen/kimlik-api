package service_controller

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/log"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	"github.com/sirupsen/logrus"

	saml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
)

func (c ServiceController) ConsumeSAMLAuthResponse(ctx context.Context, req core.Request) (*core.Response, error) {

	loggerFromContext, err := log.Of(ctx)
	if err != nil {
		return nil, core.NewError(http.StatusInternalServerError, "missing-logger-in-context")
	}
	logger := loggerFromContext.WithFields(logrus.Fields{
		"function": "consume-saml-auth-response",
	})

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		logger.WithFields(logrus.Fields{"error": err}).Error("missing-path-param-base")
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}
	logger = loggerFromContext.WithFields(logrus.Fields{
		"base": base,
	})

	samlConfigID, hasSamlConfigID := req.PathParameters["id"]
	if !hasSamlConfigID {
		logger.WithFields(logrus.Fields{"error": err}).Error("missing-path-param-saml-config-id")
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-saml-config-id")
	}
	logger = loggerFromContext.WithFields(logrus.Fields{
		"samlConfigId": samlConfigID,
	})

	var body dto.ConsumeSAMLAuthResponseRequest
	err = req.AssertBody(&body)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Error("invalid-request-body")
		return nil, err
	}

	samlConfig, err := c.DataService.GetSAMLConfig(ctx, base, samlConfigID)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Error("failed-to-get-saml-config")
		return nil, err
	}

	sp, err := getSp(samlConfig)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Error("failed-to-get-sp")
		return nil, err
	}

	assertionInfo, err := sp.RetrieveAssertionInfo(*body.SAMLResponse)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Error("saml-response-assertion-failed")
		return nil, err
	}

	if assertionInfo.WarningInfo.InvalidTime {
		logger.WithFields(logrus.Fields{"error": err}).Error("invalid-saml-response-time")
		return nil, err
	}

	if assertionInfo.WarningInfo.NotInAudience {
		logger.WithFields(logrus.Fields{"error": err}).Error("invalid-saml-response-audience")
		return nil, err
	}

	interceptorResponse, interceptorStatusCode, interceptorError := c.InterceptorService.SAMLConsume(ctx, req, samlConfig, assertionInfo)
	if interceptorError != nil {
		logger.WithFields(logrus.Fields{"error": interceptorError}).Error("interceptor-returned-error")
		return &core.Response{
			StatusCode: interceptorStatusCode,
			Body:       interceptorError,
		}, nil
	}

	return &core.Response{
		StatusCode: interceptorStatusCode,
		Body:       interceptorResponse,
	}, nil
}

func getSp(config *model.SAMLConfig) (*saml2.SAMLServiceProvider, error) {

	var rawMetadata []byte
	if config.MetadataContent != nil {
		rawMetadata = []byte(*config.MetadataContent)
	} else if config.MetadataURL != nil {
		res, err := http.Get(*config.MetadataURL)
		if err != nil {
			return nil, core.NewError(http.StatusPreconditionFailed, "failed-to-fetch-metadata reason:"+err.Error())
		}

		rawMetadata, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, core.NewError(http.StatusPreconditionFailed, "failed-to-parse-metadata reason:"+err.Error())
		}
	}

	metadata := &types.EntityDescriptor{}
	err := xml.Unmarshal(rawMetadata, metadata)
	if err != nil {
		return nil, core.NewError(http.StatusPreconditionFailed, "failed-to-read-idp-metadata, from: "+*config.MetadataURL+", reason:"+err.Error())
	}

	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				return nil, fmt.Errorf("metadata certificate(%d) must not be empty", idx)
			}

			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				return nil, core.NewError(http.StatusPreconditionFailed, "failed-to-decode-xcert-data reason:"+err.Error())
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				return nil, core.NewError(http.StatusPreconditionFailed, "failed-to-parse-idp-cert reason:"+err.Error())
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}

	// TODO We sign the AuthnRequest with a random key because Okta doesn't seem to verify these.
	randomKeyStore := dsig.RandomKeyStoreForTest()

	sp := &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:      metadata.IDPSSODescriptor.SingleSignOnServices[0].Location,
		IdentityProviderIssuer:      metadata.EntityID,
		ServiceProviderIssuer:       *config.AudienceURI,
		AssertionConsumerServiceURL: *config.AssertionConsumerServiceURL,
		SignAuthnRequests:           true,
		AudienceURI:                 *config.AudienceURI,
		IDPCertificateStore:         &certStore,
		SPKeyStore:                  randomKeyStore,
	}
	return sp, err
}
