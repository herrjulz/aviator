package modifier_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestModifier(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Modifier Suite")
}
