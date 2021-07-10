package utils

import (
	"encoding/base64"
	"fmt"
)

func StringIsBase64EncodedAndNotEmpty(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	content, err := base64.StdEncoding.DecodeString(v)

	if err != nil {
		return nil, []error{fmt.Errorf("expected %q to be base64-encoded", k)}
	}

	if len(content) == 0 {
		return nil, []error{fmt.Errorf("expected %q to not be empty", k)}
	}

	return nil, nil
}
