package rest

import (
	"strings"

	"github.com/nunchistudio/blacksmith/helper/errors"
)

/*
Options is the options the source can take as an input to be configured.
*/
type Options struct {

	// ShowMeta is used to display (or not) the metadata in the HTTP response,
	// such as the event's context and jobs details.
	ShowMeta bool

	// ShowData is used to display (or not) the data in the HTTP response. It should
	// be disabled if any sensitive data can be returned, such as private tokens.
	ShowData bool

	// Prefix allows to prefix the endpoints exposed by the source.
	//
	// Example: "/cdp"
	Prefix string
}

/*
validate ensures the options passed to initialize the source are valid.
*/
func (env *Options) validate() error {
	fail := &errors.Error{
		Message:     "source/rest: Failed to load",
		Validations: []errors.Validation{},
	}

	if env.Prefix != "" {
		if !strings.HasPrefix(env.Prefix, "/") {
			fail.Validations = append(fail.Validations, errors.Validation{
				Message: "Source prefix must start with a '/'",
				Path:    []string{"Options", "Sources", "rest", "Prefix"},
			})
		}

		if strings.HasSuffix(env.Prefix, "/") {
			fail.Validations = append(fail.Validations, errors.Validation{
				Message: "Source prefix must not end with a '/'",
				Path:    []string{"Options", "Sources", "rest", "Prefix"},
			})
		}
	}

	if len(fail.Validations) > 0 {
		return fail
	}

	return nil
}
