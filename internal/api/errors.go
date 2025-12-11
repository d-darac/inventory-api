package api

import "fmt"

type ErrorType string

const (
	ApiError            ErrorType = "api_error"
	InvalidRequestError ErrorType = "invalid_request_error"
)

type ErrorCode string

const (
	ApiKeyExpired   ErrorCode = "api_key_expired"
	CountryUnknown  ErrorCode = "country_unknown"
	CurrencyUnknown ErrorCode = "currency_unknown"
)

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

func InvalidRequestBodyMessage() string {
	return "Invalid request body. Make sure that the body is in format application/json."
}

func ApiErrorMessage() string {
	return "Something went wrong."
}

func RouteUnknownMessage(method, path string) string {
	return fmt.Sprintf("Request to unknown route (%s: %s).", method, path)
}

func MethodNotAllowedMessage(method, path string) string {
	return fmt.Sprintf("Method '%s' not allowed on %s.", method, path)
}
