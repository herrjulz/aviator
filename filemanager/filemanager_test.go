package filemanager_test

import (
	. "github.com/JulzDiverse/aviator/filemanager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filemanager", func() {

	var store *FileStore

	BeforeEach(func() {
		store = Store()
	})

	Context("Write/ReadFile", func() {
		Context("When setting a file", func() {
			It("is available with GetFile", func() {
				store.WriteFile("{{key}}", []byte("content"))
				file, ok := store.ReadFile("{{key}}")
				Expect(ok).To(Equal(true))
				Expect(string(file[:len(file)])).To(Equal("content"))
			})
		})

		Context("When setting a file with a key in double curly braces", func() {
			It("is also available with the key without curly braces", func() {
				store.WriteFile("{{keyB}}", []byte("content B"))
				file, ok := store.ReadFile("keyB")
				Expect(ok).To(Equal(true))
				Expect(string(file[:len(file)])).To(Equal("content B"))
			})
		})

		Context("When setting a file with a key in double curly braces", func() {
			It("is available via ReadFile with the key in double curly braces", func() {
				err := store.WriteFile("{{keyC}}", []byte("content C"))
				Expect(err).ToNot(HaveOccurred())
				file, ok := store.ReadFile("{{keyC}}")
				Expect(ok).To(Equal(true))
				Expect(string(file[:len(file)])).To(Equal("content C"))
			})
		})
	})

	Context("Setting a file that already exists", func() {
		It("returns an error", func() {
			err := store.WriteFile("{{keyD}}", []byte("content D"))
			Expect(err).ToNot(HaveOccurred())

			err = store.WriteFile("{{keyD}}", []byte("content D"))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("file keyD in virtual filestore already exists")))
		})
	})

	Context("ReadFile", func() {
		It("reads a existing file from filesystem", func() {
			file, ok := store.ReadFile("integration/fake.yml")
			Expect(ok).To(Equal(true))
			Expect(string(file)).To(ContainSubstring("test:"))
		})
	})
})
