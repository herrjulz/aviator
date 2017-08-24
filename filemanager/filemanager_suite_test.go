package filemanager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFilemanager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filemanager Suite")
}
