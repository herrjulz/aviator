package mingoak_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMingoak(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mingoak Suite")
}
