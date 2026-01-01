package middleware

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"slices"
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/auth"
	"github.com/d-darac/inventory-assets/database"
)

type ctxKey string

var AuthAccountID ctxKey = "middleware.auth.accountID"

type Middleware struct {
	MaxReqSize int
	Db         *database.Queries
	Auth       struct {
		MasterKey string
		Iv        string
	}
}

type MiddlewareFuncs func(http.HandlerFunc) http.HandlerFunc

type wrappedWriter struct {
	http.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
}

func (mw *Middleware) CreateStack(xs ...MiddlewareFuncs) MiddlewareFuncs {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func (w *wrappedWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (mw *Middleware) LoggerMw(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			api.ResError(w, api.ApiErrorMessage())
			return
		}

		// TODO implement req body db store
		// requestHeaders := r.Header
		// requestBody := string(body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))

		// TODO implement res body db store
		// responseBody := wrapped.buf.String()
		if _, err := io.Copy(w, wrapped.buf); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
	})
}

func (mw *Middleware) CheckReqBodyLengthMw(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, int64(mw.MaxReqSize))
		body, err := io.ReadAll(r.Body)
		if err != nil {
			if err.Error() == "http: request body too large" {
				api.ResError(w, api.RequestTooLargeMessage())
				return
			}
			api.ResError(w, api.ApiErrorMessage())
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}

func (mw *Middleware) RecoveryMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Caught panic: %v\nStack trace:\n%s", err, string(debug.Stack()))
				api.ResError(w, api.ApiErrorMessage())
				return
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func (mw *Middleware) CheckRouteAndMethodMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathsMethods := map[string][]string{
			`^\/v1\/groups$`:                   {"GET", "POST"},
			`^\/v1\/groups\/[^\/]+$`:           {"DELETE", "GET", "PUT"},
			`^\/v1\/inventories$`:              {"GET", "POST"},
			`^\/v1\/inventories\/[^\/]+$`:      {"DELETE", "GET", "PUT"},
			`^\/v1\/items$`:                    {"GET", "POST"},
			`^\/v1\/items\/[^\/]+$`:            {"DELETE", "GET", "PUT"},
			`^\/v1\/item_identifiers$`:         {"GET", "POST"},
			`^\/v1\/item_identifiers\/[^\/]+$`: {"DELETE", "GET", "PUT"},
		}
		if err := validateRoute(r.Method, r.URL.Path, pathsMethods); err != nil {
			api.ResError(w, err)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (mw *Middleware) ApiKeyAuthMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKeyString, err := auth.GetApiKey(r.Header)
		if err != nil {
			api.ResError(w, err)
			return
		}

		secret, err := auth.EncryptApiKeySecret(apiKeyString, mw.Auth.MasterKey, mw.Auth.Iv)
		if err != nil {
			api.ResError(w, err)
			return
		}

		apiKey, err := mw.Db.GetApiKeyAccAndExp(context.Background(), secret)
		if err != nil {
			api.ResError(w, &api.AppError{
				Message: "Invalid api key.",
				Status:  http.StatusUnauthorized,
				Type:    api.InvalidRequestError,
			})
			return
		}

		if apiKey.ExpiresAt.Valid && time.Now().After(apiKey.ExpiresAt.Time) {
			api.ResError(w, &api.AppError{
				Message: "Invalid api key.",
				Status:  http.StatusUnauthorized,
				Type:    api.InvalidRequestError,
			})
			return
		}

		ctx := context.WithValue(r.Context(), AuthAccountID, apiKey.Account)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	}
}

func validateRoute(reqMethod, reqPath string, pathsMethods map[string][]string) error {
	methods := []string{}
	for kPath, vMethods := range pathsMethods {
		r := regexp.MustCompile(kPath)
		matched := r.MatchString(reqPath)
		if matched {
			methods = vMethods
			break
		}
	}
	if len(methods) == 0 {
		return api.RouteUnknownMessage(reqMethod, reqPath)
	}

	if !slices.Contains(methods, reqMethod) {
		return api.MethodNotAllowedMessage(reqMethod, reqPath)
	}
	return nil
}
