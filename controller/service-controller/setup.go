package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"

	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
)

func (c ServiceController) Setup(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.SetupRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	////Generate RSA keys
	//privKey, err := GenerateRSAKeys(2048)
	//if err != nil {
	//	fmt.Println("Error generating RSA keys:", err)
	//	return
	//}
	//
	//// Convert the private and public keys to PEM strings
	//privKeyString := ConvertPrivateKeyToPEM(privKey)
	//pubKeyString := ConvertPublicKeyToPEM(&privKey.PublicKey)
	//
	//// Print the PEM strings (for demonstration purposes)
	//fmt.Println("Private Key (PEM):\n", privKeyString)
	//fmt.Println("Public Key (PEM):\n", pubKeyString)

	item, err := c.DataService.CreateTenantInfo(ctx, base, &model.TenantInfo{
		Name: &body.TenantName,
	})
	if err != nil {
		return nil, err
	}

	integrationSettings, err := c.DataService.CreateIntegrationSettings(ctx, base, &model.IntegrationSettings{})
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: map[string]interface{}{
			"tenantInfo":          item,
			"integrationSettings": integrationSettings,
		},
	}, nil
}

//// Generate RSA keys
//func GenerateRSAKeys(bits int) (*rsa.PrivateKey, error) {
//	privKey, err := rsa.GenerateKey(rand.Reader, bits)
//	if err != nil {
//		return nil, err
//	}
//	return privKey, nil
//}
//
//// Convert private key to PEM string
//func ConvertPrivateKeyToPEM(privKey *rsa.PrivateKey) string {
//	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
//	privBlock := &pem.Block{
//		Type:  "RSA PRIVATE KEY",
//		Bytes: privBytes,
//	}
//	return string(pem.EncodeToMemory(privBlock))
//}
//
//// Convert public key to PEM string
//func ConvertPublicKeyToPEM(pubKey *rsa.PublicKey) string {
//	pubBytes, err := x509.MarshalPKIXPublicKey(pubKey)
//	if err != nil {
//		panic("Error marshalling public key: " + err.Error())
//	}
//	pubBlock := &pem.Block{
//		Type:  "PUBLIC KEY",
//		Bytes: pubBytes,
//	}
//	return string(pem.EncodeToMemory(pubBlock))
//}
