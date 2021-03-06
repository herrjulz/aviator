package printer_test

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/JulzDiverse/aviator"
	. "github.com/JulzDiverse/aviator/printer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Printer", func() {
	var (
		opts     aviator.MergeConf
		expected string
		warnings []string
		to       string
	)

	BeforeEach(func() {
		opts = aviator.MergeConf{
			Files: []string{"file", "file2"},
			Prune: []string{"props", "meta"},
		}
		expected = `@G{SPRUCE MERGE:}
	@C{--prune} props
	@C{--prune} meta
	file
	file2
	@G{to: dest}

	@Y{WARNINGS:}
	@y{skipped}:@Y{x}
	@y{skipped}:@Y{y}


`

		warnings = []string{"skipped:x", "skipped:y"}
		to = "dest"
	})

	Context("BeautifulPrint", func() {
		It("prints the expected output", func() {
			output := captureOutput(BeautyfulPrint, opts, to, warnings, true, fmt.Printf)
			Expect(output).To(Equal(expected))
		})
	})
})

func captureOutput(f func(aviator.MergeConf, string, []string, bool, Print), opts aviator.MergeConf, to string, warnings []string, verbose bool, printf Print) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f(opts, to, warnings, verbose, printf)
	os.Stdout = old
	var buf bytes.Buffer
	w.Close()
	io.Copy(&buf, r)
	return buf.String()
}
