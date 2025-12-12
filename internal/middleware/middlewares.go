package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/d-darac/inventory-api/internal/api"
	"github.com/d-darac/inventory-api/internal/common"
)

type Middleware struct {
	MaxReqSize int
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
			api.ResError(w, http.StatusInternalServerError, common.ApiErrorMessage(), common.ApiError, nil, err)
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
				api.ResError(w, http.StatusBadRequest, common.RequestTooLargeMessage(), common.InvalidRequestError, nil, err)
				return
			}
			api.ResError(w, http.StatusInternalServerError, common.ApiErrorMessage(), common.ApiError, nil, err)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}

func (mw *Middleware) ValidateJsonMw(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if len(body) > 0 {
			if err != nil {
				api.ResError(w, http.StatusInternalServerError, common.ApiErrorMessage(), common.ApiError, nil, err)
				return
			}

			if !json.Valid(body) {
				api.ResError(w, http.StatusBadRequest, common.InvalidRequestBodyMessage(), common.InvalidRequestError, nil, nil)
				return
			}
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
				api.ResError(w, http.StatusInternalServerError, common.ApiErrorMessage(), common.ApiError, nil, nil)
				return
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func (mw *Middleware) CheckRouteAndMethodMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		if strings.Contains(wrapped.buf.String(), "404 page not found") {
			errRes := &api.ErrResponse{
				Message: common.RouteUnknownMessage(r.Method, r.URL.Path),
				Type:    common.InvalidRequestError,
			}

			res := struct {
				Error *api.ErrResponse `json:"error"`
			}{
				Error: errRes,
			}

			data, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				log.Printf("error marshaling JSON: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			wrapped.buf.Reset()
			wrapped.buf.Write(data)
		}

		if wrapped.statusCode == http.StatusMethodNotAllowed {
			errRes := &api.ErrResponse{
				Message: common.MethodNotAllowedMessage(r.Method, r.URL.Path),
				Type:    common.InvalidRequestError,
			}

			res := struct {
				Error *api.ErrResponse `json:"error"`
			}{
				Error: errRes,
			}

			data, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				log.Printf("error marshaling JSON: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			wrapped.buf.Reset()
			wrapped.buf.Write(data)
		}

		if _, err := io.Copy(w, wrapped.buf); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
	}
}
