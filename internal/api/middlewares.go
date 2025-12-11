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

func CreateStack(xs ...MiddlewareFuncs) MiddlewareFuncs {
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

func (cfg *Middleware) LoggerMw(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			errRes := common.NewErrorResponse(common.ApiErrorMessage(), common.ApiError, nil)
			common.ResError(w, http.StatusInternalServerError, errRes, err)
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

func (cfg *Middleware) CheckReqBodyLengthMw(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, int64(cfg.MaxReqSize))
		body, err := io.ReadAll(r.Body)
		if err != nil {
			if err.Error() == "http: request body too large" {
				errRes := common.NewErrorResponse("Request body too large", common.InvalidRequestError, nil)
				common.ResError(w, http.StatusBadRequest, errRes, err)
				return
			}
			errRes := common.NewErrorResponse(common.ApiErrorMessage(), common.ApiError, nil)
			common.ResError(w, http.StatusInternalServerError, errRes, err)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}

func (cfg *Middleware) ValidateJsonMw(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if len(body) > 0 {
			if err != nil {
				errRes := common.NewErrorResponse(common.ApiErrorMessage(), common.ApiError, nil)
				common.ResError(w, http.StatusInternalServerError, errRes, err)
				return
			}

			if !json.Valid(body) {
				errRes := common.NewErrorResponse(common.InvalidRequestBodyMessage(), common.InvalidRequestError, nil)
				common.ResError(w, http.StatusBadRequest, errRes, nil)
				return
			}
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}

func (cfg *Middleware) RecoveryMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Caught panic: %v\nStack trace:\n%s", err, string(debug.Stack()))
				errRes := common.NewErrorResponse(common.ApiErrorMessage(), common.ApiError, nil)
				common.ResError(w, http.StatusInternalServerError, errRes, nil)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func (cfg *Middleware) CheckRouteAndMethodMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		if strings.Contains(wrapped.buf.String(), "404 page not found") {
			errRes := common.NewErrorResponse(
				common.RouteUnknownMessage(r.Method, r.URL.Path),
				common.InvalidRequestError,
				nil,
			)

			res := struct {
				Error *common.ErrResponse `json:"error"`
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
			errRes := common.NewErrorResponse(
				common.MethodNotAllowedMessage(r.Method, r.URL.Path),
				common.InvalidRequestError,
				nil,
			)
			res := struct {
				Error *common.ErrResponse `json:"error"`
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
