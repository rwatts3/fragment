package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/nunchistudio/blacksmith/flow"
	"github.com/nunchistudio/blacksmith/helper/errors"
	"github.com/nunchistudio/blacksmith/source"

	"github.com/nunchistudio/blacksmith-modules/segment/segmentflow"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Alias implements the Blacksmith source.Trigger interface for the trigger
"alias". It holds the complete payload structure sent by an event and that
will be received by the gateway.
*/
type Alias struct {
	env *Options

	analytics.Alias
}

/*
String returns the string representation of the trigger Alias.
*/
func (t Alias) String() string {
	return "alias"
}

/*
Mode allows to register the trigger as a HTTP route. This means, every
time a "POST" request is executed against the route "/v1/alias", the
Extract function will run.
*/
func (t Alias) Mode() *source.Mode {
	return &source.Mode{
		Mode: source.ModeHTTP,
		UsingHTTP: &source.Route{
			Methods:  []string{"POST"},
			Path:     t.env.Prefix + "/v1/alias",
			ShowMeta: t.env.ShowMeta,
			ShowData: t.env.ShowData,
		},
	}
}

/*
Extract is the function being run when the HTTP route is triggered. It is
in charge of the "E" in the ETL process: Extract the data from the source.

The function allows to return data to flows. It is the "T" in the ETL
process: it transforms the payload from the source's trigger to given
destinations' actions.
*/
func (t Alias) Extract(tk *source.Toolkit, req *http.Request) (*source.Event, error) {

	// Create an empty payload, catch unwanted fields, and unmarshal it.
	// Return an error if any occured.
	var payload Alias
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: 400,
			Message:    "Bad Request",
			Validations: []errors.Validation{
				{
					Message: err.Error(),
					Path:    []string{"analytics", "Alias"},
				},
			},
		}
	}

	// Marshal the event's context and data using the internal method which also
	// returns the flows to run.
	subevent, fail := payload.marshal()
	if fail != nil {
		return nil, fail
	}

	// Return the context, data, and a collection of flows to run.
	return &source.Event{
		Version: "v1.0",
		Context: subevent.Context,
		Data:    subevent.Data,
		SentAt:  &payload.Timestamp,
		Flows:   subevent.Flows,
	}, nil
}

/*
marshal is an internal method to return the event's context and data alongside
flows to run for a given Alias receiver. It is handy for validating an event
triggered by the gateway or inheriting from a Batch trigger.
*/
func (t Alias) marshal() (*source.SubEvent, *errors.Error) {

	// Add the current timestamp if none was provided.
	if t.Timestamp.IsZero() {
		t.Timestamp = time.Now().UTC()
	}

	// Validate the payload using the Segment official library.
	err := t.Validate()
	if err != nil {
		fail := err.(analytics.FieldError)
		return nil, &errors.Error{
			StatusCode: 400,
			Message:    "Bad Request",
			Validations: []errors.Validation{
				{
					Message: fail.Name + " must be set",
					Path:    append(strings.Split(fail.Type, "."), fail.Name),
				},
			},
		}
	}

	// Try to marshal the context from the request payload.
	var ctx []byte
	if t.Context != nil {
		ctx, err = t.Context.MarshalJSON()
		if err != nil {
			return nil, &errors.Error{
				StatusCode: 400,
				Message:    "Bad Request",
			}
		}
	}

	// Return the context, data, and a collection of flows to run.
	return &source.SubEvent{
		Trigger: "alias",
		Context: ctx,
		Flows: []flow.Flow{
			&segmentflow.Alias{
				Alias: t.Alias,
			},
		},
	}, nil
}
