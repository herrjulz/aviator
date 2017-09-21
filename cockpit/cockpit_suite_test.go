package cockpit_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCockpit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cockpit Suite")
}
