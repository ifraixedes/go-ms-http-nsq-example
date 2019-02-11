package gateway

import (
	"errors"
	"net/http"
	"regexp"

	yaml "gopkg.in/yaml.v2"
)

var pathIDPattern = regexp.MustCompile(`\/:id(\/|$)`)

// Config is used for configuring the gateway. The zero value isn't a valid
// configuration, hence any of the constructor functions must be used for
// initializing one.
//
// NSQdAddr is optional
type Config struct {
	NSQdAddr          string
	setDriverLocation configEndpoint
	getDriver         configEndpoint
}

// NewConfigFromYAML returns a configuration initilized with the passed YAML.
//
// YAML doesn't contain any other configuraton rather than endpoints, so the
// other fields of the configuration must be set after it.
//
// The YAML must contain a property named 'urls' which is an array that must
// have 2 members one with a NSQ configuration and another with an HTTP
// configuration; the order of this elements isn't important and if there are
// more elements isn't an error, they are basically ignored. Both elements must
// have a valid 'path' and 'method' and both paths must contains a ':id' part
// on it. The method of the NSQ member must 'POST', 'PUT' or 'PATCH', the HTTP
// member must be 'GET', 'POST', 'PUT or 'PATCH'.
//
// Example of a correct configuration:
//
// urls:
//   -
//     path: "/drivers/:id/locations"
//     method: "PATCH"
//     nsq:
//       topic: "locations"
//   -
//     path: "/drivers/:id"
//     method: "GET"
//     http:
//       host: "zombie-driver"
func NewConfigFromYAML(d []byte) (Config, error) {
	var uc = struct {
		URLs []configEndpoint `yaml:"urls"`
	}{}

	var err = yaml.UnmarshalStrict(d, &uc)
	if err != nil {
		return Config{}, err
	}

	var eps = uc.URLs
	switch {
	case len(eps) < 2:
		return Config{}, errors.New("invalid configuration, missing endpoints")
	case len(eps) > 2:
		eps = eps[:2]
	}

	var c Config
	for _, e := range eps {
		if e.NSQ == nil && e.HTTP == nil {
			return Config{}, errors.New("Invalid configuration: no NSQL nor HTTP endpoint")
		}

		if len(pathIDPattern.FindAllString(e.Path, 2)) != 1 {
			return Config{}, errors.New("Invalid configuration, path doens't contain ':id' or contains more than one")
		}

		if e.NSQ != nil {
			switch e.Method {
			case http.MethodPatch, http.MethodPost, http.MethodPut:
			default:
				return Config{}, errors.New("Invalid HTTP method for set locations endpoint")
			}

			if e.NSQ.Topic == "" {
				return Config{}, errors.New("Invalid set locations endpoint, topic cannot be empty")
			}

			c.setDriverLocation = e
			continue
		}

		switch e.Method {
		case http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodPut:
		default:
			return Config{}, errors.New("Invalid HTTP method for get driver endpoint")
		}

		if e.HTTP.Host == "" {
			return Config{}, errors.New("Invalid get driver endpoint, host cannot be empty")
		}

		c.getDriver = e
	}

	if c.getDriver == (configEndpoint{}) || c.setDriverLocation == (configEndpoint{}) {
		return Config{}, errors.New("Invalid configuration, missing one endpoint")
	}

	return c, nil
}

type configEndpoint struct {
	Path   string      `yaml:"path"`
	Method string      `yaml:"method"`
	NSQ    *configNSQ  `yaml:"nsq"`
	HTTP   *configHTTP `yaml:"http"`
}

type configNSQ struct {
	Topic string `yaml:"topic"`
}

type configHTTP struct {
	Host string `yaml:"host"`
}
