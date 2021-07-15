package utils

import (
	"encoding/base64"
	"fmt"
)

func StringIsBase64Encoded(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	_, err := base64.StdEncoding.DecodeString(v)

	if err != nil {
		return nil, []error{fmt.Errorf("expected %q to be base64-encoded", k)}
	}

	return nil, nil
}
