---
title: Segment Clone
enterprise: false
time: 25
level: Beginner
modules:
  - amplitude
  - mailchimp
  - segment
links:
  - name: Repository
    url: https://github.com/nunchistudio/fragment
---

# Segment Clone

The result of this tutorial is called Fragment. It is available as an [open-source
"Template" repository on GitHub](https://github.com/nunchistudio/fragment).

In this tutorial we are going to create a consistent and reliable data engineering
solution with Blacksmith for in-house Customer Data Platform. Instead of creating
a specification from scratch, we will rely on the [Segment
Specification](https://segment.com/docs/connections/spec/). The API created will
therefore be compatible and consumable by the Segment SDKs, just like the [Segment
API](https://segment.com/docs/connections/sources/catalog/libraries/server/http-api/)
you already use and love.

## Creating the application

The first step is to create a Blacksmith application. The best way to do it is
by generating a new one from the CLI. Make sure to generate it inside your `GOPATH`
to avoid confusion when working with dependencies.
```bash
$ cd $GOPATH/src/github.com/<username>
$ blacksmith generate application fragment --path ./fragment

```

**Related resources:**
- Getting started >
  [Installation](/blacksmith/start/onboarding/install)
- Getting started >
  [Creating an application](/blacksmith/start/firstapp/create)
- CLI reference >
  [`generate application`](/blacksmith/cli/generate-application)

## Creating the source `rest`

The Segment API exposes six endpoints via a HTTP API. To mimic this behavior and
features, we'll create a source `rest`, which will be the home of our HTTP routes.
```bash
$ blacksmith generate source --name rest \
  --path ./sources/rest

```

**Related resources:**
- Guides for ETL with Go >
  [Sources](/blacksmith/etl/extraction/sources)
- CLI reference >
  [`generate source`](/blacksmith/cli/generate-source)

### The six type of events

The Segment Specification relies on six type of events. Each event exposes a HTTP
endpoint, which will be registered in the Blacksmith application as triggers using
the `http` mode under the `rest` source we just created:
- **Identify:** Who is the customer?

  **Method:** `POST`
  **Path:** `/v1/identify`

- **Track:** What are they doing?

  **Method:** `POST`
  **Path:** `/v1/track`

- **Page:** What web page are they on?

  **Method:** `POST`
  **Path:** `/v1/page`

- **Screen:** What application screen are they on?

  **Method:** `POST`
  **Path:** `/v1/screen`

- **Group:** What account or organization are they part of?

  **Method:** `POST`
  **Path:** `/v1/group`

- **Alias:** What was their past identity?

  **Method:** `POST`
  **Path:** `/v1/alias`

### The batch endpoint

In addition to these endpoints, the `batch` method lets you send a series of the
defined methods in a single batch, saving on outbound requests.

### Zoom in on the `Identify` trigger

Let's create a trigger for the `Identify` method inside the source's directory:
```bash
$ blacksmith generate trigger --name identify --mode http \
  --path ./sources/rest

```

**Related resources:**
- Guides for ETL with Go >
  [Triggers](/blacksmith/etl/extraction/triggers)
- Guides for ETL with Go >
  [Triggers HTTP](/blacksmith/etl/extraction/triggers-http)
- CLI reference >
  [`generate trigger`](/blacksmith/cli/generate-trigger)

Then, register this trigger in its parent source:
```go
func (s *REST) Triggers() map[string]source.Trigger {
  return map[string]source.Trigger{
    "identify": Identify{
      env: s.env,
    },
  }
}

```

The generated file for the trigger looks like this:
```go
package rest

import (
  "encoding/json"
  "net/http"
  "strings"
  "time"

  "github.com/nunchistudio/blacksmith/source"
  "github.com/nunchistudio/blacksmith/helper/errors"

  "gopkg.in/segmentio/analytics-go.v3"
)

/*
Identify implements the Blacksmith source.Trigger interface for the trigger
"identify". It holds the complete payload structure sent by an event and that
will be received by the gateway.
*/
type Identify struct {
  env *Options

  analytics.Identify
}

/*
String returns the string representation of the trigger Identify.
*/
func (t Identify) String() string {
  return "identify"
}

/*
Mode allows to register the trigger as a HTTP route. This means, every
time a "POST" request is executed against the route "/v1/identify", the
Extract function will run.
*/
func (t Identify) Mode() *source.Mode {
  return &source.Mode{
    Mode: source.ModeHTTP,
    UsingHTTP: &source.Route{
      Methods:  []string{"POST"},
      Path:     t.env.Prefix + "/v1/identify",
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
func (t Identify) Extract(tk *source.Toolkit, req *http.Request) (*source.Event, error) {

  // Create an empty payload, catch unwanted fields, and unmarshal it.
  // Return an error if any occured.
  var payload Identify
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
          Path:    []string{"analytics", "Identify"},
        },
      },
    }
  }

  // Add the current timestamp if none was provided.
  if payload.Timestamp.IsZero() {
    payload.Timestamp = time.Now().UTC()
  }

  // Validate the payload using the Segment official library.
  err := payload.Validate()
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
  if payload.Context != nil {
    ctx, err = payload.Context.MarshalJSON()
    if err != nil {
      return nil, &errors.Error{
        StatusCode: 400,
        Message:    "Bad Request",
      }
    }
  }

  // Try to marshal the data from the request payload.
  var data []byte
  if payload.Traits != nil {
    data, err = json.Marshal(&payload.Traits)
    if err != nil {
      return nil, &errors.Error{
        StatusCode: 400,
        Message:    "Bad Request",
      }
    }
  }

  // Return the context and data. See the section below for details about
  // the flow to execute.
  return &source.Event{
    Version: "v1.0",
    Context: ctx,
    Data:    data,
    SentAt:  &payload.Timestamp,
  }, nil
}

```

## The Segment flow

The `segment` module exposes a package named `segmentflow`. It is a collection
of flows — one for each type of event — Transforming an event's data for all the
destinations following the Segment Specification. As of now we support `amplitude`,
`mailchimp`, and `segment`. This means if a destination with one of the following
identifiers is registered in the application, the Segment flow will automatically
create the appropriate jobs.

The `Extract` function of a trigger can return this flow to execute the appropriate
actions. Given the `Identify` trigger saw earlier:
```go
func (t Identify) Extract(tk *source.Toolkit, req *http.Request) (*source.Event, error) {

  // ...

  return &source.Event{
    Version: "v1.0",
    Context: ctx,
    Data:    data,
    Flows: []flow.Flow{
      &segmentflow.Identify{
        Identify: payload.Identify,
      },
    },
    SentAt: &payload.Timestamp,
  }, nil
}

```

## Configuring the destinations

We offer different Go modules following the Segment Specification. Each of them
includes a package for Loading data to its destination:
- [Module `amplitude`](/blacksmith/tutorials/amplitude) exposes a package
  `amplitudedestination` for Loading data to Amplitude using the HTTP API;
- [Module `mailchimp`](/blacksmith/tutorials/mailchimp) exposes a package
  `mailchimpdestination` for Loading data to Mailchimp using the HTTP API;
- [Module `segment`](/blacksmith/tutorials/segment) exposes a package
  `segmentdestination` for Loading data to Segment using the HTTP API.

Each destination can be configured in the
[`*blacksmith.Options`](https://pkg.go.dev/github.com/nunchistudio/blacksmith?tab=doc#Options):
```go
package main

import (
  "github.com/nunchistudio/blacksmith"
  "github.com/nunchistudio/blacksmith/destination"

  "github.com/nunchistudio/blacksmith-modules/amplitude/amplitudedestination"
  "github.com/nunchistudio/blacksmith-modules/mailchimp/mailchimpdestination"
  "github.com/nunchistudio/blacksmith-modules/segment/segmentdestination"
)

func Init() *blacksmith.Options {

  var options = &blacksmith.Options{

    // ...

    Destinations: []destination.Destination{
      amplitudedestination.New(&amplitudedestination.Options{
        APIKey: os.Getenv("AMPLITUDE_API_KEY"),
      }),
      mailchimpdestination.New(&mailchimpdestination.Options{
        APIKey:            os.Getenv("MAILCHIMP_API_KEY"),
        DatacenterID:      os.Getenv("MAILCHIMP_DATACENTER"),
        AudienceID:        os.Getenv("MAILCHIMP_AUDIENCE"),
        EnableDoubleOptIn: true,
      }),
      segmentdestination.New(&segmentdestination.Options{
        WriteKey: os.Getenv("SEGMENT_WRITE_KEY"),
      }),
    },
  }

  return options
}

```

## Environment variables

When creating the application, all the Blacksmith options are already set to work
in a development environment. You only to configure the sources and destinations
you wish to work with.

The easiest way to set the options is to use environment variables. You can copy
the file `.env.example` into a new one named `.env`. It contains all the necessary
keys to be set. The `.env` file will not be versioned since it is ignored by the
`.gitignore`.

Given the variables desired for the destinations, it should look like this:
```bash
AMPLITUDE_API_KEY=
MAILCHIMP_API_KEY=
MAILCHIMP_DATACENTER=
MAILCHIMP_AUDIENCE=
SEGMENT_WRITE_KEY=
NATS_SERVER_URL=nats://host.docker.internal:4222
POSTGRES_STORE_URL=postgres://fragment:fragment@host.docker.internal:5432/fragment?sslmode=disable
POSTGRES_SUPERVISOR_URL=postgres://fragment:fragment@host.docker.internal:5432/fragment?sslmode=disable
POSTGRES_WANDERER_URL=postgres://fragment:fragment@host.docker.internal:5432/fragment?sslmode=disable

```

Enjoy the Enterprise Edition for free with the following license details!
```bash
BLACKSMITH_LICENSE_KEY=FRAGMENT
BLACKSMITH_LICENSE_TOKEN=FRAGMENT

```

## Docker environment

[As mentioned in the "Getting started"](/blacksmith/start/firstapp/docker),
Blacksmith leverages Docker for environment parity.

Fragment uses a `Docker-compose.yml` to make development a breeze. This file is
not required for running the application but is here for convenience to help you
get started even faster in development. It contains:
- a PostgreSQL database for using the `postgres` driver for the `store`,
  `supervisor`, and `wanderer` adapters;
- a NATS server for using the `nats` driver for the `pubsub` adapter.

You can customize the stack as much as you need, and run it with:
```bash
$ docker-compose up -d

```

## Running the application

Since the Docker stack is now up and running, we can run the Fragment application
with the Blacksmith CLI:
```bash
$ blacksmith start --bind 9090:9090 --bind 9091:9091

```

Or, you can leverage the `Dockerfile` with Docker to build and run the image:
```bash
$ docker build -t fragment ./
$ docker run --env-file .env -p 9090:9090 -p 9091:9091 fragment

```

## Your first event

Because the `rest` source follows the Segment API, we can leverage one of the Segment
client to test our application. In a Node.js application, one can use the [Segment
SDK](https://segment.com/docs/connections/sources/catalog/libraries/server/node/)
for example:
```js
var Analytics = require('analytics-node');
var analytics = new Analytics('THIS_IS_NOT_USED', {
  host: 'http://localhost:9090',
});

analytics.identify({
  userId: '7439857439857',
  traits: {
    firstName: 'John',
    lastName: 'Doe',
    email: 'johndoe@example.com'
  }
});

```

Or more straightforward with `curl`:
```bash
$ curl --request POST \
  --url http://localhost:9090/v1/identify \
  --header 'Content-Type: application/json' \
  --data '{
    "userId": "7439857439857",
    "traits": {
      "firstName": "John",
      "lastName": "Doe",
      "email": "johndoe@example.com"
    }
  }'

```

You can now visit the Blacksmith Dashboard at <http://localhost:9091/admin> and
view all the events and jobs flowing in your Customer Data Platform. If we look
for jobs, we can see that 3 jobs have been created:

![Blacksmith Dashboard](https://nunchi.studio/images/blacksmith/segment.001.png)

You can click on a job to see its details. We can understand why the `Identify`
job failed in Mailchimp:

![Blacksmith Dashboard](https://nunchi.studio/images/blacksmith/segment.002.png)

The job is marked as `discarded` and not as `failed` because the job can never
succeed.

## Going further

In this tutorial we set the foundations for creating a Customer Data Platform using
Blacksmith and following the Segment Specification.

Here are some ideas for going further:
- Add a source with triggers leveraging the Pub / Sub adapter ([mode
  `sub`](/blacksmith/etl/extraction/triggers-sub)) in order to collect events
  from a message queue in addition to the current triggers (which are of [mode
  `http`](/blacksmith/etl/extraction/triggers-http)).
- Triggers can return an action to insert data in a warehouse using the [`sqlike`
  module](/blacksmith/tutorials/sqlike).
