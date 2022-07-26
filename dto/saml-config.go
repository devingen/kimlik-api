package dto

import "github.com/devingen/kimlik-api/model"

type CreateSAMLConfigRequest = model.SAMLConfig
type UpdateSAMLConfigRequest = model.SAMLConfig

type BuildSAMLAuthURLResponse struct {
	AuthURL *string `json:"authURL"`
}

type ConsumeSAMLAuthResponseRequest struct {
	SAMLResponse *string `json:"samlResponse" validate:"required"`
}
