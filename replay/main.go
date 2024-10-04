package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/history/v1"
	"go.temporal.io/sdk/client"
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
	// panic("try to fail a workflow")

	// return "", errors.New("try to fail a workflow")
	var result string
	err := workflow.ExecuteActivity(ctx, Activity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
	// time.Sleep(1 * time.Second)

	logger.Info("HelloWorld workflow completed.", "result", result)

	return result, nil
}

func Activity(ctx context.Context, name string) (string, error) {
	// logger := activity.GetLogger(ctx)
	// logger.Info("Activity", "name", name)
	return "Hellosasdf " + name + "!", nil
	// panic("FAIL LOCAL ACTIVITY")
}
func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// Replay workflow
	fmt.Println("Replaying workflow")
	// d7226966-0465-4e31-9fdd-36cc88c59446
	// 0f7858a6-0953-4be1-bac5-6af0f2e9a897
	err = ReplayWorkflow(context.Background(), c, "hello_world_workflowID", "45d87bda-4d79-4ce9-a206-4163f1ef759c")
	if err != nil {
		log.Fatalln("Unable to replay workflow", err)
	}
}

func GetWorkflowHistory(ctx context.Context, client client.Client, id, runID string) (*history.History, error) {
	var hist history.History
	iter := client.GetWorkflowHistory(ctx, id, runID, false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
	for iter.HasNext() {
		event, err := iter.Next()
		if err != nil {
			return nil, err
		}
		hist.Events = append(hist.Events, event)
	}
	return &hist, nil
}

func ReplayWorkflow(ctx context.Context, client client.Client, id, runID string) error {
	hist, err := GetWorkflowHistory(ctx, client, id, runID)
	if err != nil {
		return err
	}
	replayer := worker.NewWorkflowReplayer()
	replayer.RegisterWorkflow(Workflow)
	return replayer.ReplayWorkflowHistory(nil, hist)
}
