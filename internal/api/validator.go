package api

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		Validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *Validator) ValidateRequestParams(params interface{}) *ErrorListResponse {
	err := v.Validate.Struct(params)
	if err != nil {
		var errListRes ErrorListResponse
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

func fieldErrToErrRes(e validator.FieldError) ErrorResponse {
	param := e.StructField()
	paramType := e.Type().String()
	errRes := ErrorResponse{
		Code:    ParameterInvalid,
		Message: ParameterInvalidMessage(param),
		Type:    InvalidRequestError,
		Param:   ToSnakeCase(param),
	}

	switch e.ActualTag() {
	case "excluded_with":
		errRes.Message = ExclusiveParamsMessage(errRes.Param, ToSnakeCase(e.Param()))
		errRes.Param = ""
	case "gt":
		if paramType == "string" {
			errRes.Code = StringLengthExceeded
			errRes.Message = StringLengthNotMetMessage(errRes.Param, e.Param())
		}
		errRes.Message = ValueNotGtMessage(errRes.Param, e.Param())
	case "gte":
		if paramType == "string" {
			errRes.Code = StringLengthExceeded
			errRes.Message = StringLengthNotMetMessage(errRes.Param, e.Param())
		}
		errRes.Message = ValueNotGteMessage(errRes.Param, e.Param())
	case "lt":
		if paramType == "string" {
			errRes.Code = StringLengthExceeded
			errRes.Message = StringLengthExceededMessage(errRes.Param, e.Param())
		}
		errRes.Message = ValueNotLtMessage(errRes.Param, e.Param())
	case "lte":
		if paramType == "string" {
			errRes.Code = StringLengthExceeded
			errRes.Message = StringLengthExceededMessage(errRes.Param, e.Param())
		}
		errRes.Message = ValueNotLteMessage(errRes.Param, e.Param())
	case "required":
		errRes.Code = ParameterMissing
		errRes.Message = ParameterMissingMessage(errRes.Param)
	default:
		errRes.Message = ParameterInvalidMessage(errRes.Param)
	}
	return errRes
}
