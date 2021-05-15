package main

import (
	"net/http"
	"os"

	"github.com/nunchistudio/blacksmith"
	"github.com/nunchistudio/blacksmith/adapter/pubsub"
	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/adapter/supervisor"
	"github.com/nunchistudio/blacksmith/adapter/wanderer"
	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/flow/source"
	"github.com/nunchistudio/blacksmith/service"

	"github.com/nunchistudio/blacksmith-modules/amplitude/amplitudedestination"
	"github.com/nunchistudio/blacksmith-modules/mailchimp/mailchimpdestination"
	"github.com/nunchistudio/blacksmith-modules/segment/segmentdestination"

	"github.com/nunchistudio/fragment/sources/rest"

	"github.com/rs/cors"
)

/*
Init is the entrypoint of Fragment, a Blacksmith application. It is used by
the Blacksmith CLI to load the application as a Go plugin.
*/
func Init() *blacksmith.Options {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	var options = &blacksmith.Options{
		Gateway: &service.Options{
			Admin: &service.Admin{
				Enabled:       false,
				WithDashboard: false,
			},
			Middleware: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					res.Header().Set("Content-Type", "application/json")
					res.Header().Set("Access-Control-Allow-Origin", "*")

					c.ServeHTTP(res, req, next.ServeHTTP)
				})
			},
		},
		Scheduler: &service.Options{
			Admin: &service.Admin{
				Enabled:       true,
				WithDashboard: true,
			},
		},

		Store: &store.Options{
			From: "postgres",
		},

		PubSub: &pubsub.Options{
			From: "nats",
		},
		Supervisor: &supervisor.Options{
			From: "postgres",
		},
		Wanderer: &wanderer.Options{
			From: "postgres",
		},

		Sources: []source.Source{
			rest.New(&rest.Options{
				ShowMeta: true,
				ShowData: true,
				Prefix:   "",
			}),
		},

		Destinations: []destination.Destination{
			amplitudedestination.New(&amplitudedestination.Options{
				Realtime: true,
				APIKey:   os.Getenv("AMPLITUDE_API_KEY"),
			}),
			mailchimpdestination.New(&mailchimpdestination.Options{
				Realtime:          true,
				APIKey:            os.Getenv("MAILCHIMP_API_KEY"),
				DatacenterID:      os.Getenv("MAILCHIMP_DATACENTER"),
				AudienceID:        os.Getenv("MAILCHIMP_AUDIENCE"),
				EnableDoubleOptIn: false,
			}),
			segmentdestination.New(&segmentdestination.Options{
				Realtime: true,
				WriteKey: os.Getenv("SEGMENT_WRITE_KEY"),
			}),
		},
	}

	return options
}
