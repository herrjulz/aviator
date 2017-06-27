package spruce_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSpruce(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Spruce Suite")
}
