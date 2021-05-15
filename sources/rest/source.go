package rest

import (
	"time"

	"github.com/nunchistudio/blacksmith/flow/source"
	"github.com/nunchistudio/blacksmith/helper/logger"
)

/*
REST implements the Blacksmith source.Source interface for the source "rest".
*/
type REST struct {
	env     *Options
	options *source.Options
}

/*
New returns a valid Blacksmith source.Source for REST.
*/
func New(env *Options) source.Source {

	// Validate the environment options passed by the application.
	// Stop the process if any error is returned.
	if err := env.validate(); err != nil {
		logger.Default.Fatal(err)
		return nil
	}

	return &REST{
		env: env,
		options: &source.Options{
			DefaultVersion: "v1.0",
			Versions: map[string]time.Time{
				"v1.0": time.Time{},
			},
		},
	}
}

/*
String returns the string representation of the source REST.
*/
func (s *REST) String() string {
	return "rest"
}

/*
Options returns common source options for REST. They will be shared across
every triggers of this source, except when overridden.
*/
func (s *REST) Options() *source.Options {
	return s.options
}

/*
Triggers return a list of triggers the source REST is able to handle.
*/
func (s *REST) Triggers() map[string]source.Trigger {
	return map[string]source.Trigger{
		"identify": Identify{
			env: s.env,
		},
		"track": Track{
			env: s.env,
		},
		"group": Group{
			env: s.env,
		},
		"alias": Alias{
			env: s.env,
		},
		"page": Page{
			env: s.env,
		},
		"screen": Screen{
			env: s.env,
		},
		"batch": Batch{
			env: s.env,
		},
	}
}
