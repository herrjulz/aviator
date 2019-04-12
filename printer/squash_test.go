package printer_test

import (
	"bytes"
	"fmt"
	"io"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/JulzDiverse/aviator/printer"
)

var _ = Describe("Squash", func() {
	var (
		files    []string
		expected string
		to       string
		output   string
	)

	BeforeEach(func() {
		files = []string{
			"file1",
			"file2",
			"file3",
		}

		expected = `@M{SQUASH FILES:}
	@w{file1}
	@w{file2}
	@w{file3}
	@M{to: dest}
`
		to = "dest"
	})

	JustBeforeEach(func() {
		output = captureOutputSquash(BeautyPrintSquash, files, to, fmt.Printf)
	})

	Context("BeautyPrintSquash", func() {
		It("prints the expected output", func() {
			Expect(output).To(Equal(expected))
		})
	})
})

func captureOutputSquash(f func([]string, string, Print), files []string, to string, printf Print) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f(files, to, printf)
	os.Stdout = old
	var buf bytes.Buffer
	w.Close()
	io.Copy(&buf, r)
	return buf.String()
}
