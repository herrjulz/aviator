package evaluator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEvaluator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Evaluator Suite")
}
