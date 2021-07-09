package utils

import (
    "fmt"
    "strings"
)

func HumaniseList(elems []string) string {
    l := len(elems)
    switch l {
    case 0:
        return ""
    case 1:
        return elems[0]
    default:
        return fmt.Sprintf("%s & %s", strings.Join(elems[: l - 1], ", "), elems[l - 1])
    }
}