package kimlikruntime

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	core_model "github.com/devingen/api-core/model"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func AssertBaseAndBody(ctx context.Context, req dvnruntime.Request, bodyValue interface{}) (string, error) {
	// retrieve base from path
	base := req.PathParameters["base"]
	if base == "" {
		return "", core_model.NewError(http.StatusBadRequest, "base-missing")
	}

	// assert body is present
	if req.Body == "" {
		return "", core_model.NewError(http.StatusBadRequest, "body-missing")
	}

	// parse body
	err := dvnruntime.ParseBody(req.Body, &bodyValue)
	if err != nil {
		return base, err
	}

	// validate body
	if validate == nil {
		validate = validator.New()
	}
	err = validate.Struct(bodyValue)

	// return proper validation error
	if err != nil {
		switch castedError := err.(type) {
		case validator.ValidationErrors:
			return base, core_model.NewError(http.StatusBadRequest, castedError.Error())
		default:
			return base, err
		}
	}

	return base, nil
}
