package worker

import (
	"github.com/BogdanYanov/go-workers-pool/job"
	"sync"
)

// Dispatcher is the interface that wraps the basic dispatcher methods
type Dispatcher interface {
	Run()
	Stop()
}

// Warehouse is the abstraction of the warehouse and implementation of dispatcher object.
type Warehouse struct {
	workers           []Worker
	jobQueue          job.Queue
	workerPool        chan chan job.Job
	workerStopped     *sync.WaitGroup
	dispatcherStopped *sync.WaitGroup
	allStop           chan struct{}
}

// WarehouseInit initialize warehouse with some starting settings
func WarehouseInit(maxWorkers int, queue job.Queue) Dispatcher {
	if maxWorkers < 1 || queue == nil {
		return nil
	}
	w := &Warehouse{
		workers:           make([]Worker, maxWorkers, maxWorkers),
		jobQueue:          queue,
		workerPool:        make(chan chan job.Job, maxWorkers),
		workerStopped:     &sync.WaitGroup{},
		dispatcherStopped: &sync.WaitGroup{},
		allStop:           make(chan struct{}),
	}

	for i := 0; i < maxWorkers; i++ {
		w.workers[i] = NewEmployee(i+1, w.workerPool, w.workerStopped)
	}

	return w
}

// Run start the warehouse and check for job
func (w *Warehouse) Run() {
	for i := 0; i < len(w.workers); i++ {
		w.workerStopped.Add(1)
		go w.workers[i].Watch()
	}

	w.dispatcherStopped.Add(1)
	go w.dispatch()
}

func (w *Warehouse) dispatch() {
	for {
		select {
		case newJob := <-w.jobQueue:
			jobChannel := <-w.workerPool
			jobChannel <- newJob
		case <-w.allStop:
			for i := 0; i < len(w.workers); i++ {
				w.workers[i].Stop()
			}
			w.workerStopped.Wait()
			w.dispatcherStopped.Done()
			return
		}
	}
}

// Stop stops the warehouse job executing
func (w *Warehouse) Stop() {
	w.allStop <- struct{}{}
	w.dispatcherStopped.Wait()
}
