package warehouse

import (
	"github.com/BogdanYanov/go-workers-pool/work"
	"reflect"
	"sync"
	"testing"
)

func TestNewUnloader(t *testing.T) {
	var (
		assignedWork = make(chan work.Work)
		workDone     = &sync.WaitGroup{}
	)

	type args struct {
		ID           int
		assignedWork chan work.Work
		workDone     *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
		want Worker
	}{
		{
			name: "NewUnloader() case 1",
			args: args{
				ID:           0,
				assignedWork: assignedWork,
				workDone:     workDone,
			},
			want: &Unloader{
				ID:                  0,
				assignedWork:        assignedWork,
				numProductsUnloaded: 0,
				workDone:            workDone,
			},
		},
		{
			name: "NewUnloader() case 2",
			args: args{
				ID:           1,
				assignedWork: assignedWork,
				workDone:     workDone,
			},
			want: &Unloader{
				ID:                  1,
				assignedWork:        assignedWork,
				numProductsUnloaded: 0,
				workDone:            workDone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnloader(tt.args.ID, tt.args.assignedWork, tt.args.workDone); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnloader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnloader_Work(t *testing.T) {
	var (
		productsNum  = 4
		truck        = work.NewTruck(int32(productsNum))
		assignedWork = make(chan work.Work, productsNum)
		workDone     = &sync.WaitGroup{}
	)

	for i := 0; i < productsNum; i++ {
		assignedWork <- truck
	}

	type fields struct {
		ID                  int
		assignedWork        chan work.Work
		numProductsUnloaded int
		workDone            *sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Work() case 1",
			fields: fields{
				ID:                  0,
				assignedWork:        assignedWork,
				numProductsUnloaded: 0,
				workDone:            workDone,
			},
			want: productsNum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Unloader{
				ID:                  tt.fields.ID,
				assignedWork:        tt.fields.assignedWork,
				numProductsUnloaded: tt.fields.numProductsUnloaded,
				workDone:            tt.fields.workDone,
			}
			workDone.Add(productsNum)
			go u.Work()
			workDone.Wait()
			if got := u.CountAmountWorkDone(); got != tt.want {
				t.Errorf("CountAmountWorkDone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnloader_GetID(t *testing.T) {
	type fields struct {
		ID                  int
		assignedWork        chan work.Work
		numProductsUnloaded int
		workDone            *sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "GetID() case 1",
			fields: fields{
				ID:                  0,
				assignedWork:        nil,
				numProductsUnloaded: 0,
				workDone:            nil,
			},
			want: 0,
		},
		{
			name: "GetID() case 2",
			fields: fields{
				ID:                  -1,
				assignedWork:        nil,
				numProductsUnloaded: 0,
				workDone:            nil,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Unloader{
				ID:                  tt.fields.ID,
				assignedWork:        tt.fields.assignedWork,
				numProductsUnloaded: tt.fields.numProductsUnloaded,
				workDone:            tt.fields.workDone,
			}
			if got := u.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}
