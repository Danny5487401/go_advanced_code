package _5_ginkgo_test

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// 初始化 ginkgo bootstrap --nodot

func Test05Ginkgo(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "05Ginkgo Suite")
}
