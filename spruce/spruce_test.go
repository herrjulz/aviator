package spruce_test

import (
	"github.com/JulzDiverse/aviator"
	"github.com/JulzDiverse/aviator/filemanager"
	. "github.com/JulzDiverse/aviator/spruce"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spruce", func() {

	var spruce *SpruceClient

	BeforeEach(func() {
		spruce = NewWithFileFilemanager(
			filemanager.Store(true, false), true,
		)
	})

	Context("CmdMergeEval", func() {
		It("simple merge two files", func() {
			opts := aviator.MergeConf{
				Files: []string{
					"../processor/integration/yamls/base.yml",
					"../processor/integration/yamls/fake.yml",
				},
			}

			result, err := spruce.MergeWithOptsRaw(opts)

			Expect(err).To(BeNil())
			value, _ := result["word"]
			value2, _ := result["the"]
			Expect(value).To(Equal("yo!"))
			Expect(value2).To(Equal("base"))
		})

		It("should be able to prune", func() {
			opts := aviator.MergeConf{
				Files: []string{
					"../processor/integration/yamls/base.yml",
					"../processor/integration/yamls/fake.yml",
				},
				Prune: []string{
					"the",
				},
			}

			result, err := spruce.MergeWithOptsRaw(opts)
			Expect(err).To(BeNil())
			value, _ := result["the"]
			Expect(value).To(BeNil())
		})
	})
})
