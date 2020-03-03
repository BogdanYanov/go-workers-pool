package warehouse

import (
	"github.com/BogdanYanov/go-workers-pool/work"
	_ "runtime"
	"sync"
)

// Worker is an interface that contains methods for implementation as a worker
type Worker interface {
	GetID() int
	Work()
	Stop()
	CountAmountWorkDone() int
}

// Employee is an employee abstraction that do unloading to warehouse.
type Employee struct {
	ID           int
	numWorkDone  int
	assignedWork chan work.Work
	workDone     *sync.WaitGroup
	stopped      *sync.WaitGroup
	quit         chan struct{}
}

// NewEmployee creates new employee.
func NewEmployee(id int, assignedWork chan work.Work, workDone *sync.WaitGroup, stopped *sync.WaitGroup) Worker {
	return &Employee{
		ID:           id,
		numWorkDone:  0,
		assignedWork: assignedWork,
		workDone:     workDone,
		stopped:      stopped,
		quit:         make(chan struct{}),
	}
}

// Work do incoming work from a warehouse work channel.
func (e *Employee) Work() {
	go func() {
		for {
			select {
			case job := <-e.assignedWork:
				job.Do()
				e.numWorkDone++
				e.workDone.Done()
				//runtime.Gosched()
			case <-e.quit:
				e.stopped.Done()
				return
			}
		}
	}()
}

// Stop stops the employee.
func (e *Employee) Stop() {
	e.quit <- struct{}{}
}

// GetID returns id of employee.
func (e *Employee) GetID() int {
	return e.ID
}

// CountAmountWorkDone return the amount of work done by the employee.
func (e *Employee) CountAmountWorkDone() int {
	return e.numWorkDone
}
