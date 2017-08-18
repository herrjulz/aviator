package validator_test

import (
	"github.com/JulzDiverse/aviator/cockpit"
	. "github.com/JulzDiverse/aviator/validator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {

	var cfg cockpit.Spruce
	var validator *Validator

	BeforeEach(func() {
		cfg = cockpit.Spruce{
			Base: "base.yml",
			Merge: []cockpit.Merge{
				cockpit.Merge{
					With: cockpit.With{},
				},
			},
			To: "target.yml",
		}

		validator = New()
	})

	Context("Merge Validator", func() {
		Context("Merge Top Level Combinations", func() {
			Context("When 'with.files' is defined", func() {
				It("returns an error when with_in is also defined", func() {
					cfg.Merge[0].With.Files = []string{"fake"}
					cfg.Merge[0].WithIn = "path/"

					err := validator.ValidateSpruce([]cockpit.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'with', 'with_in', and 'with_all_in' are discrete parameters and cannot be defined together")))
				})

				It("returns an error when with_all_in is also defined", func() {
					cfg.Merge[0].With.Files = []string{"fake"}
					cfg.Merge[0].WithAllIn = "path/"

					err := validator.ValidateSpruce([]cockpit.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'with', 'with_in', and 'with_all_in' are discrete parameters and cannot be defined together")))
				})
			})

			Context("When 'with.WithIn' is defined", func() {
				It("returns an error when with_all_in is also defined", func() {
					cfg.Merge[0].WithIn = "path/one/"
					cfg.Merge[0].WithAllIn = "path/two/"

					err := validator.ValidateSpruce([]cockpit.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'with', 'with_in', and 'with_all_in' are discrete parameters and cannot be defined together")))
				})
			})
		})

		Context("With", func() {
			Context("When 'with' is not defined", func() {
				It("returns an error if 'in_dir' is defined", func() {
					cfg.Merge[0].With.InDir = "path/"

					err := validator.ValidateSpruce([]cockpit.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'with.in_dir' or 'with.skip_non_existing' can only be declared in combination with 'with.files'")))
				})

			})
		})

		Context("WithIn & WithAllIn", func() {
			Context("When 'with_in' nor 'with_all_in' is defined", func() {
				It("returns an error if 'except' is defined", func() {
					cfg.Merge[0].Except = []string{"fake", "fake2"}

					err := validator.ValidateSpruce([]cockpit.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'merge.except' is only allowed in combination with 'merge.with_in' or 'merge.with_all_in'")))
				})
			})

			Context("When 'with_in' is  defined", func() {
				It("returns NO error if 'except' is defined", func() {
					cfg.Merge[0].WithIn = "path/"
					cfg.Merge[0].Except = []string{"fake", "fake2"}

					err := validator.ValidateSpruce([]cockpit.Spruce{cfg})
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("When 'with_all_in' is  defined", func() {
				It("returns NO error if 'except' is defined", func() {
					cfg.Merge[0].WithAllIn = "path/"
					cfg.Merge[0].Except = []string{"fake", "fake2"}

					err := validator.ValidateSpruce([]cockpit.Spruce{cfg})
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

		Context("Regexp", func() {
			Context("when 'regexp' is defined", func() {
				It("returns an error if 'with', 'with_in', or 'with_all_in' is NOT defined", func() {
					cfg.Merge[0].Regexp = ".*.(yml)"

					err := validator.ValidateSpruce([]cockpit.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'merge.regexp' is only allowed in combination with 'merge.with', 'merge.with_in' or 'merge.with_all_in'")))
				})
			})
		})
	})
})
