package job

import (
	"sync/atomic"
)

// Executor is the interface that wraps basic methods for use it in job struct
type Executor interface {
	Execute()
}

// Job contain the Executor implementing structure
type Job struct {
	Operation Executor
}

// NewJob creates new job
func NewJob(job Executor) Job {
	return Job{Operation: job}
}

// Truck is the abstraction of product delivery to warehouse.
type Truck struct {
	productsNum int32
}

// NewTruck create new truck.
func NewTruck(productsNum int32) *Truck {
	if productsNum < 1 {
		return nil
	}
	return &Truck{productsNum: productsNum}
}

// AvailableWork returns the number of accomplishment of the necessary work.
func (t *Truck) AvailableWork() int32 {
	return t.productsNum
}

func (t *Truck) unload() {
	if t.productsNum != 0 {
		atomic.AddInt32(&(t.productsNum), -1)
	}
}

// Execute performs unloading of a truck for one product.
func (t *Truck) Execute() {
	t.unload()
}
