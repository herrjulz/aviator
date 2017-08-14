package processor_test

import (
	. "github.com/JulzDiverse/aviator/processor"
	"github.com/JulzDiverse/aviator/validator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processor", func() {

	//var spruceProcessor *SpruceProcessor
	var spruceConfig validator.Spruce

	Describe("Process() returns a byte and an error)", func() {
		BeforeEach(func() {
			spruceConfig = validator.Spruce{
				Base: "input.yml",
				Merge: []validator.Merge{
					validator.Merge{
						With: validator.With{
							Files: []string{"file.yml"},
						},
					},
				},
				To: "result.yml",
			}
		})

		It("", func() {
			_, err := Process(spruceConfig)
			Expect(err).ToNot(HaveOccurred())
		})

	})
})
