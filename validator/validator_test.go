package validator_test

import (
	"github.com/JulzDiverse/aviator"
	. "github.com/JulzDiverse/aviator/validator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {

	var cfg aviator.Spruce
	var validator *Validator

	BeforeEach(func() {
		cfg = aviator.Spruce{
			Base: "base.yml",
			Merge: []aviator.Merge{
				aviator.Merge{
					With: aviator.With{},
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

					err := validator.ValidateSpruce([]aviator.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'with', 'with_in', and 'with_all_in' are discrete parameters and cannot be defined together")))
				})

				It("returns an error when with_all_in is also defined", func() {
					cfg.Merge[0].With.Files = []string{"fake"}
					cfg.Merge[0].WithAllIn = "path/"

					err := validator.ValidateSpruce([]aviator.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'with', 'with_in', and 'with_all_in' are discrete parameters and cannot be defined together")))
				})
			})

			Context("When 'with.WithIn' is defined", func() {
				It("returns an error when with_all_in is also defined", func() {
					cfg.Merge[0].WithIn = "path/one/"
					cfg.Merge[0].WithAllIn = "path/two/"

					err := validator.ValidateSpruce([]aviator.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'with', 'with_in', and 'with_all_in' are discrete parameters and cannot be defined together")))
				})
			})
		})

		Context("With", func() {
			Context("When 'with' is not defined", func() {
				It("returns an error if 'in_dir' is defined", func() {
					cfg.Merge[0].With.InDir = "path/"

					err := validator.ValidateSpruce([]aviator.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'with.in_dir' or 'with.skip_non_existing' can only be declared in combination with 'with.files'")))
				})
			})
		})

		Context("WithIn & WithAllIn", func() {
			Context("When 'with_in' nor 'with_all_in' is defined", func() {
				It("returns an error if 'except' is defined", func() {
					cfg.Merge[0].Except = []string{"fake", "fake2"}

					err := validator.ValidateSpruce([]aviator.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'merge.except' is only allowed in combination with 'merge.with_in' or 'merge.with_all_in'")))
				})
			})

			Context("When 'with_in' is  defined", func() {
				It("returns NO error if 'except' is defined", func() {
					cfg.Merge[0].WithIn = "path/"
					cfg.Merge[0].Except = []string{"fake", "fake2"}

					err := validator.ValidateSpruce([]aviator.Spruce{cfg})
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("When 'with_all_in' is  defined", func() {
				It("returns NO error if 'except' is defined", func() {
					cfg.Merge[0].WithAllIn = "path/"
					cfg.Merge[0].Except = []string{"fake", "fake2"}

					err := validator.ValidateSpruce([]aviator.Spruce{cfg})
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

		Context("Regexp", func() {
			Context("when 'regexp' is defined", func() {
				It("returns an error if 'with', 'with_in', or 'with_all_in' is NOT defined", func() {
					cfg.Merge[0].Regexp = ".*.(yml)"

					err := validator.ValidateSpruce([]aviator.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'merge.regexp' is only allowed in combination with 'merge.with', 'merge.with_in' or 'merge.with_all_in'")))
				})
			})
		})
	})

	Context("ForEach Validator", func() {
		Context("Top level combination", func() {
			Context("'files' and 'in' are mutually exclusive", func() {
				It("returns an error if both parameters are declared", func() {
					cfg.ForEach.Files = []string{"file", "file2"}
					cfg.ForEach.In = "path/"

					err := validator.ValidateSpruce([]aviator.Spruce{cfg})
					Expect(err).To(HaveOccurred())

					Expect(err).To(MatchError(ContainSubstring(
						"INVALID SYNTAX: Mutually exclusive parameters declared 'for_each.in' and 'for_each.files'",
					)))
				})
			})

			Context("'files' combinations", func() {
				Context("When 'files' is not declared", func() {
					It("returns an error if 'in_dir' is declared", func() {
						cfg.ForEach.InDir = "path/"

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).To(HaveOccurred())

						Expect(err).To(MatchError(ContainSubstring(
							"INVALID SYNTAX: 'for_each.in_dir' and 'for_each.skip_non_existing' can only be declared in combination with 'for_each.files'",
						)))
					})

					It("returns an error if 'skip_non_existing' is enabled", func() {
						cfg.ForEach.Skip = true

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).To(HaveOccurred())

						Expect(err).To(MatchError(ContainSubstring(
							"INVALID SYNTAX: 'for_each.in_dir' and 'for_each.skip_non_existing' can only be declared in combination with 'for_each.files'",
						)))
					})

					It("returns an error if 'skip_non_existing' and 'in_dir' are declared", func() {
						cfg.ForEach.Skip = true
						cfg.ForEach.InDir = "path/"

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).To(HaveOccurred())

						Expect(err).To(MatchError(ContainSubstring(
							"INVALID SYNTAX: 'for_each.in_dir' and 'for_each.skip_non_existing' can only be declared in combination with 'for_each.files'",
						)))
					})
				})

				Context("When 'files' is declared", func() {
					It("can be combined with 'in_dir'", func() {
						cfg.ForEach.InDir = "path/"
						cfg.ForEach.Files = []string{"file", "file2"}

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).ToNot(HaveOccurred())
					})

					It("can be combined with 'skip_non_existing'", func() {
						cfg.ForEach.Skip = true
						cfg.ForEach.Files = []string{"file", "file2"}

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).ToNot(HaveOccurred())
					})

					It("can be combined with 'skip_non_existing' and with 'in_dir'", func() {
						cfg.ForEach.Skip = true
						cfg.ForEach.InDir = "path/"
						cfg.ForEach.Files = []string{"file", "file2"}

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).ToNot(HaveOccurred())
					})
				})
			})

			Context("'in' combinations", func() {
				Context("When 'in' is not declared", func() {
					It("returns an error if 'except' is declared", func() {
						cfg.ForEach.Except = []string{"fake"}

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).To(HaveOccurred())

						Expect(err).To(MatchError(ContainSubstring(
							"INVALID SYNTAX: 'for_each.except' and 'for_each.include_sub_dirs' can only be declared in combination with 'for_each.in'",
						)))
					})

					It("returns an error if 'include_sub_dirs' is declared", func() {
						cfg.ForEach.SubDirs = true

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).To(HaveOccurred())

						Expect(err).To(MatchError(ContainSubstring(
							"INVALID SYNTAX: 'for_each.except' and 'for_each.include_sub_dirs' can only be declared in combination with 'for_each.in'",
						)))
					})
				})

				Context("When 'include_sub_dirs' is not declared", func() {
					It("returns an error if 'copy_parents' is enabled", func() {
						cfg.ForEach.CopyParents = true

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).To(HaveOccurred())

						Expect(err).To(MatchError(ContainSubstring(
							"INVALID SYNTAX: 'for_each.copy_parents', 'for_each.enable_matching', 'for_each.for_all' can only be declared in combination with 'for_each.inlcude_sub_dirs'",
						)))
					})
				})
			})
			Context("regexp combinations", func() {
				Context("When 'for_each.in' or 'for_each.files' not declared", func() {
					It("It returns an error if declared 'for_ach.regexp'", func() {
						cfg.ForEach.Regexp = ".*.(yml)"

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).To(HaveOccurred())

						Expect(err).To(MatchError(ContainSubstring("INVALID SYNTAX: 'for_each.regexp' is only allowed in combination with 'for_each.in', 'for_each.files'")))
					})
				})

				Context("When 'for_each.in' is declared", func() {
					It("can be combined with 'for_each.regexp'", func() {
						cfg.ForEach.Regexp = ".*.(yml)"
						cfg.ForEach.In = "path/"

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).ToNot(HaveOccurred())
					})
				})
				Context("When 'for_each.files' is declared", func() {
					It("can be combined with 'for_each.regexp'", func() {
						cfg.ForEach.Regexp = ".*.(yml)"
						cfg.ForEach.Files = []string{"fake"}

						err := validator.ValidateSpruce([]aviator.Spruce{cfg})
						Expect(err).ToNot(HaveOccurred())
					})
				})
			})
		})
	})
})
