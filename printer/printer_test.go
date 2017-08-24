package printer_test

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/JulzDiverse/aviator/cockpit"
	. "github.com/JulzDiverse/aviator/printer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Printer", func() {
	var opts cockpit.MergeConf
	var expected string

	BeforeEach(func() {
		opts = cockpit.MergeConf{
			Files:    []string{"file", "file2"},
			Prune:    []string{"props", "meta"},
			Warnings: []string{"skipped:x", "skipped:y"},
			To:       "dest",
		}
		expected = `SPRUCE MERGE:
	@C{--prune} props
	@C{--prune} meta
	file
	file2
	@G{to: dest}

	@Y{WARNINGS:}
	@y{skipped}:@Y{x}
	@y{skipped}:@Y{y}


`
	})

	Context("BeautifulPrint", func() {
		It("prints the expected output", func() {
			output := captureOutput(BeautyfulPrint, opts, fmt.Printf, true)
			Expect(output).To(Equal(expected))
		})
	})
})

func captureOutput(f func(cockpit.MergeConf, Print, bool), opts cockpit.MergeConf, printf Print, verbose bool) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f(opts, printf, verbose)
	os.Stdout = old
	var buf bytes.Buffer
	w.Close()
	io.Copy(&buf, r)
	return buf.String()
}
