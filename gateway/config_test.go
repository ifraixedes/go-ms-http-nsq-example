package gateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigFromYAML(t *testing.T) {
	var tcases = []struct {
		desc   string
		d      []byte
		assert func(*testing.T, Config, error)
	}{
		{
			desc: "OK: 2 elements",
			d:    []byte(configYAMLOK),
			assert: func(t *testing.T, c Config, err error) {
				require.NoError(t, err)

				assert.Equal(t, Config{
					setDriverLocation: configEndpoint{
						Path:   "/drivers/:id/locations",
						Method: http.MethodPatch,
						NSQ: &configNSQ{
							Topic: "locations",
						},
					},
					getDriver: configEndpoint{
						Path:   "/drivers/:id",
						Method: http.MethodGet,
						HTTP: &configHTTP{
							Host: "zombie-driver",
						},
					},
				}, c)
			},
		},
		{
			desc: "OK: 2 elements other order",
			d:    []byte(configYAMLOKOtherOrder),
			assert: func(t *testing.T, c Config, err error) {
				require.NoError(t, err)

				assert.Equal(t, Config{
					setDriverLocation: configEndpoint{
						Path:   "/drivers/:id/locations",
						Method: http.MethodPost,
						NSQ: &configNSQ{
							Topic: "locations",
						},
					},
					getDriver: configEndpoint{
						Path:   "/drivers/:id",
						Method: http.MethodGet,
						HTTP: &configHTTP{
							Host: "zombie-driver",
						},
					},
				}, c)
			},
		},
		{
			desc: "OK: more than 2 elements",
			d:    []byte(configYAMLOKMoreElements),
			assert: func(t *testing.T, c Config, err error) {
				require.NoError(t, err)

				assert.Equal(t, Config{
					setDriverLocation: configEndpoint{
						Path:   "/drivers/:id/locations",
						Method: http.MethodPatch,
						NSQ: &configNSQ{
							Topic: "locations",
						},
					},
					getDriver: configEndpoint{
						Path:   "/drivers/:id",
						Method: http.MethodGet,
						HTTP: &configHTTP{
							Host: "zombie-driver",
						},
					},
				}, c)
			},
		},
		{
			desc: "error: invalid method",
			d:    []byte(configYAMLErrorInvalidMethod),
			assert: func(t *testing.T, c Config, err error) {
				assert.Error(t, err)
			},
		},
		{
			desc: "error: missing element",
			d:    []byte(configYAMLErrorMissingElement),
			assert: func(t *testing.T, c Config, err error) {
				assert.Error(t, err)
			},
		},
		{
			desc: "error: missing ID",
			d:    []byte(configYAMLErrorMissingID),
			assert: func(t *testing.T, c Config, err error) {
				assert.Error(t, err)
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			var c, err = NewConfigFromYAML(tc.d)
			tc.assert(t, c, err)
		})
	}
}

const configYAMLOK = `
urls:
  -
    path: "/drivers/:id/locations"
    method: "PATCH"
    nsq:
      topic: "locations"
  -
    path: "/drivers/:id"
    method: "GET"
    http:
      host: "zombie-driver"
`

const configYAMLOKOtherOrder = `
urls:
  -
    path: "/drivers/:id"
    method: "GET"
    http:
      host: "zombie-driver"
  -
    path: "/drivers/:id/locations"
    method: "POST"
    nsq:
      topic: "locations"
`

const configYAMLOKMoreElements = `
urls:
  -
    path: "/drivers/:id/locations"
    method: "PATCH"
    nsq:
      topic: "locations"
  -
    path: "/drivers/:id"
    method: "GET"
    http:
      host: "zombie-driver"
  -
    path: "/drivers/:id/locations"
    method: "GET"
`

const configYAMLErrorInvalidMethod = `
urls:
  -
    path: "/drivers/:id/locations"
    method: "GET"
    nsq:
      topic: "locations"
  -
    path: "/drivers/:id"
    method: "GET"
    http:
      host: "zombie-driver"
`

const configYAMLErrorMissingElement = `
urls:
  -
    path: "/drivers/:id"
    method: "GET"
    http:
      host: "zombie-driver"
`

const configYAMLErrorMissingID = `
urls:
  -
    path: "/drivers/id/locations"
    method: "PATCH"
    nsq:
      topic: "locations"
  -
    path: "/drivers/:id"
    method: "GET"
    http:
      host: "zombie-driver"
`
