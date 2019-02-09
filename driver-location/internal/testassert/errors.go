package testassert

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.fraixed.es/errors"
)

// ErrorWithCode asserts if an error has the specified error code and metadata.
func ErrorWithCode(t *testing.T, err error, c errors.Code, mds ...errors.MD) bool {
	t.Helper()

	if !assert.Error(t, err) {
		return false
	}

	if !assert.True(t, errors.Is(err, c), "unexpected error code") {
		return false
	}

	var emsg = fmt.Sprintf("%+v", err)
	for _, md := range mds {
		if !assert.Contains(t, emsg, md.K, "unexpected metadata") {
			return false
		}

		if !assert.Contains(t, emsg, fmt.Sprintf("%+v", md.V), "unexpected metadata") {
			return false
		}
	}

	return true
}
