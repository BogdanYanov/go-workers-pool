package main

import (
	"fmt"
	"github.com/BogdanYanov/go-workers-pool/warehouse"
	"github.com/BogdanYanov/go-workers-pool/work"
)

func main() {
	var trucks = []work.Work{
		work.NewTruck(100000),
		work.NewTruck(50000),
		work.NewTruck(200000),
		work.NewTruck(20000),
	}

	wh := warehouse.NewWarehouse()

	for i := 0; i < len(trucks); i++ {
		wh.AddWork(trucks[i])
	}

	wh.StartWork(4)

	for i := 0; i < len(trucks); i++ {
		fmt.Printf("Truck #%d: available products - %d\n", i, trucks[i].AvailableWork())
	}

	fmt.Println("Products in warehouse -", wh.ProductsInStock())
}
