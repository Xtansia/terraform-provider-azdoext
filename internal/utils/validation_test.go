package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringIsBase64Encoded(t *testing.T) {
	t.Run(
		"input_not_string", func(t *testing.T) {
			warns, errors := StringIsBase64Encoded(32, "someKey")
			require.Empty(t, warns)
			require.Len(t, errors, 1)
			require.EqualError(t, errors[0], "expected type of \"someKey\" to be string")
		},
	)
	t.Run(
		"input_not_base64", func(t *testing.T) {
			warns, errors := StringIsBase64Encoded("some gibberish", "someKey")
			require.Empty(t, warns)
			require.Len(t, errors, 1)
			require.EqualError(t, errors[0], "expected \"someKey\" to be base64-encoded")
		},
	)
	t.Run(
		"input_is_base64", func(t *testing.T) {
			warns, errors := StringIsBase64Encoded("c29tZSBnaWJiZXJpc2g=", "someKey")
			require.Empty(t, warns)
			require.Empty(t, errors)
		},
	)
}
