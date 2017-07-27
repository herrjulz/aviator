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
  - with_in: path/to/dir/
    except:
    - file2.yml
  to: result.yml

- base: result.yml
  merge:
  - with_in: another/path/
  for_each_in: path/to/dir/
  except:
  - file2.yml
  regexp: ".*.(yml)"
  to_dir: some/tmp/dir/

- base: some/tmp/dir/file1.yml
  merge:
  - with_in: path/
  for_each:
  - foo.yml
  - bar.yml
  to_dir: foo/bar/

- base: base.yml
  prune:
  - some
  - properties
  merge:
  - with_in: foo/
  walk_through: foo/bar/
  to_dir: final/`
		})

		Context("aviator.yml is valid", func() {
			Context("spruce section", func() {
				It("returns a processor object with specified configuration", func() {
					processor, err := New([]byte(aviatorYaml))
					Expect(err).ToNot(HaveOccurred())

					Expect(len(processor.Aviator.Spruce[0].Merge[0].With.Files)).To(Equal(2))
					Expect(processor.Aviator.Spruce[0].Merge[1].WithIn).To(Equal("path/to/dir/"))
					Expect(len(processor.Aviator.Spruce[0].Merge[1].Except)).To(Equal(1))
					Expect(processor.Aviator.Spruce[0].Merge[0].Regexp).To(Equal(".*.(yml)"))
					Expect(processor.Aviator.Spruce[0].Merge[0].Skip).To(Equal(true))
					Expect(processor.Aviator.Spruce[0].To).To(Equal("result.yml"))
					Expect(len(processor.Aviator.Spruce[2].ForEach)).To(Equal(2))
					Expect(processor.Aviator.Spruce[1].ForEachIn).To(Equal("path/to/dir/"))
					Expect(len(processor.Aviator.Spruce[1].Except)).To(Equal(1))

				})

			})

		})

	})

})
