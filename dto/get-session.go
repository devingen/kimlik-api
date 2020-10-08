package dto

import "github.com/devingen/kimlik-api/model"

type GetSessionResponse struct {
	User *model.User `json:"user"`
}
