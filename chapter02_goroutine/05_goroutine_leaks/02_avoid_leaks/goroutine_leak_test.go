package main

import (
	"testing"

	"go.uber.org/goleak"
)

func TestGetDataWithGoleak(t *testing.T) {
	defer goleak.VerifyNone(t)
	main()
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
