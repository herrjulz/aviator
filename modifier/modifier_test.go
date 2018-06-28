package modifier_test

import (
	"github.com/JulzDiverse/aviator"
	fakes "github.com/JulzDiverse/aviator/aviatorfakes"
	. "github.com/JulzDiverse/aviator/modifier"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Modifier", func() {

	var goml *fakes.FakeGomlClient
	var mod aviator.Modify
	var modifier *Modifier

	JustBeforeEach(func() {
		goml = new(fakes.FakeGomlClient)
		modifier = NewModifier(goml)
	})

	Context("Modify", func() {
		Context("Delete", func() {
			BeforeEach(func() {
				mod = aviator.Modify{}
			})

			It("should call Delete with the right deletion string", func() {
				mod.Delete = []string{"some.yaml.path"}
				_, err := modifier.Modify([]byte(`test`), mod)
				Expect(err).ToNot(HaveOccurred())
				_, path := goml.DeleteArgsForCall(0)
				Expect(path).To(Equal("some.yaml.path"))
			})

			PIt("should return an error when passing an empty string", func() {
				mod.Delete = []string{""}
				_, err := modifier.Modify([]byte(`test`), mod)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("modification path not provided")))
			})
		})

		Context("Set", func() {
			BeforeEach(func() {
				mod = aviator.Modify{
					Set: []aviator.PathVal{aviator.PathVal{}},
				}
			})

			It("should call Set with the provided path", func() {
				mod.Set[0].Path = "some.yaml.path"
				mod.Set[0].Value = "val"
				_, err := modifier.Modify([]byte(`test`), mod)
				Expect(err).ToNot(HaveOccurred())
				_, path, val := goml.SetArgsForCall(0)
				Expect(path).To(Equal("some.yaml.path"))
				Expect(val).To(Equal("val"))
			})

			PIt("should return an error when passing an empty string", func() {
				mod.Set[0].Path = ""
				_, err := modifier.Modify([]byte(`test`), mod)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("modification path not provided")))
			})
		})

		Context("Update", func() {
			BeforeEach(func() {
				mod = aviator.Modify{
					Update: []aviator.PathVal{aviator.PathVal{}},
				}
			})

			It("should call Update with the provided path", func() {
				mod.Update[0].Path = "some.yaml.path"
				mod.Update[0].Value = "val"
				_, err := modifier.Modify([]byte(`test`), mod)
				Expect(err).ToNot(HaveOccurred())
				_, path, val := goml.UpdateArgsForCall(0)
				Expect(path).To(Equal("some.yaml.path"))
				Expect(val).To(Equal("val"))
			})

			PIt("should return an error when passing an empty string", func() {
				mod.Update[0].Path = ""
				_, err := modifier.Modify([]byte(`test`), mod)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("modification path not provided")))
			})
		})
	})
})
