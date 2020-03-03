package warehouse

import (
	"github.com/BogdanYanov/go-workers-pool/work"
	"testing"
)

func TestNewWarehouse(t *testing.T) {
	wh := NewWarehouse()
	if wh == nil {
		t.Errorf("NewWarehouse() must return non-nil value")
	}
}

func TestWarehouse_SendWork(t *testing.T) {
	var (
		truck1 = work.NewTruck(100)
		truck2 = work.NewTruck(250)
	)

	type args struct {
		newWork    work.Work
		workersNum int
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "SendWork() case 1",
			args: args{
				newWork:    truck1,
				workersNum: 2,
			},
			want: 0,
		},
		{
			name: "SendWork() case 2",
			args: args{
				newWork:    truck2,
				workersNum: 4,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wh := NewWarehouse()
			wh.Start(tt.args.workersNum)

			wh.SendWork(tt.args.newWork)

			wh.Stop()

			if got := tt.args.newWork.AvailableWork(); int(got) != tt.want {
				t.Errorf("SendWork(): got = %v, want - %v", got, tt.want)
			}
		})
	}
}

func TestWarehouse_Start(t *testing.T) {
	type args struct {
		workersNum int
	}
	var tests = []struct {
		name               string
		args               args
		do                 func(w *Warehouse)
		wantNumWorkWorkers int
		wantIsStopped      bool
		wantNumWorkers     int
	}{
		{
			name: "Start() case 1",
			args: args{
				workersNum: 0,
			},
			do:                 func(w *Warehouse) {},
			wantNumWorkWorkers: 0,
			wantIsStopped:      true,
			wantNumWorkers:     0,
		},
		{
			name: "Start() case 2",
			args: args{
				workersNum: 2,
			},
			do:                 func(w *Warehouse) { w.isStopped = false },
			wantNumWorkWorkers: 0,
			wantIsStopped:      false,
			wantNumWorkers:     0,
		},
		{
			name: "Start() case 3",
			args: args{
				workersNum: 3,
			},
			do:                 func(w *Warehouse) {},
			wantNumWorkWorkers: 3,
			wantIsStopped:      false,
			wantNumWorkers:     3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wh := NewWarehouse()
			tt.do(wh)
			wh.Start(tt.args.workersNum)
			if wh.nowWork != tt.wantNumWorkWorkers {
				t.Errorf("Start(): nowWork: got = %d, want - %d", wh.nowWork, tt.wantNumWorkWorkers)
			}
			if wh.isStopped != tt.wantIsStopped {
				t.Errorf("Start(): isStopped: got = %v, want - %v", wh.isStopped, tt.wantIsStopped)
			}
			if len(wh.workers) != tt.wantNumWorkers {
				t.Errorf("Start(): length of workers: got = %d, want- %d", len(wh.workers), tt.wantNumWorkers)
			}
		})
	}
}

func TestWarehouse_Stop(t *testing.T) {
	type args struct {
		workersNum int
	}
	var tests = []struct {
		name          string
		args          args
		do            func(w *Warehouse, workersNum int)
		wantIsStopped bool
		wantNowWorks  int
	}{
		{
			name: "Stop() case 1",
			args: args{workersNum: 2},
			do: func(w *Warehouse, workersNum int) {
				w.Start(workersNum)
			},
			wantIsStopped: true,
			wantNowWorks:  0,
		},
		{
			name: "Stop() case 2",
			args: args{workersNum: 0},
			do: func(w *Warehouse, workersNum int) {
				w.Start(workersNum)
			},
			wantIsStopped: true,
			wantNowWorks:  0,
		},
		{
			name:          "Stop() case 3",
			args:          args{workersNum: 2},
			do:            func(w *Warehouse, workersNum int) {},
			wantIsStopped: true,
			wantNowWorks:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wh := NewWarehouse()
			tt.do(wh, tt.args.workersNum)
			wh.Stop()
			if wh.isStopped != tt.wantIsStopped {
				t.Errorf("Stop(): isStopped: got = %v, want - %v", wh.isStopped, tt.wantIsStopped)
			}
			if wh.nowWork != tt.wantNowWorks {
				t.Errorf("Stop(): nowWorks: got = %v, want - %v", wh.nowWork, tt.wantNowWorks)
			}
		})
	}
}

