package aviator_test

import (
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
    chain:
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

		Context("CreateSpruceCommand", func() {
			It("Should create an expected spruce command", func() {
				command := ProcessChain(ReadYaml([]byte(file)).Spruce[0])
				Expect(command).NotTo(BeEmpty())
				Expect(command).Should(HaveLen(4))
				Expect(command[2]).To(Equal("base.yml"))
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
  chain:
  - with:
      files:
      - another.yml
      in_dir: ../integration/yamls/
  - with:
      files:
      - ../integration/yamls/addons/sub1/file1.yml
  to: ../integration/tmp/tmp.yml
- base: ../integration/tmp/tmp.yml
  chain:
  - with:
      files:
      - ../integration/yamls/yet-another.yml
  to: ../integration/tmp/result.yml`
		})

		Context("SpruceToFile", func() {
			It("Should generate a file", func() {
				avi := ReadYaml([]byte(file))
				SpruceToFile(ProcessChain(avi.Spruce[0]), avi.Spruce[0].DestFile)
				Expect("../integration/tmp/tmp.yml").To(BeAnExistingFile())
			})
		})

		Context("ProcessSpruceChain", func() {
			It("Should generate a result.yml file", func() {
				ProcessSprucePlan(ReadYaml([]byte(file)).Spruce)
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
})
