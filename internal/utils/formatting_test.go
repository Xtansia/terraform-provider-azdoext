package utils

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestHumaniseList_EmptyInput(t *testing.T) {
    result := HumaniseList(nil)
    require.Equal(t, "", result, "Result should be an empty string")
}

func TestHumaniseList_SingleInput(t *testing.T) {
    result := HumaniseList([]string{"foobar"})
    require.Equal(t, "foobar", result, "Result should just be the singular input")
}

func TestHumaniseList_TwoInputs(t *testing.T) {
    result := HumaniseList([]string{"foobar", "bazbar"})
    require.Equal(t, "foobar & bazbar", result, "Result should be two inputs joined by an ampersand")
}

func TestHumaniseList_ManyInputs(t *testing.T) {
    result := HumaniseList([]string{"foobar", "bazbar", "barrybar", "garrybar"})
    require.Equal(t, "foobar, bazbar, barrybar & garrybar", result, "Result should be all except last joined by command and last joined by an ampersand")
}