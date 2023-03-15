package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/helloworld"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: "192.168.8.42:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "hello-world", worker.Options{})

	//w.RegisterWorkflow(helloworld.Workflow)
	//w.RegisterActivity(helloworld.Activity)

	w.RegisterWorkflow(helloworld.MyWorkflow)
	w.RegisterActivity(helloworld.TaskA)
	w.RegisterActivity(helloworld.TaskB)
	w.RegisterActivity(helloworld.TaskC)
	w.RegisterActivity(helloworld.TaskD)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
