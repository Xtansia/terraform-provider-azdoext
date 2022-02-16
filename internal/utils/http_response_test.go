package utils

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/stretchr/testify/require"
)

func azdoErr(sc *int, m *string) azuredevops.WrappedError {
	return azuredevops.WrappedError{
		StatusCode: sc,
		Message:    m,
	}
}

func TestResponseWasNotFound(t *testing.T) {
	test := func(input error, expected bool) func(*testing.T) {
		return func(t *testing.T) {
			res := ResponseWasNotFound(input)
			require.Equal(t, expected, res)
		}
	}

	t.Run("nil_error", test(nil, false))
	t.Run("non_azdo_error", test(fmt.Errorf("not an AZDO error"), false))
	t.Run("azdo_error_nil_status_and_message", test(azdoErr(nil, nil), false))
	t.Run("azdo_error_not_found_status", test(azdoErr(NewInt(http.StatusNotFound), nil), true))
	t.Run("azdo_error_bad_request_status_nil_message", test(azdoErr(NewInt(http.StatusBadRequest), nil), false))
	t.Run(
		"azdo_error_bad_request_status_irrelevant_message",
		test(azdoErr(NewInt(http.StatusBadRequest), NewString("irrelevant message")), false),
	)
	t.Run(
		"azdo_error_bad_request_status_VS800075",
		test(azdoErr(NewInt(http.StatusBadRequest), NewString("something about VS800075")), true),
	)
}

func TestResponseWasStatusCode(t *testing.T) {
	test := func(input error, expected bool) func(*testing.T) {
		return func(t *testing.T) {
			res := ResponseWasStatusCode(input, http.StatusConflict)
			require.Equal(t, expected, res)
		}
	}

	t.Run("nil_error", test(nil, false))
	t.Run("non_azdo_error", test(fmt.Errorf("not an AZDO error"), false))
	t.Run("azdo_error_nil_status", test(azdoErr(nil, nil), false))
	t.Run("azdo_error_matching_status", test(azdoErr(NewInt(http.StatusConflict), nil), true))
	t.Run("azdo_error_non_matching_status", test(azdoErr(NewInt(http.StatusBadRequest), nil), false))
}

func TestResponseContainsStatusMessage(t *testing.T) {
	test := func(input error, expected bool) func(*testing.T) {
		return func(t *testing.T) {
			res := ResponseContainsStatusMessage(input, "some string")
			require.Equal(t, expected, res)
		}
	}

	t.Run("nil_error", test(nil, false))
	t.Run("non_azdo_error", test(fmt.Errorf("not an AZDO error"), false))
	t.Run("azdo_error_nil_message", test(azdoErr(nil, nil), false))
	t.Run("azdo_error_containing_message", test(azdoErr(nil, NewString("I contain some string, yes")), true))
	t.Run("azdo_error_not_containing_message", test(azdoErr(nil, NewString("I do not contain the string")), false))
}
