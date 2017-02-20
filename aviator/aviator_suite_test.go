package aviator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAviator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aviator Suite")
}
