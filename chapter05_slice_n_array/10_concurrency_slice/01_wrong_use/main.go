package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		slc = []int{}
		n   = 10000
		wg  sync.WaitGroup
	)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			slc = append(slc, i)
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("len:", len(slc))
	fmt.Println("done")
}

// Output:
/*
len: 9050
done

*/
