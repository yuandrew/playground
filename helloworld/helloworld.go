package main

import (
	"time"

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
	// err := workflow.ExecuteActivity(ctx, Activity, name).Get(ctx, &result)
	// if err != nil {
	// 	logger.Error("Activity failed.", "Error", err)
	// 	return "", err
	// }
	// time.Sleep(1 * time.Second)

	logger.Info("HelloWorld workflow completed.", "result", result)

	return result, nil
}

//	func Activity(ctx context.Context, name string) (string, error) {
//		logger := activity.GetLogger(ctx)
//		logger.Info("Activityasdf", "name", name)
//		return "Hello " + name + "!", nil
//		// panic("FAIL LOCAL ACTIVITY")
//	}
