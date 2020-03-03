package main

import (
	"fmt"
	"github.com/BogdanYanov/go-workers-pool/warehouse"
	"github.com/BogdanYanov/go-workers-pool/work"
	"time"
)

func main() {
	wh := warehouse.NewWarehouse()

	truck1 := work.NewTruck(100000)

	wh.Start(10)
	wh.SendWork(truck1)

	fmt.Println(truck1.AvailableWork())
	wh.WorkersInfo()

	time.Sleep(time.Second)

	truck2 := work.NewTruck(100000)

	wh.ChangeNumWorkers(1)
	wh.SendWork(truck2)

	fmt.Println(truck2.AvailableWork())
	wh.WorkersInfo()

	time.Sleep(2 * time.Second)

	truck3 := work.NewTruck(100000)

	wh.ChangeNumWorkers(8)
	wh.SendWork(truck3)

	fmt.Println(truck3.AvailableWork())
	wh.WorkersInfo()

	truck4 := work.NewTruck(100000)
	wh.ChangeNumWorkers(12)
	wh.SendWork(truck4)

	fmt.Println(truck4.AvailableWork())
	wh.WorkersInfo()

	wh.Stop()

	wh.ProductsInStock()
}
