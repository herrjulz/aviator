package processor_test

import (
	. "github.com/JulzDiverse/aviator/processor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processor", func() {
	var aviatorYaml string

	Context("New", func() {
		BeforeEach(func() {
			aviatorYaml = `spruce:
- base: base.yml
  merge:
  - with:
      files:
      - file1.yml
      - file2.yml
    regexp: ".*.(yml)"
    in_dir: /path/
    skip_non_existing: true
  - with_in: /path/to/dir/
    except: true
  to: result.yml`
		})

		Context("aviator.yml is valid", func() {
			It("returns a processor object with specified configuration", func() {
				processor, err := New([]byte(aviatorYaml))
				Expect(err).ToNot(HaveOccurred())

				Expect(len(processor.Aviator.Spruce[0].Merge[0].With.Files)).To(Equal(2))
			})
		})

	})

})
