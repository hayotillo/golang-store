package store

import (
	"errors"
	"fmt"
)

var ErrUserNotAuthenticated = errors.New("user not authenticated")
var ErrLoginOrPasswordIncorrect = errors.New("login or password incorrect")
var ErrRequiredDataNotFount = errors.New("required data not fount")
var ErrRecordDuplicate = errors.New("record duplicate")
var ErrProductDoNotBeSync = errors.New("product do not be sync")
var ErrUserPermissionDenied = errors.New("user permission denied")
var ErrFileSystemPermissionDenied = errors.New("file system permission denied")
var ErrImageResolutionRequired = errors.New("image resolution required")
var ErrImageFormatNotSupported = errors.New("image format not supported")
var ErrRecordNotFound = errors.New("record not found")
var ErrFileNotFound = errors.New("file not found")
var ErrRecordAlreadyExists = errors.New("record already exists")
var ErrRecordDoNotBeDelete = errors.New("record do not be delete")

func ErrManyDataIncorrect(field string, value string) error {
	return errors.New(fmt.Sprintf("field: %s have incorrect value: %s", field, value))
}
