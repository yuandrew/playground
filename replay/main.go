package main

import (
	"context"
	"fmt"

	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/history/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	// "go.temporal.io/server/common/log"
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
		panic(fmt.Sprintf("Unable to create client: %v", err))
	}
	defer c.Close()

	// Replay workflow
	fmt.Println("Replaying workflow")
	// d7226966-0465-4e31-9fdd-36cc88c59446
	// 0f7858a6-0953-4be1-bac5-6af0f2e9a897
	getID := "hello_world_workflowID"
	runID := "a9199bec-e8e3-4a9e-8b1f-0735df071379"
	// workflow.SetQueryHandler
	queryType := "__stack_trace"
	fmt.Println("queryType: ", queryType)
	// queryType := "__open_sessions"
	// queryType := "__temporal_workflow_metadata"
	// var val string
	// err = ReplayWorkflow(context.Background(), c, getID, runID)
	// err = c.SignalWorkflow(context.Background(), getID, runID, "tick", nil)
	val, err := c.QueryWorkflow(context.Background(), getID, runID, queryType, nil)
	if err != nil {
		panic(fmt.Sprintf("aaaUnable to replay workflow: %v", err))
	}
	var result string
	err = val.Get(&result)
	fmt.Println("Query val: ", val)
	fmt.Println("Query result: ", result)
	if err != nil {
		panic(fmt.Sprintf("Unable to replay workflow: %v", err))
	}

	fmt.Println("done")
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
