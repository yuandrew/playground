package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	tallyprom "github.com/uber-go/tally/v4/prometheus"

	"go.temporal.io/sdk/client"
	sdktally "go.temporal.io/sdk/contrib/tally"
)

func main() {
	metricsHandler := setupPrometheus()
	// These lines in any order will cause an error where the metrics have different labels
	metricsHandler.Counter("testing").Inc(1)
	metricsHandler.WithTags(map[string]string{"someName": "someValue"}).Counter("testing").Inc(1)
}

// Just setup stuff for Prometheus below here
var PrometheusSanitizeOptions = tally.SanitizeOptions{
	NameCharacters:       tally.ValidCharacters{Ranges: tally.AlphanumericRange, Characters: []rune{'_'}},
	KeyCharacters:        tally.ValidCharacters{Ranges: tally.AlphanumericRange, Characters: []rune{'_'}},
	ValueCharacters:      tally.ValidCharacters{Ranges: tally.AlphanumericRange, Characters: []rune{'_'}},
	ReplacementCharacter: tally.DefaultReplacementCharacter,
}

func setupPrometheus() client.MetricsHandler {
	cfg := tallyprom.Configuration{
		ListenAddress: "localhost:9090",
	}
	reporter, _ := cfg.NewReporter(tallyprom.ConfigurationOptions{
		Registry: prometheus.NewRegistry(),
		OnError: func(err error) {
			fmt.Printf("BOOM: %+v\n", err)
		},
	})
	scopeOpts := tally.ScopeOptions{CachedReporter: reporter, SanitizeOptions: &PrometheusSanitizeOptions}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)
	return sdktally.NewMetricsHandler(scope)
}
