package v1

import (
	"fmt"
	"testing"
)

func TestHttpParams(t *testing.T) {
	hp := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	httpParams := HttpParams(hp)
	encodeParam := httpParams.Encode()
	fmt.Println(encodeParam)
}
