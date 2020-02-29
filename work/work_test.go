package work

import (
	"sync"
	"testing"
)

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

func TestTruck_Do(t1 *testing.T) {
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
			name: "Do() case 1",
			fields: fields{
				productsNum: 100,
			},
			workersNum: 5,
			want:       95,
		},
		{
			name: "Do() case 2",
			fields: fields{
				productsNum: 5,
			},
			workersNum: 5,
			want:       0,
		},
		{
			name: "Do() case 3",
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
					t.Do()
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
