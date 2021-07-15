package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHumaniseList(t *testing.T) {
	test := func(input []string, expected string) func(*testing.T) {
		return func(t *testing.T) {
			result := HumaniseList(input)
			require.Equal(t, expected, result)
		}
	}

	t.Run("none", test(nil, ""))
	t.Run("one", test([]string{"foobar"}, "foobar"))
	t.Run("two", test([]string{"foobar", "bazbar"}, "foobar & bazbar"))
	t.Run("many", test([]string{"foobar", "bazbar", "barrybar", "garrybar"}, "foobar, bazbar, barrybar & garrybar"))
}
