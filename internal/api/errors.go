package api

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ErrorType string

const (
	ApiError            ErrorType = "api_error"
	InvalidRequestError ErrorType = "invalid_request_error"
)

type ErrorCode string

const (
	ApiKeyExpired        ErrorCode = "api_key_expired"
	CountryUnknown       ErrorCode = "country_unknown"
	CurrencyUnknown      ErrorCode = "currency_unknown"
	ParameterInvalid     ErrorCode = "parameter_invalid"
	ParameterMissing     ErrorCode = "parameter_missing"
	ResourceNotFound     ErrorCode = "resource_not_found"
	StringLengthExceeded ErrorCode = "string_length_exceeded"
	StringLengthNotMet   ErrorCode = "string_length_not_met"
)

func ApiErrorMessage() string {
	return "Something went wrong."
}

func CountryUnknownMessage(input string) string {
	return fmt.Sprintf("Country '%s' is unknown. "+
		"Try using a 2-character alphanumeric country code instead, such as 'US', 'IE', or 'GB'. "+
		"A complete list of official country codes is available at https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#Officially_assigned_code_elements",
		input,
	)
}

func CurrencyUnknownMessage(input string) string {
	return fmt.Sprintf("Currency '%s' is unknown. "+
		"Try using a 3-character alphanumeric currency code instead, such as 'USD', 'EUR', or 'GBP'. "+
		"A complete list of official currency codes is available at https://en.wikipedia.org/wiki/ISO_4217#Active_codes_(list_one)",
		input,
	)
}

func ExclusiveParamsMessage(a, b string) string {
	return fmt.Sprintf("Received both '%s' and '%s' parameters. Pass one at a time.", a, b)
}

func InvalidIdMessage(value, resource string) string {
	return fmt.Sprintf("Provided value '%s' is not a valid %s id.", value, resource)
}

func InvalidRequestBodyMessage(e error) string {
	return fmt.Sprintf("Invalid request body: %s.", e.Error())
}

func MethodNotAllowedMessage(method, path string) string {
	return fmt.Sprintf("Method '%s' not allowed on %s.", method, path)
}

func NotFoundMessage(id uuid.UUID, resource string) string {
	return fmt.Sprintf("%s with id '%s' not found.", cases.Title(language.English).String(resource), id)
}

func ParameterInvalidMessage(param string) string {
	return fmt.Sprintf("Parameter invalid: '%s'.", param)
}

func ParameterMissingMessage(param string) string {
	return fmt.Sprintf("Missing required param: '%s'.", param)
}

func RequestTooLargeMessage() string {
	return "Request body too large."
}

func RouteUnknownMessage(method, path string) string {
	return fmt.Sprintf("Request to unknown route (%s: %s).", method, path)
}

func StringLengthExceededMessage(param, v string) string {
	return fmt.Sprintf("The length of the '%s' parameter cannot be greater than %s characters.", param, v)
}

func StringLengthNotMetMessage(param, v string) string {
	return fmt.Sprintf("The length of the '%s' parameter must be at least %s characters.", param, v)
}

func ValueNotGtMessage(param, v string) string {
	return fmt.Sprintf("Value of '%s' param must be greater than '%s'.", param, v)
}

func ValueNotGteMessage(param, v string) string {
	return fmt.Sprintf("Value of '%s' param must be greater than or equal to '%s'.", param, v)
}

func ValueNotLtMessage(param, v string) string {
	return fmt.Sprintf("Value of '%s' param must be less than '%s'.", param, v)
}

func ValueNotLteMessage(param, v string) string {
	return fmt.Sprintf("Value of '%s' param must be less than or equal to '%s'.", param, v)
}
