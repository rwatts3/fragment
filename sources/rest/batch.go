package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nunchistudio/blacksmith/helper/errors"
	"github.com/nunchistudio/blacksmith/source"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Batch implements the Blacksmith source.Trigger interface for the trigger
"batch". It holds the complete payload structure sent by an event and that
will be received by the gateway.
*/
type Batch struct {
	env *Options

	Events    interface{}        `json:"batch"`
	Context   *analytics.Context `json:"context,omitempty"`
	Timestamp time.Time          `json:"timestamp,omitempty"`
}

/*
String returns the string representation of the trigger Batch.
*/
func (t Batch) String() string {
	return "batch"
}

/*
Mode allows to register the trigger as a HTTP route. This means, every
time a "POST" request is executed against the route "/v1/batch", the
Extract function will run.
*/
func (t Batch) Mode() *source.Mode {
	return &source.Mode{
		Mode: source.ModeHTTP,
		UsingHTTP: &source.Route{
			Methods:  []string{"POST"},
			Path:     t.env.Prefix + "/v1/batch",
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
func (t Batch) Extract(tk *source.Toolkit, req *http.Request) (*source.Event, error) {

	// Create an empty payload, catch unwanted fields, and unmarshal it.
	// Return an error if any occured.
	var payload Batch
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: 400,
			Message:    "Bad Request",
			Validations: []errors.Validation{
				{
					Message: err.Error(),
					Path:    []string{"analytics", "Batch"},
				},
			},
		}
	}

	// Add the current timestamp if none was provided.
	if payload.Timestamp.IsZero() {
		payload.Timestamp = time.Now().UTC()
	}

	// The events are registered in a slice. Since we do not know the type of
	// each one, consider them as interfaces.
	events := payload.Events.([]interface{})

	// Go through each event passed in the payload. We create a slice of sub-events
	// to register every events found in the batch object. By returning sub-events,
	// these said sub-events will have this batch event as their parent event.
	var subEvents = []*source.SubEvent{}
	for i := range events {
		var fail *errors.Error
		var subevent *source.SubEvent

		// Make sure Go can parse the map of the event.
		event, ok := events[i].(map[string]interface{})
		if ok == false {
			continue
		}

		// Marshal the event so we can then unmarshal it.
		b, err := json.Marshal(event)
		if err != nil {
			continue
		}

		// Determine the event type based on the "type" key.
		eventType, ok := event["type"].(string)
		if ok == false {
			continue
		}

		// Unmarshal the event with the appropriate struct and create the
		// flow for the corresponding event.
		switch eventType {
		case "identify":
			var e Identify
			json.Unmarshal(b, &e)
			subevent, fail = e.marshal()
			if fail != nil {
				tk.Logger.Error(fail)
			}

		case "track":
			var e Track
			json.Unmarshal(b, &e)
			subevent, fail = e.marshal()
			if fail != nil {
				tk.Logger.Error(fail)
			}

		case "group":
			var e Group
			json.Unmarshal(b, &e)
			subevent, fail = e.marshal()
			if fail != nil {
				tk.Logger.Error(fail)
			}

		case "alias":
			var e Alias
			json.Unmarshal(b, &e)
			subevent, fail = e.marshal()
			if fail != nil {
				tk.Logger.Error(fail)
			}

		case "page":
			var e Page
			json.Unmarshal(b, &e)
			subevent, fail = e.marshal()
			if fail != nil {
				tk.Logger.Error(fail)
			}

		case "screen":
			var e Screen
			json.Unmarshal(b, &e)
			subevent, fail = e.marshal()
			if fail != nil {
				tk.Logger.Error(fail)
			}

		default:
			continue
		}

		// If a sub-event is present, add it to the slice of sub-events to process.
		if subevent != nil {
			subEvents = append(subEvents, subevent)
		}
	}

	// Try to marshal the context from the request payload.
	var ctx []byte
	if payload.Context != nil {
		ctx, err = payload.Context.MarshalJSON()
		if err != nil {
			return nil, &errors.Error{
				StatusCode: 400,
				Message:    "Bad Request",
			}
		}
	}

	// Return the context, data, and a collection of sub-events to process.
	return &source.Event{
		Version:   "v1.0",
		Context:   ctx,
		SubEvents: subEvents,
		SentAt:    &payload.Timestamp,
	}, nil
}
