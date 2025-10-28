package _5_ginkgo

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// 为Suite添加Specs
var _ = ginkgo.Describe("Book", func() {
	ginkgo.Context("myBook", func() {
		ginkgo.It("1", func() {
			err := Book(1)
			gomega.Expect(err).Should(gomega.BeNil())
		})

		ginkgo.It("2", func() {
			Book(2)
		})

		ginkgo.It("3", func() {
			Book(3)
		})
	})
})
