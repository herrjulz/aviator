package processor_test

import (
	"github.com/JulzDiverse/aviator/cockpit"
	. "github.com/JulzDiverse/aviator/processor"
	fakes "github.com/JulzDiverse/aviator/processor/processorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processor", func() {

	var processor *Processor
	var spruceConfig []cockpit.Spruce
	var spruceClient *fakes.FakeSpruceClient

	Describe("", func() {
		BeforeEach(func() {
			cfg := cockpit.Spruce{
				Base: "input.yml",
				Merge: []cockpit.Merge{
					cockpit.Merge{
						With: cockpit.With{
							Files: []string{"file.yml"},
						},
					},
				},
				To: "result.yml",
			}
			spruceConfig = []cockpit.Spruce{cfg}
			spruceClient = new(fakes.FakeSpruceClient)
			processor = New(spruceClient)
		})

		It("", func() {
			_, err := processor.Process(spruceConfig)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
