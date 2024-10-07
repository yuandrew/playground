package main

import (
	"context"
	"log"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	sdktally "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a Hello World workflow definition.
func Workflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	// Start a child workflow
	// cwo := workflow.ChildWorkflowOptions{
	// 	WorkflowID: "child-workflow",
	// }
	// ctx = workflow.WithChildOptions(ctx, cwo)

	// future := workflow.ExecuteChildWorkflow(ctx, ChildWorkflow)

	// Try to send a signal to the child workflow before it's started
	// This should create an invalid command sequence
	// workflow.SignalExternalWorkflow(ctx, "child-workflow", "", "signal-name", "data")

	// var result string
	// if err := future.Get(ctx, &result); err != nil {
	// 	return "", err
	// }
	// return "", workflow.Await(ctx, func() bool {
	// 	return false // Keep workflow running
	// })

	// OTHER SCENARIO
	// panic("try to fail a workflow")
	// return "", errors.New("try to fail a workflow")
	var result string
	err := workflow.ExecuteActivity(ctx, Activity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	// err = workflow.ExecuteActivity(ctx, Activity1, name).Get(ctx, &result)
	// if err != nil {
	// 	logger.Error("Activity failed.", "Error", err)
	// 	return "", err
	// }

	err = workflow.ExecuteActivity(ctx, Activity2, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", "result", result)

	return result, nil
}

func NewPrometheusScope(c prometheus.Configuration) tally.Scope {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				log.Println("error in prometheus reporter", err)
			},
		},
	)
	if err != nil {
		log.Fatalln("error creating prometheus reporter", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
		Prefix:          "temporal_samples",
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)

	log.Println("prometheus metrics scope created")
	return scope
}

func ChildWorkflow(ctx workflow.Context) (string, error) {
	return "result", nil
}

func Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity1", "name", name)
	return "Hello1 " + name + "!", nil
}
func Activity1(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity1", "name", name)
	// return "", errors.New("try to fail a workflow")
	return "Hello1 " + name + "!", nil
}
func Activity2(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity1", "name", name)
	return "Hello1 " + name + "!", nil
}

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		MetricsHandler: sdktally.NewMetricsHandler(NewPrometheusScope(prometheus.Configuration{
			ListenAddress: "0.0.0.0:9090",
			TimerType:     "histogram",
		})),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "hello-world", worker.Options{})

	w.RegisterWorkflow(Workflow)
	w.RegisterActivity(Activity)
	w.RegisterActivity(Activity1)
	w.RegisterActivity(Activity2)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
