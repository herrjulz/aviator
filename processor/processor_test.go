package processor_test

import (
	. "github.com/JulzDiverse/aviator/processor"
	"github.com/JulzDiverse/aviator/validator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processor", func() {

	//var spruceProcessor *SpruceProcessor
	var spruceConfig []validator.Spruce

	BeforeEach(func() {
		cfg := validator.Spruce{
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
		spruceConfig = []validator.Spruce{cfg}
	})

	Describe("Process() returns a byte and an error", func() {
		It("", func() {
			_, err := Process(spruceConfig)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
