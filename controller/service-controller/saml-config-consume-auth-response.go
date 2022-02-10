package service_controller

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/log"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"

	saml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
)

func (controller ServiceController) ConsumeSAMLAuthResponse(ctx context.Context, req core.Request) (interface{}, int, error) {

	fmt.Println("ConsumeSAMLAuthResponse", 1)
	loggerFromContext, err := log.Of(ctx)
	if err != nil {
		fmt.Println("ConsumeSAMLAuthResponse", 2)
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-logger-in-context")
	}
	logger := loggerFromContext.WithFields(logrus.Fields{
		"function": "consume-saml-auth-response",
	})

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		logger.WithFields(logrus.Fields{"error": err}).Error("missing-path-param-base")
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	samlConfigID, hasSamlConfigID := req.PathParameters["id"]
	if !hasSamlConfigID {
		logger.WithFields(logrus.Fields{"error": err}).Error("missing-path-param-saml-config-id")
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-path-param-saml-config-id")
	}

	var body dto.ConsumeSAMLAuthResponseRequest
	err = req.AssertBody(&body)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Error("invalid-request-body")
		return nil, 0, err
	}

	samlConfig, err := controller.DataService.GetSAMLConfig(ctx, base, samlConfigID)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Error("failed-to-get-saml-config")
		return nil, 0, err
	}

	sp, err := getSp(samlConfig)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Error("failed-to-get-sp")
		return nil, 0, err
	}

	assertionInfo, err := sp.RetrieveAssertionInfo(*body.SAMLResponse)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Error("saml-response-assertion-failed")
		return nil, 0, err
	}

	if assertionInfo.WarningInfo.InvalidTime {
		logger.WithFields(logrus.Fields{"error": err}).Error("invalid-saml-response-time")
		return nil, 0, err
	}

	if assertionInfo.WarningInfo.NotInAudience {
		logger.WithFields(logrus.Fields{"error": err}).Error("invalid-saml-response-audience")
		return nil, 0, err
	}

	interceptorResponse, interceptorStatusCode, interceptorError := controller.InterceptorService.SAMLConsume(ctx, req, samlConfig, assertionInfo)
	if interceptorError != nil {
		logger.WithFields(logrus.Fields{"error": interceptorError}).Error("interceptor-returned-error")
		return interceptorError, interceptorStatusCode, nil
	}

	return interceptorResponse, http.StatusOK, nil
}

func getSp(config *model.SAMLConfig) (*saml2.SAMLServiceProvider, error) {

	var rawMetadata []byte
	if config.MetadataContent != nil {
		rawMetadata = []byte(*config.MetadataContent)
	} else if config.MetadataURL != nil {
		res, err := http.Get(*config.MetadataURL)
		if err != nil {
			return nil, err
		}

		rawMetadata, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
	}

	metadata := &types.EntityDescriptor{}
	err := xml.Unmarshal(rawMetadata, metadata)
	if err != nil {
		return nil, err
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
				return nil, err
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				return nil, err
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}

	// TODO We sign the AuthnRequest with a random key because Okta doesn't seem to verify these.
	randomKeyStore := dsig.RandomKeyStoreForTest()

	sp := &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:      metadata.IDPSSODescriptor.SingleSignOnServices[0].Location,
		IdentityProviderIssuer:      metadata.EntityID,
		ServiceProviderIssuer:       *config.ServiceProviderIssuer,
		AssertionConsumerServiceURL: *config.AssertionConsumerServiceURL,
		SignAuthnRequests:           true,
		AudienceURI:                 *config.AudienceURI,
		IDPCertificateStore:         &certStore,
		SPKeyStore:                  randomKeyStore,
	}
	return sp, err
}
