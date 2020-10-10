package kimlikruntime

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	core_model "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/token-service"
	"net/http"
)

func AssertAuthentication(ctx context.Context, req dvnruntime.Request) (string, *token_service.TokenPayload, error) {
	tokenPayload := ctx.Value(token_service.ContextKeyTokenPayload).(*token_service.TokenPayload)
	if tokenPayload == nil {
		return "", nil, core_model.NewStatusError(http.StatusUnauthorized)
	}

	base := req.PathParameters["base"]
	if base == "" {
		return "", tokenPayload, core_model.NewError(http.StatusBadRequest, "base-missing")
	}
	return base, tokenPayload, nil
}

//func AssertUser(ctx context.Context, req dvnruntime.Request, controller service_controller.ServiceController) (string, *model.User, error) {
//	tokenPayload := ctx.Value(token_service.ContextKeyTokenPayload).(*token_service.TokenPayload)
//	if tokenPayload == nil {
//		return "", nil, core_model.NewStatusError(http.StatusUnauthorized)
//	}
//
//	base := req.PathParameters["base"]
//	if base == "" {
//		return "", nil, core_model.NewError(http.StatusBadRequest, "base-missing")
//	}
//
//	user, err := controller.Service.FindUserUserWithId(base, tokenPayload.UserId)
//	if user == nil {
//		return "", nil, core_model.NewError(http.StatusNotFound, "token-owner-not-found")
//	}
//
//	return base, user, err
//}
