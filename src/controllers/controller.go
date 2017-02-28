package controllers

import (
	"github.com/gorilla/schema"
	"github.com/unrolled/render"
	"strings"
	"github.com/gorilla/mux"
	"io"
	"bytes"
	"encoding/json"
	"net/http"
	"config"
)

const (
	HeaderAccept = "Accept"
	HeaderContentType = "Content-Type"
	MediaTypeJson = "application/json.*"
	MediaTypeForm = "application/x-www-form-urlencoded.*"
)

var (
	Renderer *render.Render
	decoder *schema.Decoder
	renderer = render.New()
)

func init() {
	Renderer = render.New(render.Options{
		Directory:     config.ActiveThemeTemplatesDir(),
		Extensions:    []string{".html"},
		IsDevelopment: config.IsDevelopment(),
	})

	decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
}

type errResponse struct {
	Code    int
	Reason  string
	Message string
}

type validationErrResponse struct {
	Code   int
	Reason string
	Errors []struct {
		Message   string
		FieldName string
	}
}

func NewErrResponse(code int, err string) errResponse {
	return errResponse{
		Code:    code,
		Reason:  http.StatusText(code),
		Message: err,
	}
}

func BadRequestErr(err string) errResponse {
	return NewErrResponse(http.StatusBadRequest, err)
}

func ServerErr(err string) errResponse {
	return NewErrResponse(http.StatusInternalServerError, err)
}

func NotFoundErr(err string) errResponse {
	return NewErrResponse(http.StatusNotFound, err)
}

func (ver *validationErrResponse) AddErr(message string) {
	if err := strings.Split(message, ":"); len(err) == 2 {
		result := struct {
			Message   string
			FieldName string
		}{
			Message:   err[1],
			FieldName: err[0],
		}

		ver.Errors = append(ver.Errors, result)
	}
}

func ValidationErr(errs string) validationErrResponse {
	result := validationErrResponse{
		Code:   422,
		Reason: http.StatusText(422),
	}

	if errors := strings.Split(errs, ";"); len(errors) > 0 {
		for _, message := range errors {
			result.AddErr(message)
		}
	} else {
		result.AddErr(errs)
	}

	return result
}

type Api interface {
	Router() *mux.Router
	RegisterEndpoints()
}

func ParseJson(r io.ReadCloser, dst interface{}) error {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	defer r.Close()
	if err := json.Unmarshal(buf.Bytes(), dst); err != nil {
		return err
	}

	return nil
}

func ParseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := decoder.Decode(dst, r.Form); err != nil {
		return err
	}

	return nil
}