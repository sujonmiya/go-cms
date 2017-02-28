package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Err struct {
	Message    string
	FieldName  string
	FieldValue string
}

type Errs []Err

func (e *Errs) Add(message, fieldName, fieldValue string) {
	err := Err{
		Message:    message,
		FieldName:  fieldName,
		FieldValue: fieldValue,
	}
	*e = append(*e, err)
}

func (e *Errs) AddInvalidError(fieldName, fieldValue string) {
	err := Err{
		Message:    "is invalid",
		FieldName:  fieldName,
		FieldValue: fieldValue,
	}
	*e = append(*e, err)
}

func (e *Errs) AddEmptyError(fieldName string) {
	*e = append(*e, Err{Message: "can not be empty", FieldName: fieldName})
}

func (e *Errs) AddMaxLengthError(max uint16, fieldName, fieldValue string) {
	err := Err{
		Message:    fmt.Sprintf("must be less than %d characters", max),
		FieldName:  fieldName,
		FieldValue: strconv.Itoa(len(fieldValue)),
	}
	*e = append(*e, err)
}

func (e *Errs) AddLengthBetweenError(min, max uint16, fieldName, fieldValue string) {
	err := Err{
		Message:    fmt.Sprintf("must be between %d to %d characters", min, max),
		FieldName:  fieldName,
		FieldValue: strconv.Itoa(len(fieldValue)),
	}
	*e = append(*e, err)
}

func (e *Errs) AddNotAllowedError(fieldName, fieldValue, pattern string) {
	err := Err{
		Message:    fmt.Sprintf("must contain allowed characters:%s", strings.TrimPrefix(strings.TrimSuffix(pattern, "]+$"), "^[")),
		FieldName:  fieldName,
		FieldValue: fieldValue,
	}
	*e = append(*e, err)
}

func (e *Errs) AddInvalidIdError(fieldName string) {
	*e = append(*e, Err{Message: "must be a valid id", FieldName: fieldName})
}

func ValidationErr(err *Errs) ValidationErrResponse {
	return ValidationErrResponse{
		Code:   422,
		Reason: "Validation Error",
		Errors: err,
	}
}

type ValidationErrResponse struct {
	Code   int
	Reason string
	Errors *Errs
}
