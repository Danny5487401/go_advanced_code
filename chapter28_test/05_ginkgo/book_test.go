package _5_ginkgo

import (
	ginkgo "github.com/onsi/ginkgo/v2"
)

// 为Suite添加Specs
var _ = ginkgo.Describe("Book", func() {
	ginkgo.Context("myBook", func() {
		ginkgo.It("1", func() {
			Book(1)
		})

		ginkgo.It("2", func() {
			Book(2)
		})

		ginkgo.It("3", func() {
			Book(3)
		})
	})
})
