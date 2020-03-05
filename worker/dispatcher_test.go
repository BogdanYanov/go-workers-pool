package worker

import (
	"github.com/BogdanYanov/go-workers-pool/job"
	"sync"
	"testing"
)

func TestEmployee_GetID(t *testing.T) {
	type fields struct {
		id int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "GetID() case 1",
			fields: fields{id: 0},
			want:   0,
		},
		{
			name:   "GetID() case 2",
			fields: fields{id: 48},
			want:   48,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEmployee(tt.fields.id, nil, nil)
			if got := e.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmployee_Scenario(t *testing.T) {
	var (
		workerPool = make(chan chan job.Job, 1)
		truck      = job.NewTruck(2)
		testJob    = job.NewJob(truck)
		stopped    = &sync.WaitGroup{}
	)

	type fields struct {
		id         int
		workerPool chan chan job.Job
		stopped    *sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Stop() case 1",
			fields: fields{
				id:         0,
				workerPool: workerPool,
				stopped:    stopped,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEmployee(tt.fields.id, tt.fields.workerPool, tt.fields.stopped)
			stopped.Add(1)
			go e.Watch()
			go func() {
				jobCh := <-workerPool
				jobCh <- testJob
				e.Stop()
			}()
			stopped.Wait()
			if truck.AvailableWork() != 1 {
				t.Errorf("Stop(): worker don`t stop. (available work = %d)", truck.AvailableWork())
			}
		})
	}
}

func TestWarehouseInit(t *testing.T) {
	var dispatcher Dispatcher

	dispatcher = WarehouseInit(0, nil)
	if dispatcher != nil {
		t.Errorf("WarehouseInit(): must return nil")
	}
	dispatcher = WarehouseInit(1, nil)
	if dispatcher != nil {
		t.Errorf("WarehouseInit(): must return nil")
	}
	dispatcher = WarehouseInit(1, job.QueueInit(0))
	if dispatcher == nil {
		t.Errorf("WarehouseInit(): must return non-nil value")
	}
}

func TestWarehouse_Scenario(t *testing.T) {
	var (
		jobQueue   job.Queue
		dispatcher Dispatcher
		testTrucks = []*job.Truck{
			job.NewTruck(100000),
			job.NewTruck(100000),
		}
		warehouseJob job.Job
	)

	jobQueue = job.QueueInit(0)

	dispatcher = WarehouseInit(100, jobQueue)

	dispatcher.Run()

	for i := 0; i < len(testTrucks); i++ {
		warehouseJob = job.NewJob(testTrucks[i])
		availableWork := testTrucks[i].AvailableWork()
		for i := 0; i < int(availableWork); i++ {
			jobQueue.SendJob(warehouseJob)
		}
	}

	dispatcher.Stop()

	for i := 0; i < len(testTrucks); i++ {
		if testTrucks[i].AvailableWork() != 0 {
			t.Errorf("Scenario failed. Don`t do all work")
		}
	}
}

func BenchmarkWarehouse_Run(b *testing.B) {
	var productsNum int32 = 100000

	type args struct {
		workersNum  int
		productsNum int32
	}
	var tests = []struct {
		name string
		args args
	}{
		{
			name: "Benchmark SendWork() case 1 - 1 worker",
			args: args{
				workersNum:  1,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark SendWork() case 2 - 2 workers",
			args: args{
				workersNum:  2,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark SendWork() case 3 - 10 workers",
			args: args{
				workersNum:  10,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark SendWork() case 4 - 50 workers",
			args: args{
				workersNum:  50,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark SendWork() case 5 - 100 workers",
			args: args{
				workersNum:  100,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark SendWork() case 6 - 500 workers",
			args: args{
				workersNum:  500,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark SendWork() case 7 - 1000 workers",
			args: args{
				workersNum:  1000,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark SendWork() case 8 - 5000 workers",
			args: args{
				workersNum:  5000,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark SendWork() case 9 - 10000 workers",
			args: args{
				workersNum:  10000,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark SendWork() case 10 - 50000 workers",
			args: args{
				workersNum:  50000,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark SendWork() case 11 - 100000 workers",
			args: args{
				workersNum:  100000,
				productsNum: productsNum,
			},
		},
	}

	for _, tt := range tests {
		queue := job.QueueInit(0)
		dispatcher := WarehouseInit(tt.args.workersNum, queue)
		dispatcher.Run()
		b.Run(tt.name, func(b *testing.B) {
			testJob := job.NewJob(job.NewTruck(tt.args.productsNum))
			for i := 0; i < int(tt.args.productsNum); i++ {
				queue.SendJob(testJob)
			}
		})
		dispatcher.Stop()
	}
}
