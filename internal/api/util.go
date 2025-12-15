package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

func GetIdFromPath(r *http.Request) (uuid.UUID, string) {
	pathValue := r.PathValue("id")
	id, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.UUID{}, InvalidIdMessage(pathValue, "account")
	}
	return id, ""
}

func JsonDecode(r *http.Request, v any, w http.ResponseWriter) *ErrorResponse {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			errRes := formErrResponse(ute)
			return errRes
		}
		return &ErrorResponse{
			Message: InvalidRequestBodyMessage(err),
			Type:    InvalidRequestError,
		}
	}
	return nil
}

func formErrResponse(e *json.UnmarshalTypeError) *ErrorResponse {
	msg := fmt.Sprintf("Type '%s' is not assignable to '%s' parameter.", e.Value, e.Field)
	return &ErrorResponse{
		Code:    ParameterInvalid,
		Message: msg,
		Param:   e.Field,
		Type:    InvalidRequestError,
	}
}

func ToSnakeCase(input string) string {
	var b strings.Builder
	var prev rune
	for i, v := range input {
		if unicode.IsLower(v) {
			b.WriteRune(v)
		} else {
			if i > 0 && (unicode.IsLower(prev) ||
				unicode.IsLower(nextRune(input[i+utf8.RuneLen(v):]))) {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(v))
		}
		prev = v
	}
	return b.String()
}

func nextRune(s string) rune { r, _ := utf8.DecodeRuneInString(s); return r }
