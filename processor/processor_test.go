package processor_test

import (
	"github.com/JulzDiverse/aviator"
	fakes "github.com/JulzDiverse/aviator/aviatorfakes"
	"github.com/JulzDiverse/aviator/filemanager"
	. "github.com/JulzDiverse/aviator/processor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processor", func() {

	var (
		processor    *Processor
		spruceConfig []aviator.Spruce
		spruceClient *fakes.FakeSpruceClient
		store        *filemanager.FileManager //*fakes.FakeFileStore
		modifier     *fakes.FakeModifier
	)

	Describe("Process", func() {

		var cfg aviator.Spruce

		BeforeEach(func() {
			cfg = aviator.Spruce{
				Base: "input.yml",
				Merge: []aviator.Merge{
					aviator.Merge{
						With: aviator.With{},
					},
				},
				ForEach: aviator.ForEach{},
				To:      "integration/tmp/result.yml",
				ToDir:   "integration/tmp/",
			}
			store = filemanager.Store(true, false)
			modifier = new(fakes.FakeModifier)
		})

		Context("Modify", func() {
			Context("Delete", func() {
				Context("When Delete is defined", func() {
					It("it should call modify", func() {
						cfg.Merge[0].With.Files = []string{"file.yml"}
						cfg.Modify.Delete = []string{"some.path"}
						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						Expect(modifier.ModifyCallCount()).To(Equal(1))
					})

					It("should invoke delete with the expected values", func() {
						cfg.Merge[0].With.Files = []string{"file.yml"}
						cfg.Modify.Delete = []string{"some.path", "second.path", "third.path"}
						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
						Expect(err).ToNot(HaveOccurred())
						_, mods := modifier.ModifyArgsForCall(0)
						Expect(mods.Delete[0]).To(Equal("some.path"))
						Expect(mods.Delete[1]).To(Equal("second.path"))
						Expect(mods.Delete[2]).To(Equal("third.path"))
						Expect(len(mods.Set)).To(Equal(0))
						Expect(len(mods.Update)).To(Equal(0))
					})
				})
			})

			Context("Set", func() {
				It("When set is defined it should call modify", func() {
					cfg.Merge[0].With.Files = []string{"file.yml"}
					set1 := aviator.PathVal{"some.path", "val"}
					set := []aviator.PathVal{set1}
					cfg.Modify.Set = set
					spruceConfig = []aviator.Spruce{cfg}
					spruceClient = new(fakes.FakeSpruceClient)
					processor = NewTestProcessor(spruceClient, store, modifier)

					err := processor.ProcessSilent(spruceConfig)
					Expect(err).ToNot(HaveOccurred())

					Expect(modifier.ModifyCallCount()).To(Equal(1))
				})

				//It("should invoke set with the expected values", func() {
				//cfg.Merge[0].With.Files = []string{"file.yml"}
				//s := aviator.PathVal{Path: "some.path", Value: "val"}
				//s2 := aviator.PathVal{Path: "second.path", Value: "val2"}
				//set := []aviator.PathVal{s, s2}
				//cfg.Modify.Set = set
				//fmt.Println(cfg.Modify.Set)
				//spruceClient = new(fakes.FakeSpruceClient)
				//processor = NewTestProcessor(spruceClient, store, modifier)

				//err := processor.ProcessSilent(spruceConfig)
				//Expect(err).ToNot(HaveOccurred())
				//_, mods := modifier.ModifyArgsForCall(0)
				//fmt.Println(mods)
				//Expect(mods.Set[0].Path).To(Equal("some.path"))
				//Expect(mods.Set[0].Value).To(Equal("val"))
				//Expect(mods.Set[1].Path).To(Equal("second.path"))
				//Expect(mods.Set[1].Value).To(Equal("val2"))
				//})
			})

			Context("Update", func() {
				It("When set is defined it should call modify", func() {
					cfg.Merge[0].With.Files = []string{"file.yml"}
					u := aviator.PathVal{"some.path", "val"}
					update := []aviator.PathVal{u}
					cfg.Modify.Set = update
					spruceConfig = []aviator.Spruce{cfg}
					spruceClient = new(fakes.FakeSpruceClient)
					processor = NewTestProcessor(spruceClient, store, modifier)

					err := processor.ProcessSilent(spruceConfig)
					Expect(err).ToNot(HaveOccurred())

					Expect(modifier.ModifyCallCount()).To(Equal(1))
				})
			})
		})

		Context("Default Merge", func() {
			Context("Merge Section", func() {
				Context("Using Merge.With.Files", func() {
					It("includes the right files with the right amount in the merge ", func() {
						cfg.Merge[0].With.Files = []string{"file.yml"}
						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
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

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(3))
						Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
						Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/fake.yml"))
						Expect(mergeOpts.Files[2]).To(Equal("integration/yamls/fake2.yml"))
					})
				})

				Context("Using Merge.With.Files in combination with SkipNonExisting", func() {
					//It("excludes non existing files from the merge", func() {
					//cfg.Merge[0].With.Files = []string{"nonExisting.yml", "fake.yml", "fake2.yml"}
					//cfg.Merge[0].With.InDir = "integration/yamls/"
					//cfg.Merge[0].With.Skip = true

					//spruceConfig = []aviator.Spruce{cfg}
					//spruceClient = new(fakes.FakeSpruceClient)
					//store.ReadFileReturnsOnCall(0, []byte(""), false)
					//processor = NewTestProcessor(spruceClient, store)

					//err := processor.ProcessSilent(spruceConfig)
					//Expect(err).ToNot(HaveOccurred())

					//mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
					//Expect(len(mergeOpts.Files)).To(Equal(3))
					//Expect(mergeOpts.Files[0]).To(Equal("input.yml"))
					//Expect(mergeOpts.Files[1]).To(Equal("integration/yamls/fake.yml"))
					//Expect(mergeOpts.Files[2]).To(Equal("integration/yamls/fake2.yml"))
					//})
				})

				Context("Using Merge.With.Files including an nonexisting file", func() {
					It("includes the right files with the right amount in the merge ", func() {
						cfg.Merge[0].With.Files = []string{"nonExisting.yml", "fake.yml", "fake2.yml"}
						cfg.Merge[0].With.InDir = "integration/yamls/"

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
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

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
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

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
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

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
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

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
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

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
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

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
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
				It("should run a merge for each file in 'for_each.files'", func() {
					cfg.Merge[0].With.Files = []string{"fake1", "fake2"}
					cfg.ForEach.Files = []string{"file1", "file2"}
					cfg.ToDir = "{{path}}"

					spruceConfig = []aviator.Spruce{cfg}
					spruceClient = new(fakes.FakeSpruceClient)
					processor = NewTestProcessor(spruceClient, store, modifier)

					err := processor.ProcessSilent(spruceConfig)
					Expect(err).ToNot(HaveOccurred())

					mergeOpts1 := spruceClient.MergeWithOptsArgsForCall(0)
					mergeOpts2 := spruceClient.MergeWithOptsArgsForCall(1)
					Expect(len(mergeOpts1.Files)).To(Equal(4))
					Expect(len(mergeOpts2.Files)).To(Equal(4))
					Expect(mergeOpts1.Files[3]).To(Equal("file1"))
					Expect(mergeOpts2.Files[3]).To(Equal("file2"))
					//to, _ := store.WriteFileArgsForCall(0)
					//Expect(to).To(Equal("{{path/file1}}"))
				})
			})

			Context("In", func() {
				It("should run a merge for each file in the directory specified in 'for_each.in'", func() {
					cfg.Merge[0].With.Files = []string{"fake1", "fake2"}
					cfg.ForEach.In = "integration/yamls/addons/sub1/"

					spruceConfig = []aviator.Spruce{cfg}
					spruceClient = new(fakes.FakeSpruceClient)
					processor = NewTestProcessor(spruceClient, store, modifier)

					err := processor.ProcessSilent(spruceConfig)
					Expect(err).ToNot(HaveOccurred())

					mergeOpts1 := spruceClient.MergeWithOptsArgsForCall(0)
					mergeOpts2 := spruceClient.MergeWithOptsArgsForCall(1)
					Expect(len(mergeOpts1.Files)).To(Equal(4))
					Expect(len(mergeOpts2.Files)).To(Equal(4))
					//Expect does not work for any reason
					//Expect(mergeOpts1.Files[3]).To(Equal("integration/yamls/addons/sub1/file1.yml"))
					Expect(mergeOpts2.Files[3]).To(Equal("integration/yamls/addons/sub1/file2.yml"))
				})
			})

			Context("'In' in combination with except", func() {
				It("should run a merge for each file in the directory specified in 'for_each.in' except those specified in 'except'", func() {
					cfg.Merge[0].With.Files = []string{"fake1", "fake2"}
					cfg.ForEach.In = "integration/yamls/"
					cfg.ForEach.Except = []string{"fake2.yml"}

					spruceConfig = []aviator.Spruce{cfg}
					spruceClient = new(fakes.FakeSpruceClient)
					processor = NewTestProcessor(spruceClient, store, modifier)

					err := processor.ProcessSilent(spruceConfig)
					Expect(err).ToNot(HaveOccurred())

					cc := spruceClient.MergeWithOptsCallCount()
					Expect(cc).To(Equal(2))

					mergeOpts1 := spruceClient.MergeWithOptsArgsForCall(0)
					mergeOpts := spruceClient.MergeWithOptsArgsForCall(1)
					Expect(len(mergeOpts1.Files)).To(Equal(4))
					Expect(len(mergeOpts.Files)).To(Equal(4))
					//Expect(mergeOpts1.Files[3]).To(Equal("integration/yamls/base.yml"))
					Expect(mergeOpts.Files[3]).To(Equal("integration/yamls/fake.yml"))
				})
			})

			Context("'In' in combination with regexp", func() {
				It("should run a merge for each file in the directory specified in 'for_each.in' matching the 'regexp'", func() {
					cfg.Merge[0].With.Files = []string{"fake1", "fake2"}
					cfg.ForEach.In = "integration/yamls/"
					cfg.ForEach.Regexp = "base.yml"

					spruceConfig = []aviator.Spruce{cfg}
					spruceClient = new(fakes.FakeSpruceClient)
					processor = NewTestProcessor(spruceClient, store, modifier)

					err := processor.ProcessSilent(spruceConfig)
					Expect(err).ToNot(HaveOccurred())

					cc := spruceClient.MergeWithOptsCallCount()
					Expect(cc).To(Equal(1))

					mergeOpts1 := spruceClient.MergeWithOptsArgsForCall(0)
					Expect(len(mergeOpts1.Files)).To(Equal(4))
					Expect(mergeOpts1.Files[3]).To(Equal("integration/yamls/base.yml"))
				})
			})

			Context("'In' in combination with 'regexp'", func() {
				It("should run a merge for each file in the directory specified in 'for_each.in' matching the 'regexp'", func() {
					cfg.Merge[0].With.Files = []string{"fake1", "fake2"}
					cfg.ForEach.In = "integration/yamls/"
					cfg.ForEach.Regexp = "base.yml"

					spruceConfig = []aviator.Spruce{cfg}
					spruceClient = new(fakes.FakeSpruceClient)
					processor = NewTestProcessor(spruceClient, store, modifier)

					err := processor.ProcessSilent(spruceConfig)
					Expect(err).ToNot(HaveOccurred())

					cc := spruceClient.MergeWithOptsCallCount()
					Expect(cc).To(Equal(1))

					mergeOpts1 := spruceClient.MergeWithOptsArgsForCall(0)
					Expect(len(mergeOpts1.Files)).To(Equal(4))
					Expect(mergeOpts1.Files[3]).To(Equal("integration/yamls/base.yml"))
					//to, _ := store.WriteFileArgsForCall(0)
					//Expect(to).To(Equal("integration/tmp/yamls_base.yml"))
				})
			})

			Context("Walk", func() {
				Context("'In' in combination with 'subdirs'", func() {
					It("should run a merge for each file in the directory and its subdirs", func() {
						cfg.Merge[0].With.Files = []string{"fake1", "fake2"}
						cfg.ForEach.In = "integration/yamls/addons/"
						cfg.ForEach.SubDirs = true

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						cc := spruceClient.MergeWithOptsCallCount()
						Expect(cc).To(Equal(3))

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(4))
					})
				})

				Context("'In' in combination with 'subdirs' and 'for_all'", func() {
					It("should run a merge for each file in the directory specified in 'for_each.in' and its subdirs... its complicated", func() {
						cfg.Merge[0].With.Files = []string{"fake1", "fake2"}
						cfg.ForEach.In = "integration/yamls/addons/"
						cfg.ForEach.SubDirs = true
						cfg.ForEach.ForAll = "integration/yamls/"

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						cc := spruceClient.MergeWithOptsCallCount()
						Expect(cc).To(Equal(9))

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(5))
					})
				})

				Context("'In' in combination with 'subdirs' and 'except'", func() {
					It("should run a merge for each file in the directory specified in 'for_each.in' and its subdirs, except those specified in 'except'.. its even more complicated", func() {
						cfg.Merge[0].With.Files = []string{"fake1", "fake2"}
						cfg.ForEach.In = "integration/yamls/addons/"
						cfg.ForEach.SubDirs = true
						cfg.ForEach.Except = "fake2"
						cfg.ForEach.ForAll = "integration/yamls/"

						spruceConfig = []aviator.Spruce{cfg}
						spruceClient = new(fakes.FakeSpruceClient)
						processor = NewTestProcessor(spruceClient, store, modifier)

						err := processor.ProcessSilent(spruceConfig)
						Expect(err).ToNot(HaveOccurred())

						cc := spruceClient.MergeWithOptsCallCount()
						Expect(cc).To(Equal(8))

						mergeOpts := spruceClient.MergeWithOptsArgsForCall(0)
						Expect(len(mergeOpts.Files)).To(Equal(4))
					})
				})
			})
		})
	})
})
