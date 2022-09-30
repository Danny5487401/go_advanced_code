package main

import (
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"math/rand"
	"time"
)

func main() {
	producerFunc := func(ctx context.Context, next chan<- rxgo.Item) {
		for _, i := range []int{1, 2} {
			next <- rxgo.Item{
				V: i,
			}
		}
	}

	mapFunc := func(_ context.Context, i interface{}) (interface{}, error) {
		time.Sleep(time.Duration(rand.Int31()))
		return i.(int)*2 + 1, nil
	}

	reducerFunc := func(_ context.Context, acc interface{}, elem interface{}) (interface{}, error) {
		if acc == nil {
			return elem, nil
		}
		return acc.(int) + elem.(int), nil
	}

	observable := rxgo.Create([]rxgo.Producer{
		producerFunc,
	}).
		Map(mapFunc, rxgo.WithCPUPool()).
		Reduce(reducerFunc)

	for item := range observable.Observe() {
		fmt.Println(item.V)
	}
}
