package warehouse

import (
	"github.com/BogdanYanov/go-workers-pool/work"
	"reflect"
	"sync"
	"testing"
)

func TestNewWarehouse(t *testing.T) {
	wh := NewWarehouse()
	if wh == nil {
		t.Errorf("NewWarehouse() must return non-nil value")
	}
}

func TestWarehouse_AddWork(t *testing.T) {
	var (
		truck1 = work.NewTruck(4)
		truck2 = work.NewTruck(5)
	)

	type fields struct {
		workers            []Worker
		workQueue          []work.Work
		workChannel        chan work.Work
		productsNumInStock int
		workDone           *sync.WaitGroup
	}
	type args struct {
		newWork []work.Work
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		pos    int
		want   work.Work
	}{
		{
			name: "AddWork() case 1",
			fields: fields{
				workers:            nil,
				workQueue:          nil,
				workChannel:        nil,
				productsNumInStock: 0,
				workDone:           nil,
			},
			args: args{
				newWork: []work.Work{truck1},
			},
			pos:  0,
			want: truck1,
		},
		{
			name: "AddWork() case 2",
			fields: fields{
				workers:            nil,
				workQueue:          nil,
				workChannel:        nil,
				productsNumInStock: 0,
				workDone:           nil,
			},
			args: args{
				newWork: []work.Work{truck2, truck1},
			},
			pos:  1,
			want: truck1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wh := &Warehouse{
				workers:            tt.fields.workers,
				workQueue:          tt.fields.workQueue,
				workChannel:        tt.fields.workChannel,
				productsNumInStock: tt.fields.productsNumInStock,
				workDone:           tt.fields.workDone,
			}

			for i := 0; i < len(tt.args.newWork); i++ {
				wh.AddWork(tt.args.newWork[i])
			}

			if got := wh.workQueue[tt.pos]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddWork(): got = %v, want - %v", got, tt.want)
			}
		})
	}
}

func TestWarehouse_StartWork(t *testing.T) {
	type fields struct {
		workQueue []work.Work
	}
	type args struct {
		workersNum int
	}
	var tests = []struct {
		name                string
		fields              fields
		args                args
		wantProductsInStock int
		wantWorkDoneZero    bool
	}{
		{
			name: "StartWork() case 1",
			fields: fields{
				workQueue: []work.Work{work.NewTruck(50)},
			},
			args: args{
				workersNum: 5,
			},
			wantProductsInStock: 50,
			wantWorkDoneZero:    false,
		},
		{
			name: "StartWork() case 2",
			fields: fields{
				workQueue: []work.Work{work.NewTruck(50)},
			},
			args: args{
				workersNum: 10,
			},
			wantProductsInStock: 50,
			wantWorkDoneZero:    false,
		},
		{
			name: "StartWork() case 3",
			fields: fields{
				workQueue: []work.Work{work.NewTruck(50)},
			},
			args: args{
				workersNum: 50,
			},
			wantProductsInStock: 50,
			wantWorkDoneZero:    true,
		},
		{
			name: "StartWork() case 4",
			fields: fields{
				workQueue: nil,
			},
			args: args{
				workersNum: 10,
			},
			wantProductsInStock: 0,
			wantWorkDoneZero:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wh := NewWarehouse()
			wh.workQueue = tt.fields.workQueue
			wh.StartWork(tt.args.workersNum)
			if wh.ProductsInStock() != tt.wantProductsInStock {
				t.Errorf("StartWork() -> ProductsInStock(): got = %d, want - %d", wh.ProductsInStock(), tt.wantProductsInStock)
			}
			for i := 0; i < len(wh.workers); i++ {
				if wh.workers[i].CountAmountWorkDone() == 0 && !tt.wantWorkDoneZero {
					t.Errorf("StartWork() -> CountAmountWorkDone(): got = %d", wh.workers[i].CountAmountWorkDone())
				}
			}
		})
	}
}

func BenchmarkWarehouse_StartWork(b *testing.B) {
	var productsNum int32 = 1000000

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
		b.ResetTimer()
		b.Run(tt.name, func(b *testing.B) {
			b.StopTimer()
			wh := NewWarehouse()
			wh.AddWork(work.NewTruck(tt.args.productsNum))
			b.StartTimer()
			wh.StartWork(tt.args.workersNum)
		})
	}

}
