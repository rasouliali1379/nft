package lambda

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	{
		in := []int{1, 2, 3, 4, 5, 6, 7, 8}
		want := []int{2, 4, 6, 8, 10, 12, 14, 16}
		if got := Map(in, func(t int) int { return t * 2 }); !reflect.DeepEqual(got, want) {
			t.Errorf("Filter() = %v, want %v", got, want)
		}
	}
	{
		in := []int{1, 2, 3, 4, 5, 6, 7, 8}
		want := []bool{false, true, false, true, false, true, false, true}
		if got := Map(in, func(t int) bool { return t%2^1 == 1 }); !reflect.DeepEqual(got, want) {
			t.Errorf("Filter() = %v, want %v", got, want)
		}
	}
}

func TestFilter(t *testing.T) {
	{
		in := []int{1, 2, 3, 4, 5, 6, 7, 8}
		want := []int{2, 4, 6, 8}
		if got := Filter(in, func(t int) bool { return t%2 == 0 }); !reflect.DeepEqual(got, want) {
			t.Errorf("Filter() = %v, want %v", got, want)
		}
	}
	{
		in := []int{1, 2, 3, 4, 5, 6, 7, 8}
		want := []int{3, 6}
		if got := Filter(in, func(t int) bool { return t%3 == 0 }); !reflect.DeepEqual(got, want) {
			t.Errorf("Filter() = %v, want %v", got, want)
		}
	}
}

func TestReduce(t *testing.T) {
	{
		initial := 9
		in := []int{1, 2, 3, 4, 5, 6, 7, 8}
		want := 45
		if got := Reduce(in, initial, func(cm, t int) int { return cm + t }); !reflect.DeepEqual(got, want) {
			t.Errorf("Filter() = %v, want %v", got, want)
		}
	}
	{
		type data struct {
			x int
			y int
		}
		initial := data{x: 4, y: 5}
		in := []int{1, 2, 3, 4, 5, 6, 7, 8}
		want := data{x: 40, y: 77}
		if got := Reduce(in, initial, func(cm data, t int) data { return data{x: cm.x + t, y: cm.y + t*2} }); !reflect.DeepEqual(got, want) {
			t.Errorf("Filter() = %v, want %v", got, want)
		}
	}
}
