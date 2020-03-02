package warehouse

import (
	"github.com/BogdanYanov/go-workers-pool/work"
	"sync"
	"testing"
)

func TestNewUnloader(t *testing.T) {
	var (
		assignedWork = make(chan work.Work)
		workDone     = &sync.WaitGroup{}
		stopped      = &sync.WaitGroup{}
	)

	u := NewUnloader(0, assignedWork, workDone, stopped)
	if u == nil {
		t.Errorf("NewUnloader(): Unloader is nil")
	}
}

func TestUnloader_Work(t *testing.T) {
	var (
		productsNum  = 100
		truck        = work.NewTruck(int32(productsNum))
		assignedWork = make(chan work.Work)
		workDone     = &sync.WaitGroup{}
		stopped      = &sync.WaitGroup{}
	)

	workDone.Add(int(truck.AvailableWork()))
	go func() {
		for i := 0; i < productsNum; i++ {
			assignedWork <- truck
		}
	}()

	type fields struct {
		ID           int
		assignedWork chan work.Work
		workDone     *sync.WaitGroup
		stopped      *sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Work() case 1",
			fields: fields{
				ID:           0,
				assignedWork: assignedWork,
				workDone:     workDone,
				stopped:      stopped,
			},
			want: productsNum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUnloader(tt.fields.ID, tt.fields.assignedWork, tt.fields.workDone, tt.fields.stopped)
			stopped.Add(1)
			u.Work()
			workDone.Wait()
			u.Stop()
			stopped.Wait()
			if got := u.CountAmountWorkDone(); got != tt.want {
				t.Errorf("CountAmountWorkDone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnloader_GetID(t *testing.T) {
	type fields struct {
		ID int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "GetID() case 1",
			fields: fields{
				ID: 0,
			},
			want: 0,
		},
		{
			name: "GetID() case 2",
			fields: fields{
				ID: -1,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUnloader(tt.fields.ID, nil, nil, nil)
			if got := u.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}
