package main

import (
	"fmt"
	"github.com/BogdanYanov/go-workers-pool/job"
	"github.com/BogdanYanov/go-workers-pool/worker"
)

func main() {

	var (
		jobQueue job.Queue
		dispatcher worker.Dispatcher
		truck *job.Truck
		warehouseJob job.Job
	)

	jobQueue = job.QueueInit(0)

	dispatcher = worker.WarehouseInit(100, jobQueue)

	dispatcher.Run()

	truck = job.NewTruck(100000)

	warehouseJob = job.NewJob(truck)

	availableWork := truck.AvailableWork()

	for i := 0; i < int(availableWork); i++ {
		jobQueue.SendJob(warehouseJob)
	}

	dispatcher.Stop()

	fmt.Println(truck.AvailableWork())
}
