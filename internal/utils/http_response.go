package utils

import (
	"net/http"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
)

func ResponseWasNotFound(err error) bool {
	// If API returns 404, resource was not found
	if ResponseWasStatusCode(err, http.StatusNotFound) {
		return true
	}

	// Some APIs return 400 BadRequest with the VS800075 error message if
	// DevOps Project doesn't exist. If parent project doesn't exist, all
	// child resources are considered "doesn't exist".
	if ResponseWasStatusCode(err, http.StatusBadRequest) {
		return ResponseContainsStatusMessage(err, "VS800075")
	}
	return false
}

// ResponseWasStatusCode was used for check if error status code was specific http status code
func ResponseWasStatusCode(err error, statusCode int) bool {
	if err == nil {
		return false
	}
	if wrapperErr, ok := err.(azuredevops.WrappedError); ok {
		if wrapperErr.StatusCode != nil && *wrapperErr.StatusCode == statusCode {
			return true
		}
	}
	return false
}

// ResponseContainsStatusMessage is used for check if error message contains specific message
func ResponseContainsStatusMessage(err error, statusMessage string) bool {
	if err == nil {
		return false
	}
	if wrapperErr, ok := err.(azuredevops.WrappedError); ok {
		if wrapperErr.Message == nil {
			return false
		}
		return strings.Contains(*wrapperErr.Message, statusMessage)
	}
	return false
}
