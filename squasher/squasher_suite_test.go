package squasher_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSquasher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Squasher Suite")
}