func BenchmarkWarehouse_StartWork(b *testing.B) {
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
			name: "Benchmark StartWork() case 1 - 1 worker",
			args: args{
				workersNum:  1,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark StartWork() case 2 - 2 workers",
			args: args{
				workersNum:  2,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark StartWork() case 3 - 10 workers",
			args: args{
				workersNum:  10,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark StartWork() case 4 - 50 workers",
			args: args{
				workersNum:  50,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark StartWork() case 5 - 100 workers",
			args: args{
				workersNum:  100,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark StartWork() case 6 - 500 workers",
			args: args{
				workersNum:  500,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark StartWork() case 7 - 1000 workers",
			args: args{
				workersNum:  1000,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark StartWork() case 8 - 5000 workers",
			args: args{
				workersNum:  5000,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark StartWork() case 9 - 10000 workers",
			args: args{
				workersNum:  10000,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark StartWork() case 10 - 50000 workers",
			args: args{
				workersNum:  50000,
				productsNum: productsNum,
			},
		},
		{
			name: "Benchmark StartWork() case 11 - 100000 workers",
			args: args{
				workersNum:  100000,
				productsNum: productsNum,
			},
		},
	}

	for _, tt := range tests {
		wh := NewWarehouse()
		wh.Start(tt.args.workersNum)
		b.Run(tt.name, func(b *testing.B) {
			wh.SendWork(work.NewTruck(tt.args.productsNum))
		})
		wh.Stop()
	}

}

func TestWarehouse_ChangeNumWorkers(t *testing.T) {
	type args struct {
		workersNum    int
		newNumWorkers int
	}
	tests := []struct {
		name         string
		args         args
		do           func(w *Warehouse)
		wantNowWorks int
	}{
		{
			name: "ChangeNumWorkers() case 1",
			args: args{
				workersNum:    2,
				newNumWorkers: 4,
			},
			do: func(w *Warehouse) {
				w.Stop()
			},
			wantNowWorks: 0,
		},
		{
			name: "ChangeNumWorkers() case 2",
			args: args{
				workersNum:    4,
				newNumWorkers: 2,
			},
			do:           func(w *Warehouse) {},
			wantNowWorks: 2,
		},
		{
			name: "ChangeNumWorkers() case 3",
			args: args{
				workersNum:    2,
				newNumWorkers: 4,
			},
			do:           func(w *Warehouse) {},
			wantNowWorks: 4,
		},
		{
			name: "ChangeNumWorkers() case 4",
			args: args{
				workersNum:    2,
				newNumWorkers: 0,
			},
			do:           func(w *Warehouse) {},
			wantNowWorks: 2,
		},
		{
			name: "ChangeNumWorkers() case 5",
			args: args{
				workersNum:    10,
				newNumWorkers: 12,
			},
			do: func(w *Warehouse) {
				w.ChangeNumWorkers(6)
			},
			wantNowWorks: 12,
		},
		{
			name: "ChangeNumWorkers() case 5",
			args: args{
				workersNum:    10,
				newNumWorkers: 8,
			},
			do: func(w *Warehouse) {
				w.ChangeNumWorkers(6)
			},
			wantNowWorks: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWarehouse()
			w.Start(tt.args.workersNum)
			defer w.Stop()
			tt.do(w)
			w.ChangeNumWorkers(tt.args.newNumWorkers)
			if w.nowWork != tt.wantNowWorks {
				t.Errorf("ChangeNumWorkers(): nowWork = %d, want - %d", w.nowWork, tt.wantNowWorks)
			}
		})
	}
}

func TestWarehouse_ProductsInStock(t *testing.T) {
	w := NewWarehouse()

	w.Start(2)

	trucks := []work.Work{
		work.NewTruck(100000),
		work.NewTruck(20000),
	}

	var totalProducts int32

	for i := 0; i < len(trucks); i++ {
		totalProducts += trucks[i].AvailableWork()
		w.SendWork(trucks[i])
	}

	w.Stop()

	if w.ProductsInStock() != int(totalProducts) {
		t.Errorf("ProductsInStock(): got = %d, want - %d", w.ProductsInStock(), totalProducts)
	}
}
