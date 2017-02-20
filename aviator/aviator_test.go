package aviator_test

import (
	"fmt"
	. "masterjulz/aviator/aviator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aviator", func() {

	var file string

	Context("Test basic function functionality", func() {

		BeforeEach(func() {
			file = `spruce:
- base: base.yml
  with:
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
				command := CreateSpruceCommand(ReadYaml([]byte(file)).Spruce[0])
				Expect(command).NotTo(BeEmpty())
				Expect(command).Should(HaveLen(4))
				Expect(command[2]).To(Equal("base.yml"))
			})
		})

		Context("ConcatFileName", func() {
			It("should concat a file name concatinaed by the parant folder name and the file name", func() {
				filePath := "path/to/some/file.yml"
				Expect(ConcatFileName(filePath)).To(Equal("some_file.yml"))
			})
		})
	})

	Context("Integration Tests: Spruce Specific Files", func() {

		BeforeEach(func() {
			file = `spruce:
- base: ../integration/yamls/base.yml
  prune:
  - meta
  with:
  - ../integration/yamls/another.yml
  to: ../integration/tmp/tmp.yml
- base: ../integration/tmp/tmp.yml
  with:
  - ../integration/yamls/yet-another.yml
  to: ../integration/tmp/result.yml

fly:
  config: pipeline.yml
  vars:
  - personal.yml`
		})

		Context("SpruceToFile", func() {
			It("Should generate a file", func() {
				avi := ReadYaml([]byte(file))
				SpruceToFile(CreateSpruceCommand(avi.Spruce[0]), avi.Spruce[0].DestFile)
				Expect("../integration/tmp/tmp.yml").To(BeAnExistingFile())
			})
		})

		Context("ProcessSpruceChain", func() {
			It("Should generate a result.yml file", func() {
				ProcessSpruceChain(ReadYaml([]byte(file)).Spruce)
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

	Context("Integration Tests: Spruce base.yml with files in a direcotry", func() {

		BeforeEach(func() {
			file = `spruce:
- base: ../integration/yamls/base.yml
  prune:
  - meta
  with_in: ../integration/yamls/addons/sub1/
  to: ../integration/tmp/result.yml

fly:
  config: pipeline.yml
  vars:
  - personal.yml`
		})

		Context("SpruceToFile", func() {
			It("Should spruce base.yml with all files in a specified directory", func() {
				avi := ReadYaml([]byte(file))
				cmd := CreateSpruceCommand(avi.Spruce[0])
				SpruceToFile(cmd, avi.Spruce[0].DestFile)

				Expect(cmd[5]).To(Equal("../integration/yamls/addons/sub1/file1.yml"))
				Expect(cmd[6]).To(Equal("../integration/yamls/addons/sub1/file2.yml"))
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

	Context("Integration Tests: Spruce base.yml for each file in a direcotry", func() {

		BeforeEach(func() {
			file = `spruce:
- base: ../integration/yamls/base.yml
  prune:
  - meta
  for_each_in: ../integration/yamls/addons/sub1/
  to_dir: ../integration/tmp/

fly:
  config: pipeline.yml
  vars:
  - personal.yml`
		})

		Context("forEachIn", func() {
			It("should spruce base.yml with each file in a specified directory seperately", func() {
				avi := ReadYaml([]byte(file))
				ForEachIn(avi.Spruce[0])
				fmt.Println("WHHHHHHHHAAATTTTTT")
				Expect("../integration/tmp/sub1_file1.yml").To(BeAnExistingFile())
				Expect("../integration/tmp/sub1_file2.yml").To(BeAnExistingFile())
			})
		})

		Context("Cleanup", func() {
			It("Should delete all files in tmp", func() {
				Cleanup("../integration/tmp/")
				Expect("../integration/tmp/sub1_file1.yml").ShouldNot(BeAnExistingFile())
				Expect("../integration/tmp/sub1_file2.yml").ShouldNot(BeAnExistingFile())
			})
		})
	})

	Context("Integration Tests: Spruce base.yml for each file in all subdirecotries", func() {

		BeforeEach(func() {
			file = `spruce:
- base: ../integration/yamls/base.yml
  prune:
  - meta
  walk_through: ../integration/yamls/addons/
  to_dir: ../integration/tmp/

fly:
  config: pipeline.yml
  vars:
  - personal.yml`
		})

		Context("Walk", func() {
			It("should spruce base.yml with each file in all subdirectories seperately", func() {
				avi := ReadYaml([]byte(file))
				Walk(avi.Spruce[0])

				Expect("../integration/tmp/sub1_file1.yml").To(BeAnExistingFile())
				Expect("../integration/tmp/sub1_file2.yml").To(BeAnExistingFile())
				Expect("../integration/tmp/sub2_file1.yml").To(BeAnExistingFile())
			})
		})

		Context("Cleanup", func() {
			It("Should delete all files in tmp", func() {
				Cleanup("../integration/tmp/")
				Expect("../integration/tmp/sub1_file1.yml").ShouldNot(BeAnExistingFile())
				Expect("../integration/tmp/sub1_file2.yml").ShouldNot(BeAnExistingFile())
				Expect("../integration/tmp/sub2_file1.yml").ShouldNot(BeAnExistingFile())
			})
		})
	})

	Context("Integration Tests: Spruce base.yml for each file in all subdirecotries", func() {

		BeforeEach(func() {
			file = `spruce:
- base: ../integration/yamls/base.yml
  prune:
  - meta
  walk_through: ../integration/yamls/addons/
  to_dir: ../integration/tmp/
  regexp: file1

fly:
  config: pipeline.yml
  vars:
  - personal.yml`
		})

		Context("Walk", func() {
			It("should only create files which matches the regexp", func() {
				avi := ReadYaml([]byte(file))
				Walk(avi.Spruce[0])

				Expect("../integration/tmp/sub1_file1.yml").To(BeAnExistingFile())
				Expect("../integration/tmp/sub1_file2.yml").NotTo(BeAnExistingFile())
				Expect("../integration/tmp/sub2_file1.yml").To(BeAnExistingFile())
			})
		})

		Context("Cleanup", func() {
			It("Should delete all files in tmp", func() {
				Cleanup("../integration/tmp/")
				Expect("../integration/tmp/sub1_file1.yml").ShouldNot(BeAnExistingFile())
				Expect("../integration/tmp/sub1_file2.yml").ShouldNot(BeAnExistingFile())
				Expect("../integration/tmp/sub2_file1.yml").ShouldNot(BeAnExistingFile())
			})
		})
	})

})
