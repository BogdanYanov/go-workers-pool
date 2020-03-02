package warehouse

import (
	"fmt"
	"github.com/BogdanYanov/go-workers-pool/work"
	"sync"
)

// Warehouse is an abstraction of the warehouse into which products are unloaded.
type Warehouse struct {
	workers      []Worker        // pool of goroutines
	assignedWork chan work.Work  // the channel to which all work is sent
	workDone     *sync.WaitGroup // WaitGroup to wait until all work is done
	stopped      *sync.WaitGroup // WaitGroup to wait until all goroutines are completed
	isStopped    bool            // boolean value of whether the warehouse is working
	nowWork      int             // working goroutines counter
}

// NewWarehouse creates new warehouse.
func NewWarehouse() *Warehouse {
	return &Warehouse{
		workers:      nil,
		assignedWork: make(chan work.Work),
		workDone:     &sync.WaitGroup{},
		stopped:      &sync.WaitGroup{},
		isStopped:    true,
		nowWork:      0,
	}
}

// Start launches unloaders work.
func (w *Warehouse) Start(workersNum int) {
	// if workersNum is less than 1, don`t start the work
	if workersNum < 1 {
		return
	}

	// if the warehouse is already works, don`t start again
	if !w.isStopped {
		return
	}

	// if the warehouse does not have workers, create them
	if w.workers == nil {
		w.workers = make([]Worker, workersNum, workersNum)
		for i := 0; i < workersNum; i++ {
			w.workers[i] = NewUnloader(i, w.assignedWork, w.workDone, w.stopped)
		}
	}

	// add all working goroutines to WaitGroup and set the value for nowWork counter
	w.stopped.Add(workersNum)
	w.nowWork = workersNum
	for i := 0; i < workersNum; i++ {
		w.workers[i].Work()
	}
	w.isStopped = false
}

// ChangeNumWorkers change number of working unloaders.
func (w *Warehouse) ChangeNumWorkers(newNumWorkers int) {
	if w.isStopped {
		return
	}

	if len(w.workers) > newNumWorkers {
		diff := len(w.workers) - newNumWorkers
		for i := 1; i <= diff; i++ {
			w.workers[len(w.workers)-i].Stop()
		}
		w.nowWork = w.nowWork - diff
	}

	if len(w.workers) < newNumWorkers {
		diff := newNumWorkers - len(w.workers)
		for i := 0; i < diff; i++ {
			w.workers = append(w.workers, NewUnloader(len(w.workers), w.assignedWork, w.workDone, w.stopped))
			w.stopped.Add(1)
			w.workers[len(w.workers)-1].Work()
		}
		w.nowWork = w.nowWork + diff
	}
}

// SendWork send work for execution.
func (w *Warehouse) SendWork(newWork work.Work) {
	availableWork := int(newWork.AvailableWork())

	w.workDone.Add(availableWork)
	go func() {
		for i := 0; i < availableWork; i++ {
			w.assignedWork <- newWork
		}
	}()

	w.workDone.Wait()
}

// Stop stops all working unloaders.
func (w *Warehouse) Stop() {
	if w.isStopped {
		return
	}
	for i := 0; i < w.nowWork; i++ {
		w.workers[i].Stop()
	}
	w.stopped.Wait()
	w.nowWork = 0
	w.isStopped = true
}

// WorkersInfo displays information about all unloaders in warehouse.
func (w *Warehouse) WorkersInfo() {
	for i := 0; i < len(w.workers); i++ {
		fmt.Printf("Worker #%d; total products unload - %d\n", w.workers[i].GetID(), w.workers[i].CountAmountWorkDone())
	}
}
