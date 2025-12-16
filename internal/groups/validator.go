package groups

import (
	"errors"

	"github.com/d-darac/inventory-api/internal/api"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func validateCreateParams(params *createParams) *api.ErrorListResponse {
	validate = validator.New(validator.WithRequiredStructEnabled())

	return validateStruct(params)
}

func validateStruct(params *createParams) *api.ErrorListResponse {
	err := validate.Struct(params)
	if err != nil {
		var errListRes api.ErrorListResponse
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				errListRes.Errors = append(errListRes.Errors, fieldErrToErrRes(e))
			}
		}
		return &errListRes
	}
	return nil
}

func fieldErrToErrRes(e validator.FieldError) api.ErrorResponse {
	param := api.ToSnakeCase(e.StructField())
	errRes := api.ErrorResponse{
		Message: api.ParameterInvalidMessage(param),
		Type:    api.InvalidRequestError,
		Param:   api.ToSnakeCase(param),
	}

	switch e.ActualTag() {
	case "required":
		errRes.Code = api.ParameterMissing
		errRes.Message = api.ParameterMissingMessage(errRes.Param)
	case "lte":
		errRes.Code = api.CharacterLimitExceeded
		errRes.Message = api.CharacterLimitExceededMessage(errRes.Param, e.Param())
	}
	return errRes
}
