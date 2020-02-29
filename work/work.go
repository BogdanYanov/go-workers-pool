package work

import (
	"sync/atomic"
)

// Work is an interface that contains methods for implementing structures as work.
type Work interface {
	Do()
	AvailableWork() int32
}

// Truck abstraction of product delivery to warehouse.
type Truck struct {
	productsNum int32
}

// NewTruck create new truck as work.
func NewTruck(productsNum int32) Work {
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

// Do performs unloading of a truck for one product.
func (t *Truck) Do() {
	t.unload()
}
