package aviator_test

import (
	"os"

	. "github.com/JulzDiverse/aviator/aviator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aviator", func() {

	var file string

	Context("Test basic function functionality", func() {

		BeforeEach(func() {
			file = `spruce:
  - base: base.yml
    merge:
    - with:
        files:
        - another.yml
    to: some/destination/file.yml

fly:
  config: pipeline.yml
  vars:
  - personal.yml`
		})

		Context("ReadYaml", func() {
			It("should return an non empty aviator struct", func() {
				Expect(ReadYaml([]byte(file))).ShouldNot(BeNil())
				Expect(len(ReadYaml([]byte(file)).Spruce)).To(Equal(1))
			})

			It("should contain a fly and a spruce object with params", func() {
				Expect(ReadYaml([]byte(file)).Fly.Config).To(Equal("pipeline.yml"))
				Expect(ReadYaml([]byte(file)).Spruce[0].Base).To(Equal("base.yml"))
			})
		})

		Context("ConcatFileName", func() {
			It("should concat a file name concatinaed by the parant folder name and the file name", func() {
				filePath := "path/to/some/file.yml"
				fn, _ := ConcatFileName(filePath)
				Expect(fn).To(Equal("some_file.yml"))
			})
		})

	})

	Context("Integration Tests: Spruce Specific Files", func() {

		BeforeEach(func() {
			file = `spruce:
- base: ../integration/yamls/base.yml
  merge:
  - with:
      files:
      - another.yml
      in_dir: ../integration/yamls/
  - with:
      files:
      - ../integration/yamls/addons/sub1/file1.yml
  to: ../integration/tmp/tmp.yml
- base: ../integration/tmp/tmp.yml
  merge:
  - with:
      files:
      - ../integration/yamls/yet-another.yml
  to: ../integration/tmp/result.yml`
		})

		Context("ProcessSpruceChain", func() {
			It("Should generate a result.yml file", func() {
				ProcessSprucePlan(ReadYaml([]byte(file)).Spruce, false, false)
				Expect("../integration/tmp/result.yml").To(BeAnExistingFile())
			})
		})

		Context("Cleanup", func() {
			It("Should delete all files in tmp", func() {
				Cleanup("../integration/tmp/")
				Expect("../integration/tmp/result.yml").ShouldNot(BeAnExistingFile())
			})
		})
	})

	Context("Resolve env vars", func() {
		var file, resultFile string
		BeforeEach(func() {
			file = `spruce:
  - base: $BASE
    merge:
    - with:
        files:
        - ${INPUT_FILE}
    to: $TMP_DIR/destination/$FILE

fly:
  config: pipeline.yml
  vars:
  - personal.yml`

			resultFile = `spruce:
  - base: base.yml
    merge:
    - with:
        files:
        - another.yml
    to: some/destination/file.yml

fly:
  config: pipeline.yml
  vars:
  - personal.yml`

		})

		It("should resolve env variables", func() {
			os.Setenv("BASE", "base.yml")
			os.Setenv("INPUT_FILE", "another.yml")
			os.Setenv("TMP_DIR", "some")
			os.Setenv("FILE", "file.yml")

			resolved := ResolveEnvVars([]byte(file))
			Expect(string(resolved)).To(Equal(resultFile))
		})

	})

})
