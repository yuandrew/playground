package main

import (
	"log"

	"github.com/yuandrew/playground" // Replace with the actual path to the playground package

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "hello-world", worker.Options{})

	w.RegisterWorkflow(playground.Workflow)
	w.RegisterActivity(playground.Activity)
	w.RegisterActivity(playground.Activity1)
	w.RegisterActivity(playground.Activity2)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
