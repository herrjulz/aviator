package spruce_test

import (
	. "github.com/JulzDiverse/aviator/spruce"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spruce", func() {

	Context("CmdMergeEval", func() {
		It("simple merge two files", func() {
			opts := MergeOpts{
				Files: []string{
					"../integration/yamls/base.yml",
					"../integration/yamls/another.yml",
				},
			}

			result, err := CmdMergeEval(opts)

			Expect(err).To(BeNil())
			value, _ := result["word"]
			value2, _ := result["the"]
			Expect(value).To(Equal("yo!"))
			Expect(value2).To(Equal("base"))
		})

		It("should be able to prune", func() {
			opts := MergeOpts{
				Files: []string{
					"../integration/yamls/base.yml",
					"../integration/yamls/another.yml",
				},
				Prune: []string{
					"the",
				},
			}

			result, err := CmdMergeEval(opts)
			Expect(err).To(BeNil())
			value, _ := result["the"]
			Expect(value).To(BeNil())
		})
	})
})
