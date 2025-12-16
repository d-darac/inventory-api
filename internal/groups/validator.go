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

func validateListParams(params *listParams) *api.ErrorListResponse {
	validate = validator.New(validator.WithRequiredStructEnabled())

	return validateStruct(params)
}

func validateStruct(params interface{}) *api.ErrorListResponse {
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
	param := e.StructField()
	paramType := e.Type().String()
	errRes := api.ErrorResponse{
		Code:    api.ParameterInvalid,
		Message: api.ParameterInvalidMessage(param),
		Type:    api.InvalidRequestError,
		Param:   api.ToSnakeCase(param),
	}

	switch e.ActualTag() {
	case "excluded_with":
		errRes.Message = api.ExclusiveParamsMessage(errRes.Param, api.ToSnakeCase(e.Param()))
		errRes.Param = ""
	case "gt":
		if paramType == "string" {
			errRes.Code = api.StringLengthExceeded
			errRes.Message = api.StringLengthNotMetMessage(errRes.Param, e.Param())
		}
		errRes.Message = api.ValueNotGtMessage(errRes.Param, e.Param())
	case "gte":
		if paramType == "string" {
			errRes.Code = api.StringLengthExceeded
			errRes.Message = api.StringLengthNotMetMessage(errRes.Param, e.Param())
		}
		errRes.Message = api.ValueNotGteMessage(errRes.Param, e.Param())
	case "lt":
		if paramType == "string" {
			errRes.Code = api.StringLengthExceeded
			errRes.Message = api.StringLengthExceededMessage(errRes.Param, e.Param())
		}
		errRes.Message = api.ValueNotLtMessage(errRes.Param, e.Param())
	case "lte":
		if paramType == "string" {
			errRes.Code = api.StringLengthExceeded
			errRes.Message = api.StringLengthExceededMessage(errRes.Param, e.Param())
		}
		errRes.Message = api.ValueNotLteMessage(errRes.Param, e.Param())
	case "required":
		errRes.Code = api.ParameterMissing
		errRes.Message = api.ParameterMissingMessage(errRes.Param)
	default:
		errRes.Message = api.ParameterInvalidMessage(errRes.Param)
	}
	return errRes
}
