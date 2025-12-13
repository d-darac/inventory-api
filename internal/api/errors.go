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
	ApiKeyExpired    ErrorCode = "api_key_expired"
	CountryUnknown   ErrorCode = "country_unknown"
	CurrencyUnknown  ErrorCode = "currency_unknown"
	ParameterInvalid ErrorCode = "parameter_invalid"
	ResourceNotFound ErrorCode = "resource_not_found"
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

func InvalidIdMessage(value, resource string) string {
	return fmt.Sprintf("Provided value '%s' is not a valid %s id.", value, resource)
}

func InvalidRequestBodyMessage() string {
	return "Invalid request body. Make sure that the body is in format application/json."
}

func MethodNotAllowedMessage(method, path string) string {
	return fmt.Sprintf("Method '%s' not allowed on %s.", method, path)
}

func NotFoundMessage(id uuid.UUID, resource string) string {
	return fmt.Sprintf("%s with id '%s' not found.", cases.Title(language.English).String(resource), id)
}

func RequestTooLargeMessage() string {
	return "Request body too large."
}

func RouteUnknownMessage(method, path string) string {
	return fmt.Sprintf("Request to unknown route (%s: %s).", method, path)
}
