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
	res, err := http.Get(config.MetadataURL)
	//res, err := http.Get("http://idp.oktadev.com/metadata")
	// res, err := http.Get("https://idp-ac-2-cs4g5.ondigitalocean.app/metadata")
	if err != nil {
		return nil, err
	}

	rawMetadata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	//rawMetadata = []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?><EntityDescriptor ID=\"_88e2b037-f946-43da-8495-dae3cef4ae57\" entityID=\"https://sts.windows.net/58f82789-6695-4f4a-abdb-357668d55cff/\" xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\"><Signature xmlns=\"http://www.w3.org/2000/09/xmldsig#\"><SignedInfo><CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\" /><SignatureMethod Algorithm=\"http://www.w3.org/2001/04/xmldsig-more#rsa-sha256\" /><Reference URI=\"#_88e2b037-f946-43da-8495-dae3cef4ae57\"><Transforms><Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\" /><Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\" /></Transforms><DigestMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#sha256\" /><DigestValue>CNVD9vtHLyYw9smPNqV9UFDN7OLqQ+ZObT59XQadukw=</DigestValue></Reference></SignedInfo><SignatureValue>VrSDXpu8B19NRS+F87Ie+uAxpggpOBOj3yXQI7qvRn7ySXb5bl4xq1mlCEPcOlm3t9f//IKBnXWFL9ui0tPoV3kpD9qsVrD3CqJn3PdpVDdfohY8fr9s4/erWq1M5lUsqmHWVeDXnlh8fR5z74hkCkmc/wmJc8WlkL3CTLAgbIuS1d5l9x6VgrCbo3zaV3cZMOyRvsQVU7+ddWcUl2Y0oXN19GUIX738Zhb30oy+Qi3K3QNEcrjBM1Pd5ykB1ph6Lp+tIAofveCSi9aCk4T0PqoRA/BVmsJ4sRc4h0dm56GHOoLH5IbcHERuh7AiZJLwWt3/c1QvyBhhZU5NzTdRKw==</SignatureValue><KeyInfo><X509Data><X509Certificate>MIIC8DCCAdigAwIBAgIQSP5mK0sN6IJGzefNl8rxQDANBgkqhkiG9w0BAQsFADA0MTIwMAYDVQQDEylNaWNyb3NvZnQgQXp1cmUgRmVkZXJhdGVkIFNTTyBDZXJ0aWZpY2F0ZTAeFw0yMjAxMTEwMjM5NTRaFw0yNTAxMTEwMjM5NThaMDQxMjAwBgNVBAMTKU1pY3Jvc29mdCBBenVyZSBGZWRlcmF0ZWQgU1NPIENlcnRpZmljYXRlMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAv08WgEsT86mCgByv/G+c6OmazgVxcSo6HvOUD8iilfu+WisMsFfmitFLx/ZHKGiLwTzK4Q0/Ekv+tTmzLcIJx5EyY/hk6Md5KYGDx2bOzaRlFzVq3Sm2fAnyHZyB5qq7koHjYUaIW+sNLxt7Xz0isWVNcOElmC6RrP16x3lNpQVLaYyiDvPVwArfoLa8GrzeF9bbCWBLvKDKDP4tT7uP/GxgY8IgNRwbdnbtqpuvQlEBsZidjdlQckMP5bg0friicsHJ/gqJ2TV7rGbWk/BrUjrngSrW1CqLa/HEgmPw9WSJAMUKYjJ2zf2OXPUDqA6K1qs65SOPQEJULP3KIqQcIQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCf1X/RQdcoizdvjun4ezF0hI+tg0yfrYMsvtumsEodQXX+rjnZ1fN3pkOceo3wsiKVeTgCceLuXLUWGQ04HUot4SU2qYI5eY+0oLM978s8l4SIHrPkj8yTb1sFuVnkmG2UMc/77wW9mKBIR7QuaocVFCHyuyMwtt0KrEV5X4+zD6yaNGtpI8cWHcSgCpoEO8lt7kf65kFOannhn0PZ90sTDXgW7kPmBXwl8TRjJBvV/b3HLzl+ZgAVkosNeJmJifZYBds//mE5edz8GwGj9Q67XNxousk3k188nWe5tRREwAqo+gw5a8y6Quz1CrUGlp2g257MSXlyVkrzsdDB88aR</X509Certificate></X509Data></KeyInfo></Signature><RoleDescriptor xsi:type=\"fed:SecurityTokenServiceType\" protocolSupportEnumeration=\"http://docs.oasis-open.org/wsfed/federation/200706\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:fed=\"http://docs.oasis-open.org/wsfed/federation/200706\"><KeyDescriptor use=\"signing\"><KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\"><X509Data><X509Certificate>MIIC8DCCAdigAwIBAgIQSP5mK0sN6IJGzefNl8rxQDANBgkqhkiG9w0BAQsFADA0MTIwMAYDVQQDEylNaWNyb3NvZnQgQXp1cmUgRmVkZXJhdGVkIFNTTyBDZXJ0aWZpY2F0ZTAeFw0yMjAxMTEwMjM5NTRaFw0yNTAxMTEwMjM5NThaMDQxMjAwBgNVBAMTKU1pY3Jvc29mdCBBenVyZSBGZWRlcmF0ZWQgU1NPIENlcnRpZmljYXRlMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAv08WgEsT86mCgByv/G+c6OmazgVxcSo6HvOUD8iilfu+WisMsFfmitFLx/ZHKGiLwTzK4Q0/Ekv+tTmzLcIJx5EyY/hk6Md5KYGDx2bOzaRlFzVq3Sm2fAnyHZyB5qq7koHjYUaIW+sNLxt7Xz0isWVNcOElmC6RrP16x3lNpQVLaYyiDvPVwArfoLa8GrzeF9bbCWBLvKDKDP4tT7uP/GxgY8IgNRwbdnbtqpuvQlEBsZidjdlQckMP5bg0friicsHJ/gqJ2TV7rGbWk/BrUjrngSrW1CqLa/HEgmPw9WSJAMUKYjJ2zf2OXPUDqA6K1qs65SOPQEJULP3KIqQcIQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCf1X/RQdcoizdvjun4ezF0hI+tg0yfrYMsvtumsEodQXX+rjnZ1fN3pkOceo3wsiKVeTgCceLuXLUWGQ04HUot4SU2qYI5eY+0oLM978s8l4SIHrPkj8yTb1sFuVnkmG2UMc/77wW9mKBIR7QuaocVFCHyuyMwtt0KrEV5X4+zD6yaNGtpI8cWHcSgCpoEO8lt7kf65kFOannhn0PZ90sTDXgW7kPmBXwl8TRjJBvV/b3HLzl+ZgAVkosNeJmJifZYBds//mE5edz8GwGj9Q67XNxousk3k188nWe5tRREwAqo+gw5a8y6Quz1CrUGlp2g257MSXlyVkrzsdDB88aR</X509Certificate></X509Data></KeyInfo></KeyDescriptor><fed:ClaimTypesOffered><auth:ClaimType Uri=\"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Name</auth:DisplayName><auth:Description>The mutable display name of the user.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/nameidentifier\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Subject</auth:DisplayName><auth:Description>An immutable, globally unique, non-reusable identifier of the user that is unique to the application for which a token is issued.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Given Name</auth:DisplayName><auth:Description>First name of the user.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Surname</auth:DisplayName><auth:Description>Last name of the user.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/identity/claims/displayname\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Display Name</auth:DisplayName><auth:Description>Display name of the user.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/identity/claims/nickname\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Nick Name</auth:DisplayName><auth:Description>Nick name of the user.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/ws/2008/06/identity/claims/authenticationinstant\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Authentication Instant</auth:DisplayName><auth:Description>The time (UTC) when the user is authenticated to Windows Azure Active Directory.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/ws/2008/06/identity/claims/authenticationmethod\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Authentication Method</auth:DisplayName><auth:Description>The method that Windows Azure Active Directory uses to authenticate users.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/identity/claims/objectidentifier\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>ObjectIdentifier</auth:DisplayName><auth:Description>Primary identifier for the user in the directory. Immutable, globally unique, non-reusable.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/identity/claims/tenantid\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>TenantId</auth:DisplayName><auth:Description>Identifier for the user's tenant.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/identity/claims/identityprovider\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>IdentityProvider</auth:DisplayName><auth:Description>Identity provider for the user.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Email</auth:DisplayName><auth:Description>Email address of the user.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/ws/2008/06/identity/claims/groups\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Groups</auth:DisplayName><auth:Description>Groups of the user.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/identity/claims/accesstoken\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>External Access Token</auth:DisplayName><auth:Description>Access token issued by external identity provider.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/ws/2008/06/identity/claims/expiration\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>External Access Token Expiration</auth:DisplayName><auth:Description>UTC expiration time of access token issued by external identity provider.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/identity/claims/openid2_id\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>External OpenID 2.0 Identifier</auth:DisplayName><auth:Description>OpenID 2.0 identifier issued by external identity provider.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/claims/groups.link\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>GroupsOverageClaim</auth:DisplayName><auth:Description>Issued when number of user's group claims exceeds return limit.</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/ws/2008/06/identity/claims/role\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>Role Claim</auth:DisplayName><auth:Description>Roles that the user or Service Principal is attached to</auth:Description></auth:ClaimType><auth:ClaimType Uri=\"http://schemas.microsoft.com/ws/2008/06/identity/claims/wids\" xmlns:auth=\"http://docs.oasis-open.org/wsfed/authorization/200706\"><auth:DisplayName>RoleTemplate Id Claim</auth:DisplayName><auth:Description>Role template id of the Built-in Directory Roles that the user is a member of</auth:Description></auth:ClaimType></fed:ClaimTypesOffered><fed:SecurityTokenServiceEndpoint><wsa:EndpointReference xmlns:wsa=\"http://www.w3.org/2005/08/addressing\"><wsa:Address>https://login.microsoftonline.com/58f82789-6695-4f4a-abdb-357668d55cff/wsfed</wsa:Address></wsa:EndpointReference></fed:SecurityTokenServiceEndpoint><fed:PassiveRequestorEndpoint><wsa:EndpointReference xmlns:wsa=\"http://www.w3.org/2005/08/addressing\"><wsa:Address>https://login.microsoftonline.com/58f82789-6695-4f4a-abdb-357668d55cff/wsfed</wsa:Address></wsa:EndpointReference></fed:PassiveRequestorEndpoint></RoleDescriptor><RoleDescriptor xsi:type=\"fed:ApplicationServiceType\" protocolSupportEnumeration=\"http://docs.oasis-open.org/wsfed/federation/200706\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:fed=\"http://docs.oasis-open.org/wsfed/federation/200706\"><KeyDescriptor use=\"signing\"><KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\"><X509Data><X509Certificate>MIIC8DCCAdigAwIBAgIQSP5mK0sN6IJGzefNl8rxQDANBgkqhkiG9w0BAQsFADA0MTIwMAYDVQQDEylNaWNyb3NvZnQgQXp1cmUgRmVkZXJhdGVkIFNTTyBDZXJ0aWZpY2F0ZTAeFw0yMjAxMTEwMjM5NTRaFw0yNTAxMTEwMjM5NThaMDQxMjAwBgNVBAMTKU1pY3Jvc29mdCBBenVyZSBGZWRlcmF0ZWQgU1NPIENlcnRpZmljYXRlMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAv08WgEsT86mCgByv/G+c6OmazgVxcSo6HvOUD8iilfu+WisMsFfmitFLx/ZHKGiLwTzK4Q0/Ekv+tTmzLcIJx5EyY/hk6Md5KYGDx2bOzaRlFzVq3Sm2fAnyHZyB5qq7koHjYUaIW+sNLxt7Xz0isWVNcOElmC6RrP16x3lNpQVLaYyiDvPVwArfoLa8GrzeF9bbCWBLvKDKDP4tT7uP/GxgY8IgNRwbdnbtqpuvQlEBsZidjdlQckMP5bg0friicsHJ/gqJ2TV7rGbWk/BrUjrngSrW1CqLa/HEgmPw9WSJAMUKYjJ2zf2OXPUDqA6K1qs65SOPQEJULP3KIqQcIQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCf1X/RQdcoizdvjun4ezF0hI+tg0yfrYMsvtumsEodQXX+rjnZ1fN3pkOceo3wsiKVeTgCceLuXLUWGQ04HUot4SU2qYI5eY+0oLM978s8l4SIHrPkj8yTb1sFuVnkmG2UMc/77wW9mKBIR7QuaocVFCHyuyMwtt0KrEV5X4+zD6yaNGtpI8cWHcSgCpoEO8lt7kf65kFOannhn0PZ90sTDXgW7kPmBXwl8TRjJBvV/b3HLzl+ZgAVkosNeJmJifZYBds//mE5edz8GwGj9Q67XNxousk3k188nWe5tRREwAqo+gw5a8y6Quz1CrUGlp2g257MSXlyVkrzsdDB88aR</X509Certificate></X509Data></KeyInfo></KeyDescriptor><fed:TargetScopes><wsa:EndpointReference xmlns:wsa=\"http://www.w3.org/2005/08/addressing\"><wsa:Address>https://sts.windows.net/58f82789-6695-4f4a-abdb-357668d55cff/</wsa:Address></wsa:EndpointReference></fed:TargetScopes><fed:ApplicationServiceEndpoint><wsa:EndpointReference xmlns:wsa=\"http://www.w3.org/2005/08/addressing\"><wsa:Address>https://login.microsoftonline.com/58f82789-6695-4f4a-abdb-357668d55cff/wsfed</wsa:Address></wsa:EndpointReference></fed:ApplicationServiceEndpoint><fed:PassiveRequestorEndpoint><wsa:EndpointReference xmlns:wsa=\"http://www.w3.org/2005/08/addressing\"><wsa:Address>https://login.microsoftonline.com/58f82789-6695-4f4a-abdb-357668d55cff/wsfed</wsa:Address></wsa:EndpointReference></fed:PassiveRequestorEndpoint></RoleDescriptor><IDPSSODescriptor protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\"><KeyDescriptor use=\"signing\"><KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\"><X509Data><X509Certificate>MIIC8DCCAdigAwIBAgIQSP5mK0sN6IJGzefNl8rxQDANBgkqhkiG9w0BAQsFADA0MTIwMAYDVQQDEylNaWNyb3NvZnQgQXp1cmUgRmVkZXJhdGVkIFNTTyBDZXJ0aWZpY2F0ZTAeFw0yMjAxMTEwMjM5NTRaFw0yNTAxMTEwMjM5NThaMDQxMjAwBgNVBAMTKU1pY3Jvc29mdCBBenVyZSBGZWRlcmF0ZWQgU1NPIENlcnRpZmljYXRlMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAv08WgEsT86mCgByv/G+c6OmazgVxcSo6HvOUD8iilfu+WisMsFfmitFLx/ZHKGiLwTzK4Q0/Ekv+tTmzLcIJx5EyY/hk6Md5KYGDx2bOzaRlFzVq3Sm2fAnyHZyB5qq7koHjYUaIW+sNLxt7Xz0isWVNcOElmC6RrP16x3lNpQVLaYyiDvPVwArfoLa8GrzeF9bbCWBLvKDKDP4tT7uP/GxgY8IgNRwbdnbtqpuvQlEBsZidjdlQckMP5bg0friicsHJ/gqJ2TV7rGbWk/BrUjrngSrW1CqLa/HEgmPw9WSJAMUKYjJ2zf2OXPUDqA6K1qs65SOPQEJULP3KIqQcIQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCf1X/RQdcoizdvjun4ezF0hI+tg0yfrYMsvtumsEodQXX+rjnZ1fN3pkOceo3wsiKVeTgCceLuXLUWGQ04HUot4SU2qYI5eY+0oLM978s8l4SIHrPkj8yTb1sFuVnkmG2UMc/77wW9mKBIR7QuaocVFCHyuyMwtt0KrEV5X4+zD6yaNGtpI8cWHcSgCpoEO8lt7kf65kFOannhn0PZ90sTDXgW7kPmBXwl8TRjJBvV/b3HLzl+ZgAVkosNeJmJifZYBds//mE5edz8GwGj9Q67XNxousk3k188nWe5tRREwAqo+gw5a8y6Quz1CrUGlp2g257MSXlyVkrzsdDB88aR</X509Certificate></X509Data></KeyInfo></KeyDescriptor><SingleLogoutService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"https://login.microsoftonline.com/58f82789-6695-4f4a-abdb-357668d55cff/saml2\" /><SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"https://login.microsoftonline.com/58f82789-6695-4f4a-abdb-357668d55cff/saml2\" /><SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"https://login.microsoftonline.com/58f82789-6695-4f4a-abdb-357668d55cff/saml2\" />")

	metadata := &types.EntityDescriptor{}
	err = xml.Unmarshal(rawMetadata, metadata)
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
			// fmt.Println(xcert.Data)
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

	//assertionConsumerServiceURL := "/" + base + "/saml-configs/" + config.ID.Hex() + "/consume"
	sp := &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:      metadata.IDPSSODescriptor.SingleSignOnServices[0].Location,
		IdentityProviderIssuer:      metadata.EntityID,
		ServiceProviderIssuer:       config.ServiceProviderIssuer,       //"https://mentornity.com/saml/sp",
		AssertionConsumerServiceURL: config.AssertionConsumerServiceURL, // "http://localhost:1001/saml/consume", // https://mentornity.com/saml/consume
		SignAuthnRequests:           true,
		AudienceURI:                 config.AudienceURI, //"https://mentornity.com/saml",
		IDPCertificateStore:         &certStore,
		SPKeyStore:                  randomKeyStore,
	}
	return sp, err
}
