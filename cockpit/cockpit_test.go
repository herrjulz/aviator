package cockpit_test

import (
	"errors"
	"os"

	fakes "github.com/JulzDiverse/aviator/aviatorfakes"
	. "github.com/JulzDiverse/aviator/cockpit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cockpit", func() {
	var aviatorYaml string
	var spruceProcessor *fakes.FakeSpruceProcessor
	var flyExecuter *fakes.FakeFlyExecuter
	var validator *fakes.FakeValidator
	var cockpit *Cockpit
	var vars map[string]string

	BeforeEach(func() {
		spruceProcessor = new(fakes.FakeSpruceProcessor)
		flyExecuter = new(fakes.FakeFlyExecuter)
		validator = new(fakes.FakeValidator)
		cockpit = Init(spruceProcessor, flyExecuter, validator)
		vars = map[string]string{}
	})

	Context("New", func() {
		Context("aviator.yml parsing", func() {
			var aviator *Aviator

			Context("Spruce", func() {
				Context("Merge Section", func() {
					It("is able to read all 'with' related properties", func() {
						aviatorYaml = `spruce:
- base: base.yml
  merge:
  - with:
      files:
      - file1.yml
      - file2.yml
      skip_non_existing: true
    regexp: ".*.(yml)"
    in_dir: /path/
  - with_in: path/to/dir/
    except:
    - file2.yml
  - with_all_in: path/
  to: result.yml`

						var err error
						aviator, err = cockpit.NewAviator([]byte(aviatorYaml), vars)
						Expect(err).ToNot(HaveOccurred())

						Expect(len(aviator.AviatorYaml.Spruce[0].Merge[0].With.Files)).To(Equal(2))
						Expect(aviator.AviatorYaml.Spruce[0].Merge[1].WithIn).To(Equal("path/to/dir/"))
						Expect(len(aviator.AviatorYaml.Spruce[0].Merge[1].Except)).To(Equal(1))
						Expect(aviator.AviatorYaml.Spruce[0].Merge[0].Regexp).To(Equal(".*.(yml)"))
						Expect(aviator.AviatorYaml.Spruce[0].Merge[0].With.Skip).To(Equal(true))
						Expect(aviator.AviatorYaml.Spruce[0].To).To(Equal("result.yml"))
						Expect(aviator.AviatorYaml.Spruce[0].Merge[2].WithAllIn).To(Equal("path/"))
					})

					It("is able to read all 'cherry_pick' and 'skip_eval' properties", func() {
						aviatorYaml = `spruce:
- base: some/tmp/dir/file1.yml
  cherry_pick:
  - one
  - two
  - three
  merge:
  - with_in: path/
  skip_eval: true
  to_dir: foo/bar/`

						var err error
						aviator, err = cockpit.NewAviator([]byte(aviatorYaml), vars)
						Expect(err).ToNot(HaveOccurred())

						Expect(len(aviator.AviatorYaml.Spruce[0].CherryPicks)).To(Equal(3))
						Expect(aviator.AviatorYaml.Spruce[0].SkipEval).To(Equal(true))
					})

					Context("Environment Variabels & Curly Braces", func() {
						It("is able resolve environment variables", func() {
							os.Setenv("ENV_VAR", "envVar")
							os.Setenv("ANOTHER_VAR", "another")
							os.Setenv("RESULT", "result")
							aviatorYaml = `spruce:
- base: $ENV_VAR
  merge:
  - with:
      files:
      - $ANOTHER_VAR
  to: $RESULT`

							var err error
							aviator, err = cockpit.NewAviator([]byte(aviatorYaml), vars)
							Expect(err).ToNot(HaveOccurred())
							Expect(aviator.AviatorYaml.Spruce[0].Base).To(Equal("envVar"))
							Expect(aviator.AviatorYaml.Spruce[0].Merge[0].With.Files[0]).To(Equal("another"))
							Expect(aviator.AviatorYaml.Spruce[0].To).To(Equal("result"))
						})

						It("returns an error if variables are not set", func() {
							os.Setenv("ENV_VAR", "envVar")
							os.Setenv("ANOTHER_VAR", "another")
							os.Unsetenv("RESULT")
							aviatorYaml = `spruce:
- base: $ENV_VAR
  merge:
  - with:
      files:
      - $ANOTHER_VAR
  to: $RESULT`

							var err error
							aviator, err = cockpit.NewAviator([]byte(aviatorYaml), vars)
							Expect(err).To(HaveOccurred())
						})

						It("is able to parse '{{}}'", func() {
							os.Setenv("ENV_VAR", "envVar")
							os.Setenv("ANOTHER_VAR", "another")
							os.Setenv("RESULT", "result")
							aviatorYaml = `spruce:
- base: input.yml 
  merge:
  - with:
      files:
      - {{identifier}}
  to: {{result}}`

							var err error
							aviator, err = cockpit.NewAviator([]byte(aviatorYaml), vars)
							Expect(err).ToNot(HaveOccurred())
						})
					})
				})

				Context("ForEach Section", func() {
					It("is able to parse all for_each.files related properties", func() {
						aviatorYaml = `spruce:
- base: result.yml
  merge:
  - with_in: another/path/
  for_each:
    files:
    - file1.yml
    - file2.yml
    skip_non_existing: true
    in_dir: path/
  to_dir: some/tmp/dir/`

						var err error
						aviator, err = cockpit.NewAviator([]byte(aviatorYaml), vars)
						Expect(err).ToNot(HaveOccurred())

						Expect(len(aviator.AviatorYaml.Spruce[0].ForEach.Files)).To(Equal(2))
						Expect(aviator.AviatorYaml.Spruce[0].ForEach.Skip).To(Equal(true))
						Expect(aviator.AviatorYaml.Spruce[0].ForEach.InDir).To(Equal("path/"))
						Expect(aviator.AviatorYaml.Spruce[0].ToDir).To(Equal("some/tmp/dir/"))
					})

					It("is able to parse all for_each.in related properties", func() {
						aviatorYaml = `spruce:
- base: result.yml
  merge:
  - with_in: another/path/
  for_each:
    in: path/
    except:
    - file.yml
    - file2.yml
    regexp: ".*.(yml)"
  to_dir: some/tmp/dir/`

						var err error
						aviator, err = cockpit.NewAviator([]byte(aviatorYaml), vars)
						Expect(err).ToNot(HaveOccurred())

						Expect(aviator.AviatorYaml.Spruce[0].ForEach.In).To(Equal("path/"))
						Expect(len(aviator.AviatorYaml.Spruce[0].ForEach.Except)).To(Equal(2))
						Expect(aviator.AviatorYaml.Spruce[0].ForEach.Regexp).To(Equal(".*.(yml)"))
					})

					It("is able to parse all for_each.in related properties in combination with enable_sub_dirs", func() {
						aviatorYaml = `spruce:
- base: result.yml
  merge:
  - with_in: another/path/
  for_each:
    in: path/
    include_sub_dirs: true
    copy_parents: true
    enable_matching: true
  to_dir: some/tmp/dir/`

						var err error
						aviator, err = cockpit.NewAviator([]byte(aviatorYaml), vars)
						Expect(err).ToNot(HaveOccurred())

						Expect(aviator.AviatorYaml.Spruce[0].ForEach.SubDirs).To(Equal(true))
						Expect(aviator.AviatorYaml.Spruce[0].ForEach.CopyParents).To(Equal(true))
						Expect(aviator.AviatorYaml.Spruce[0].ForEach.EnableMatching).To(Equal(true))
					})
				})
			})

			Context("fly section", func() {
				BeforeEach(func() {
					aviatorYaml = `fly:
  name: pipelineName
  target: targetName
  config: configFile
  expose: true
  load_vars_from:
  - credentials.yml
  vars:
    key: value`
				})

				It("is able to read all properties from the fly section", func() {
					var err error
					aviator, err = cockpit.NewAviator([]byte(aviatorYaml), vars)
					Expect(err).ToNot(HaveOccurred())

					Expect(aviator.AviatorYaml.Fly.Name).To(Equal("pipelineName"))
					Expect(aviator.AviatorYaml.Fly.Target).To(Equal("targetName"))
					Expect(aviator.AviatorYaml.Fly.Config).To(Equal("configFile"))
					Expect(aviator.AviatorYaml.Fly.Expose).To(BeTrue())
					Expect(len(aviator.AviatorYaml.Fly.Vars)).To(Equal(1))
					Expect(aviator.AviatorYaml.Fly.Var["key"]).To(Equal("value"))
				})
			})
		})
	})

	Context("spruce section processor", func() {
		var aviator *Aviator
		BeforeEach(func() {
			aviatorYaml = `spruce:
- base: input.yml
  merge:
  - with_in: some/dir/
  to: output.yml`

			spruceProcessor.ProcessWithOptsReturns(errors.New("uups"))
		})

		It("returns a valid error message", func() {
			var err error
			aviator, err = cockpit.NewAviator([]byte(aviatorYaml), vars)
			Expect(err).ToNot(HaveOccurred())
			err = aviator.ProcessSprucePlan(false, false)

			Expect(err).To(MatchError(ContainSubstring("Processing Spruce Plan FAILED")))
			Expect(err).To(MatchError(ContainSubstring("uups")))
		})
	})
})
