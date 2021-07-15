package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	mapper := func(s string) string {
		return "foo" + s
	}

	t.Run(
		"empty", func(t *testing.T) {
			result := MapStrings(nil, mapper)
			require.Len(t, result, 0, "Result should be empty")
		},
	)
	t.Run(
		"many", func(t *testing.T) {
			result := MapStrings([]string{"bar", "baz", "ban"}, mapper)
			require.Equal(
				t, []string{"foobar", "foobaz", "fooban"}, result,
				"Result should contain all input elements with mapper applied to them",
			)
		},
	)
}
