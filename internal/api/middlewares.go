package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
}

func (w *wrappedWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (cfg *ApiConfig) LoggerMw(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			resError(w, http.StatusInternalServerError, ApiErrorMessage(), ApiError, nil, err)
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

func (cfg *ApiConfig) CheckReqBodyLengthMw(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, int64(cfg.MaxReqSize))
		body, err := io.ReadAll(r.Body)
		if err != nil {
			if err.Error() == "http: request body too large" {
				resError(w, http.StatusBadRequest, "Request body too large", InvalidRequestError, nil, err)
				return
			}
			resError(w, http.StatusInternalServerError, ApiErrorMessage(), ApiError, nil, err)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) ValidateJsonMw(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if len(body) > 0 {
			if err != nil {
				if err.Error() == "EOF" {
					resError(w, http.StatusBadRequest, InvalidRequestBodyMessage(), InvalidRequestError, nil, nil)
					return
				}
				resError(w, http.StatusInternalServerError, ApiErrorMessage(), ApiError, nil, err)
				return
			}

			if !json.Valid(body) {
				resError(w, http.StatusBadRequest, InvalidRequestBodyMessage(), InvalidRequestError, nil, nil)
				return
			}
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) RecoveryMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg := "Caught panic: %v, Stack trace: %s"
				log.Printf(msg, err, string(debug.Stack()))
				resError(w, http.StatusInternalServerError, ApiErrorMessage(), ApiError, nil, nil)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func (cfg *ApiConfig) CheckRouteAndMethodMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		if strings.Contains(wrapped.buf.String(), "404 page not found") {
			errRes := struct {
				Error ErrResponse `json:"error"`
			}{
				Error: ErrResponse{
					Message: RouteUnknownMessage(r.Method, r.URL.Path),
					Type:    InvalidRequestError,
				},
			}

			data, err := json.MarshalIndent(errRes, "", "  ")
			if err != nil {
				log.Printf("error marshaling JSON: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			wrapped.buf.Reset()
			wrapped.buf.Write(data)
		}

		if wrapped.statusCode == http.StatusMethodNotAllowed {
			errRes := struct {
				Error ErrResponse `json:"error"`
			}{
				Error: ErrResponse{
					Message: MethodNotAllowedMessage(r.Method, r.URL.Path),
					Type:    InvalidRequestError,
				},
			}

			data, err := json.MarshalIndent(errRes, "", "  ")
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
