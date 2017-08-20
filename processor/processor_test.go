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

	Describe("Process", func() {
		var cfg cockpit.Spruce
		Context("Default Merge", func() {
			BeforeEach(func() {
				cfg = cockpit.Spruce{
					Base: "input.yml",
					Merge: []cockpit.Merge{
						cockpit.Merge{
							With: cockpit.With{},
						},
					},
					To: "result.yml",
				}

			})
			Context("Merge Section", func() {
				Context("Using Merge.With.Files", func() {
					It("includes the right files with the right amount in the merge ", func() {
						cfg.Merge[0].With.Files = []string{"file.yml"}
						spruceConfig = []cockpit.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = New(spruceClient)

						_, err := processor.Process(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(2))
						Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
						Expect(mergeOpts.Files[1]).To(Equal("file.yml"))
					})
				})

				Context("Using Merge.With.Files in combination with InDir", func() {
					It("includes the right files with the right amount in the merge ", func() {
						cfg.Merge[0].With.Files = []string{"fake.yml", "fake2.yml"}
						cfg.Merge[0].With.InDir = "integration/yamls/"

						spruceConfig = []cockpit.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = New(spruceClient)

						_, err := processor.Process(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(3))
						Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
						Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/fake.yml"))
						Expect(mergeOpts.Files[2]).To(Equal("integration/yamls/fake2.yml"))
					})
				})

				Context("Using Merge.With.Files in combination with SkipNonExisting", func() {
					It("excludes non existing files from the merge", func() {
						cfg.Merge[0].With.Files = []string{"nonExisting.yml", "fake.yml", "fake2.yml"}
						cfg.Merge[0].With.InDir = "integration/yamls/"
						cfg.Merge[0].With.Skip = true

						spruceConfig = []cockpit.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = New(spruceClient)

						_, err := processor.Process(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(3))
						Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
						Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/fake.yml"))
						Expect(mergeOpts.Files[2]).To(Equal("integration/yamls/fake2.yml"))
					})
				})

				Context("Using Merge.With.Files including an nonexisting file", func() {
					It("includes the right files with the right amount in the merge ", func() {
						cfg.Merge[0].With.Files = []string{"nonExisting.yml", "fake.yml", "fake2.yml"}
						cfg.Merge[0].With.InDir = "integration/yamls/"

						spruceConfig = []cockpit.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = New(spruceClient)

						_, err := processor.Process(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(4))
						Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
						Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/nonExisting.yml"))
						Expect(mergeOpts.Files[2]).To(Equal("integration/yamls/fake.yml"))
					})
				})

				Context("Using Merge.WithIn", func() {
					It("includes all files within a directory, but not subdirectories ", func() {
						cfg.Merge[0].WithIn = "integration/yamls/"

						spruceConfig = []cockpit.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = New(spruceClient)

						_, err := processor.Process(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(4))
						Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
						Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/base.yml"))
						Expect(mergeOpts.Files[2]).To(Equal("integration/yamls/fake.yml"))
						Expect(mergeOpts.Files[3]).To(Equal("integration/yamls/fake2.yml"))
					})
				})

				Context("Using Merge.WithIn in combination with Except", func() {
					It("includes all files within a directory, except files listed in Except ", func() {
						cfg.Merge[0].WithIn = "integration/yamls/"
						cfg.Merge[0].Except = []string{"base.yml", "fake.yml"}

						spruceConfig = []cockpit.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = New(spruceClient)

						_, err := processor.Process(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(2))
						Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
						Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/fake2.yml"))
					})
				})

				Context("Using Merge.WithIn in combination with Regexp", func() {
					It("includes only files within a directory matching the regexp", func() {
						cfg.Merge[0].WithIn = "integration/yamls/"
						cfg.Merge[0].Regexp = "base.yml"

						spruceConfig = []cockpit.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = New(spruceClient)

						_, err := processor.Process(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(2))
						Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
						Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/base.yml"))
					})
				})

				Context("Using Merge.WithIn in combination with Regexp and Except", func() {
					It("includes only files within a directory matching the regexp and not part of Except array", func() {
						cfg.Merge[0].WithIn = "integration/yamls/"
						cfg.Merge[0].Regexp = "fake.*.yml"
						cfg.Merge[0].Except = []string{"fake.yml"}

						spruceConfig = []cockpit.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = New(spruceClient)

						_, err := processor.Process(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(2))
						Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
						Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/fake2.yml"))
					})
				})

				Context("Using Merge.WithAllIn", func() {
					It("includes all files within a directory and all subdirectories", func() {
						cfg.Merge[0].WithAllIn = "integration/yamls/"

						spruceConfig = []cockpit.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = New(spruceClient)

						_, err := processor.Process(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(7))
						Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/addons/sub1/file1.yml"))
					})
				})

				Context("Using Merge.WithAllIn in combination with Regexp", func() {
					It("includes all files within a directory and all subdirectories matching the regexp", func() {
						cfg.Merge[0].WithAllIn = "integration/yamls/"
						cfg.Merge[0].Regexp = "file.*.yml"

						spruceConfig = []cockpit.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = New(spruceClient)

						_, err := processor.Process(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(4))
						Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
						Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/addons/sub1/file1.yml"))
						Expect(mergeOpts.Files[2]).To(Equal("integration/yamls/addons/sub1/file2.yml"))
						Expect(mergeOpts.Files[3]).To(Equal("integration/yamls/addons/sub2/file1.yml"))
					})
				})
			})
		})

		Context("ForEach", func() {
			Context("Files", func() {

			})
		})
	})
})
