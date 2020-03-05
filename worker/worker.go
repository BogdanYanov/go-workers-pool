package worker

import (
	"github.com/BogdanYanov/go-workers-pool/job"
	"runtime"
	"sync"
)

// Worker is the interface that wraps basic Worker methods.
type Worker interface {
	Watch()
	Stop()
	GetID() int
}

// Employee is the employee abstraction that do unloading to warehouse and implementing Worker interface.
type Employee struct {
	id         int
	jobChannel chan job.Job
	workerPool chan chan job.Job
	stopped    *sync.WaitGroup
	quit       chan struct{}
}

// NewEmployee creates new employee.
func NewEmployee(id int, workerPool chan chan job.Job, stopped *sync.WaitGroup) Worker {
	return &Employee{
		id:         id,
		jobChannel: make(chan job.Job),
		workerPool: workerPool,
		stopped:    stopped,
		quit:       make(chan struct{}),
	}
}

// Watch start checking channel for a job or for an exit from watching.
func (e *Employee) Watch() {
	for {
		e.workerPool <- e.jobChannel
		select {
		case j := <-e.jobChannel:
			j.Operation.Execute()
			runtime.Gosched()
		case <-e.quit:
			e.stopped.Done()
			return
		}
	}
}

// Stop stops the employee.
func (e *Employee) Stop() {
	go func() {
		e.quit <- struct{}{}
	}()
}

// GetID returns id of employee.
func (e *Employee) GetID() int {
	return e.id
}
