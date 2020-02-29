package warehouse

import (
	"github.com/BogdanYanov/go-workers-pool/work"
	"runtime"
	"sync"
)

var counter int32

// Worker is an interface that contains methods for implementation as a worker
type Worker interface {
	GetID() int
	Work()
	CountAmountWorkDone() int
}

// Unloader is an unloader abstraction that do unloading to warehouse.
type Unloader struct {
	ID                  int
	assignedWork        chan work.Work
	numProductsUnloaded int
	workDone            *sync.WaitGroup
}

// NewUnloader creates new unloader.
func NewUnloader(ID int, assignedWork chan work.Work, workDone *sync.WaitGroup) Worker {
	return &Unloader{
		ID:           ID,
		assignedWork: assignedWork,
		workDone:     workDone,
	}
}

// Work do incoming work from a warehouse work queue.
func (u *Unloader) Work() {
	for work := range u.assignedWork {
		work.Do()
		u.numProductsUnloaded++
		u.workDone.Done()
		runtime.Gosched()
	}
}

// CountAmountWorkDone return the amount of work done by the worker.
func (u *Unloader) CountAmountWorkDone() int {
	return u.numProductsUnloaded
}

// GetID returns worker id.
func (u *Unloader) GetID() int {
	return u.ID
}
