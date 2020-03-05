package job

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewJob(t *testing.T) {
	var (
		truck             *Truck
		testJob           Job
		wantAvailableWork int32
	)

	truck = NewTruck(10)
	wantAvailableWork = truck.AvailableWork() - 1

	testJob = NewJob(truck)

	if testJob.Operation == nil {
		t.Errorf("NewJob(): Operation is nil")
	}

	testJob.Operation.Execute()

	if truck.AvailableWork() != wantAvailableWork {
		t.Errorf("NewJob(): Execute() = %d, want - %d", truck.AvailableWork(), wantAvailableWork)
	}
}

func TestNewTruck(t *testing.T) {
	truck := NewTruck(0)
	if truck != nil {
		t.Errorf("NewTruck() must be nil")
	}

	truck = NewTruck(160)
	if truck == nil {
		t.Errorf("NewTruck() must be not nil")
	}
}

func TestTruck_AvailableWork(t1 *testing.T) {
	type fields struct {
		productsNum int32
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		{
			name:   "AvailableWork() case 1",
			fields: fields{productsNum: 200},
			want:   199,
		},
		{
			name:   "AvailableWork() case 2",
			fields: fields{productsNum: 1},
			want:   0,
		},
		{
			name:   "AvailableWork() case 3",
			fields: fields{productsNum: 0},
			want:   0,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Truck{
				productsNum: tt.fields.productsNum,
			}
			t.unload()
			if got := t.AvailableWork(); got != tt.want {
				t1.Errorf("AvailableWork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruck_Execute(t1 *testing.T) {
	type fields struct {
		productsNum int32
	}
	tests := []struct {
		name       string
		fields     fields
		workersNum int
		want       int
	}{
		{
			name: "Execute() case 1",
			fields: fields{
				productsNum: 100,
			},
			workersNum: 5,
			want:       95,
		},
		{
			name: "Execute() case 2",
			fields: fields{
				productsNum: 5,
			},
			workersNum: 5,
			want:       0,
		},
		{
			name: "Execute() case 3",
			fields: fields{
				productsNum: 4,
			},
			workersNum: 5,
			want:       0,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Truck{
				productsNum: tt.fields.productsNum,
			}
			wg := &sync.WaitGroup{}
			wg.Add(tt.workersNum)
			for i := 0; i < tt.workersNum; i++ {
				go func(wg *sync.WaitGroup) {
					t.Execute()
					wg.Done()
				}(wg)
			}
			wg.Wait()
			if int(t.productsNum) != tt.want {
				t1.Errorf("Do(): productsNum = %d, want - %d", t.productsNum, tt.want)
			}
		})
	}
}

func TestQueueInit(t *testing.T) {
	type args struct {
		maxJobs int
	}
	tests := []struct {
		name         string
		args         args
		wantCapacity int
	}{
		{
			name: "QueueInit() case 1",
			args: args{
				maxJobs: 0,
			},
			wantCapacity: 0,
		},
		{
			name: "QueueInit() case 2",
			args: args{
				maxJobs: 1,
			},
			wantCapacity: 1,
		},
		{
			name: "QueueInit() case 3",
			args: args{
				maxJobs: 50,
			},
			wantCapacity: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := QueueInit(tt.args.maxJobs)
			if got := cap(queue); got != tt.wantCapacity {
				t.Errorf("QueueInit(): capacity = %d, want - %d", got, tt.wantCapacity)
			}
		})
	}
}

func TestQueue_SendJob(t *testing.T) {
	var (
		testQueue = QueueInit(0)
		testJob1  = NewJob(NewTruck(10))
		testJob2  = NewJob(NewTruck(20))
	)

	type args struct {
		job Job
	}
	tests := []struct {
		name string
		q    Queue
		args args
		want Job
	}{
		{
			name: "SendJob() case 1",
			q:    testQueue,
			args: args{
				job: testJob1,
			},
			want: testJob1,
		},
		{
			name: "SendJob() case 2",
			q:    testQueue,
			args: args{
				job: testJob2,
			},
			want: testJob2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Job

			wg := &sync.WaitGroup{}
			wg.Add(1)

			go func() {
				got = <-tt.q
				wg.Done()
			}()

			tt.q.SendJob(tt.args.job)

			wg.Wait()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendJob(): got = %v, want - %v", got, tt.want)
			}
		})
	}
}
