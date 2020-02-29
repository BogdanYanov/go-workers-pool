package warehouse

import (
	"github.com/BogdanYanov/go-workers-pool/work"
	"sync"
)

const workQueueCapacity = 10

// Warehouse abstraction of the warehouse into which products are unloaded.
type Warehouse struct {
	workers            []Worker
	workQueue          []work.Work
	workChannel        chan work.Work
	productsNumInStock int
	workDone           *sync.WaitGroup
}

// NewWarehouse creates new warehouse.
func NewWarehouse() *Warehouse {
	return &Warehouse{
		workDone:    &sync.WaitGroup{},
		workChannel: make(chan work.Work),
	}
}

// AddWork add work for execution.
func (wh *Warehouse) AddWork(newWork work.Work) {
	if wh.workQueue == nil {
		wh.workQueue = make([]work.Work, 0, workQueueCapacity)
	}
	wh.workQueue = append(wh.workQueue, newWork)
}

// StartWork launches goroutines and performs work from the work queue array.
func (wh *Warehouse) StartWork(workersNum int) {
	if wh.workQueue == nil || len(wh.workQueue) == 0 {
		return
	}

	//wh.workers = make([]Worker, 0, workersNum)
	wh.workers = make([]Worker, workersNum, workersNum)
	wh.workChannel = make(chan work.Work)

	var availableWork int32

	for i := 0; i < len(wh.workQueue); i++ {
		availableWork += wh.workQueue[i].AvailableWork()
	}

	for i := 0; i < cap(wh.workers); i++ {
		//wh.workers = append(wh.workers, NewUnloader(i, wh.workChannel, wh.wg))
		wh.workers[i] = NewUnloader(i, wh.workChannel, wh.workDone)
	}

	wh.workDone.Add(int(availableWork))
	for i := 0; i < len(wh.workers); i++ {
		go wh.workers[i].Work()
	}

	for i := 0; i < len(wh.workQueue); {
		currentAvailableWork := wh.workQueue[i].AvailableWork()
		for j := 0; j < int(currentAvailableWork); j++ {
			wh.workChannel <- wh.workQueue[i]
		}
		wh.workQueue = wh.workQueue[1:]
	}

	close(wh.workChannel)
	wh.workDone.Wait()

	for i := 0; i < len(wh.workers); i++ {
		wh.productsNumInStock += wh.workers[i].CountAmountWorkDone()
	}
}

// ProductsInStock return the number of products in the warehouse.
func (wh *Warehouse) ProductsInStock() int {
	return wh.productsNumInStock
}
