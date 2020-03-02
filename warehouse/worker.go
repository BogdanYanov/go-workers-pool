package warehouse

import (
	"github.com/BogdanYanov/go-workers-pool/work"
	"runtime"
	"sync"
)

// Worker is an interface that contains methods for implementation as a worker
type Worker interface {
	GetID() int
	Work()
	Stop()
	CountAmountWorkDone() int
}

// Unloader is an unloader abstraction that do unloading to warehouse.
type Unloader struct {
	ID                  int
	numProductsUnloaded int
	assignedWork        chan work.Work
	workDone            *sync.WaitGroup
	stopped             *sync.WaitGroup
	quit                chan struct{}
}

// NewUnloader creates new unloader.
func NewUnloader(id int, assignedWork chan work.Work, workDone *sync.WaitGroup, stopped *sync.WaitGroup) Worker {
	return &Unloader{
		ID:                  id,
		numProductsUnloaded: 0,
		assignedWork:        assignedWork,
		workDone:            workDone,
		stopped:             stopped,
		quit:                make(chan struct{}),
	}
}

// Work do incoming work from a warehouse work channel.
func (u *Unloader) Work() {
	go func() {
		for {
			select {
			case job := <-u.assignedWork:
				job.Do()
				u.numProductsUnloaded++
				u.workDone.Done()
				runtime.Gosched()
			case <-u.quit:
				u.stopped.Done()
				return
			}
		}
	}()
}

// Stop stops the unloader.
func (u *Unloader) Stop() {
	u.quit <- struct{}{}
}

// GetID returns id of unloader.
func (u *Unloader) GetID() int {
	return u.ID
}

// CountAmountWorkDone return the amount of work done by the unloader.
func (u *Unloader) CountAmountWorkDone() int {
	return u.numProductsUnloaded
}
