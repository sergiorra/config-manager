package config_manager

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	ErrInvalidConfig struct{}

	ErrEmptyPrefix struct{}

	ErrMandatoryField struct {
		fieldName string
	}

	ErrTypeNotSupported struct {
		fieldKind string
	}

	ErrMultiLevelNestedStruct struct{}

	ErrBindFlags struct {
		inner error
	}

	ErrBindEnvVar struct {
		inner error
	}

	ErrNotFound struct {
		fileName string
	}

	ErrReadFile struct {
		inner error
	}

	ErrUnmarshalConfig struct {
		inner error
	}
)

// NewErrInvalidConfig returns a new invalid configuration error instance
func NewErrInvalidConfig() *ErrInvalidConfig {
	return &ErrInvalidConfig{}
}

func (err *ErrInvalidConfig) Error() string {
	return "invalid configuration struct"
}

// NewErrEmptyPrefix returns a new empty prefix error instance
func NewErrEmptyPrefix() *ErrEmptyPrefix {
	return &ErrEmptyPrefix{}
}

func (err *ErrEmptyPrefix) Error() string {
	return "environment variables prefix can not be an empty string"
}

// NewErrMandatoryField returns a new mandatory field error instance
func NewErrMandatoryField(fieldName string) *ErrMandatoryField {
	return &ErrMandatoryField{fieldName}
}

func (err *ErrMandatoryField) Error() string {
	return fmt.Sprintf("not value set by mandatory host field %s", err.fieldName)
}

// NewErrTypeNotSupported returns a new type not supported error instance
func NewErrTypeNotSupported(fieldKind string) *ErrTypeNotSupported {
	return &ErrTypeNotSupported{fieldKind}
}

func (err *ErrTypeNotSupported) Error() string {
	return fmt.Sprintf("config type %s not supported", err.fieldKind)
}

// NewErrMultiLevelNestedStruct returns a new multi level nested struct error instance
func NewErrMultiLevelNestedStruct() *ErrMultiLevelNestedStruct {
	return &ErrMultiLevelNestedStruct{}
}

func (err *ErrMultiLevelNestedStruct) Error() string {
	return "multiple level nested struct not supported"
}

// NewErrBindFlags returns a new bind flags error instance
func NewErrBindFlags(inner error, fieldName string) *ErrBindFlags {
	return &ErrBindFlags{inner: errors.Wrapf(inner, "flag not bind by field %s", fieldName)}
}

func (err *ErrBindFlags) Error() string {
	return err.inner.Error()
}

// NewErrBindEnvVar returns a new bind environment variable error instance
func NewErrBindEnvVar(inner error, fieldName string) *ErrBindEnvVar {
	return &ErrBindEnvVar{inner: errors.Wrapf(inner, "environment variable not bind by field %s", fieldName)}
}

func (err *ErrBindEnvVar) Error() string {
	return err.inner.Error()
}

// NewErrNotFound returns a new not found error instance
func NewErrNotFound(name string) *ErrNotFound {
	return &ErrNotFound{name}
}

func (err *ErrNotFound) Error() string {
	return fmt.Sprintf("file not found %s", err.fileName)
}

// NewErrReadFile returns a new read file error instance
func NewErrReadFile(inner error, fileName string) *ErrReadFile {
	return &ErrReadFile{inner: errors.Wrapf(inner, "read in file %s", fileName)}
}

func (err *ErrReadFile) Error() string {
	return err.inner.Error()
}

// NewErrUnmarshalConfig returns a new unmarshal configuration error instance
func NewErrUnmarshalConfig(inner error, filename string) *ErrUnmarshalConfig {
	return &ErrUnmarshalConfig{inner: errors.Wrapf(inner, "unmarshal config file %s into struct", filename)}
}

func (err *ErrUnmarshalConfig) Error() string {
	return err.inner.Error()
}
