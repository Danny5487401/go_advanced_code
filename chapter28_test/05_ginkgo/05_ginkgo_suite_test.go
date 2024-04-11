package zero5_ginkgo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Test05Ginkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "05Ginkgo Suite")
}
