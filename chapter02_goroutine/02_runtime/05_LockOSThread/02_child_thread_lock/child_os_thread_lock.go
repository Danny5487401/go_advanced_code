// This Go program demonstrates the behavior of runtime.LockOSThread() and runtime.UnlockOSThread().
// runtime.LockOSThread() forces the wiring of a Goroutine to a system thread.
// No other goroutine can run in the same thread, unless runtime.UnlockOSThread() is called at some point.
//
// According to the manual of phthread_self, when many threads are created, the system may reassign an ID that was used by a terminated thread to a new thread.
//
// This programs shows that thread IDs are indeed reused (tested on Linux and Mac), but that the thread itself is actually destroyed.
// To prove this, we store used thread IDs in a global variable, and some data in each local thread using pthread_setspecific().
//
// When we hint the Go runtime to reuse the threads by calling runtime.UnlockOSThread(), we can see that the local data is still available when a thread is reused.

package main

/*
#include <stdlib.h>
#include <pthread.h>

int setspecific(pthread_key_t key, int i) {
	int *ptr = calloc(1, sizeof(int)); // memory leak on purpose
	*ptr = i;
	return pthread_setspecific(key, ptr);
}
*/
import "C"
import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

const nbGoroutines = 1000

type goroutine struct {
	num int    // App specific goroutine ID
	id  uint64 // Internal goroutine ID (debug only, do not rely on this in real programs)
}

var seenThreadIDs map[C.pthread_t]goroutine = make(map[C.pthread_t]goroutine, nbGoroutines+1)
var seenThreadIDsMutex sync.RWMutex

// getGID gets the current goroutine ID (copied from https://blog.sgmansfield.com/2015/12/goroutine-ids/)
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// isThreadIDReused checks if the passed thread ID has already be used before
func isThreadIDReused(t1 C.pthread_t, currentGoroutine goroutine) bool {
	seenThreadIDsMutex.RLock()
	defer seenThreadIDsMutex.RUnlock()
	for t2, previousGoroutine := range seenThreadIDs {
		if C.pthread_equal(t1, t2) != 0 {
			fmt.Printf("Thread ID reused (previous goroutine: %v, current goroutine: %v)\n", previousGoroutine, currentGoroutine)

			return true
		}
	}

	return false
}

func main() {
	runtime.LockOSThread()
	seenThreadIDsMutex.Lock()
	seenThreadIDs[C.pthread_self()] = goroutine{0, getGID()}
	seenThreadIDsMutex.Unlock()

	// It could be better to use C.calloc() to prevent the GC to destroy the key
	var tlsKey C.pthread_key_t
	if C.pthread_key_create(&tlsKey, nil) != 0 {
		panic("problem creating pthread key")
	}

	for i := 1; i <= nbGoroutines; i++ {
		go func(i int) {
			runtime.LockOSThread()
			// Uncomment the following line to see how the runtime behaves when threads can be reused
			//defer runtime.UnlockOSThread()

			// Check if data has already been associated with this thread
			oldI := C.pthread_getspecific(tlsKey)
			if oldI != nil {
				fmt.Printf("Thread reused, getspecific not empty (%d)\n", *(*C.int)(oldI))
			}

			g := goroutine{i, getGID()}

			// Get the current thread ID
			t := C.pthread_self()
			isThreadIDReused(t, g)

			// Associate some data to the local thread
			if C.setspecific(tlsKey, C.int(i)) != 0 {
				panic("problem setting specific")
			}

			// Add the current thread ID in the list of already used IDs
			seenThreadIDsMutex.Lock()
			defer seenThreadIDsMutex.Unlock()
			seenThreadIDs[C.pthread_self()] = g
		}(i)
	}
}
