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

	Context("Set/GetFile", func() {
		Context("When setting a file", func() {
			It("is available with GetFile", func() {
				store.SetFile("key", []byte("content"))
				file, ok := store.GetFile("key")
				Expect(ok).To(Equal(true))
				Expect(string(file[:len(file)])).To(Equal("content"))
			})

			It("is available using GetFile with a key inside double curly braces", func() {
				store.SetFile("keyA", []byte("content A"))
				file, ok := store.GetFile("{{keyA}}")
				Expect(ok).To(Equal(true))
				Expect(string(file[:len(file)])).To(Equal("content A"))
			})
		})

		Context("When setting a file with a key in double curly braces", func() {
			It("is available via GetFile", func() {
				store.SetFile("{{keyB}}", []byte("content B"))
				file, ok := store.GetFile("keyB")
				Expect(ok).To(Equal(true))
				Expect(string(file[:len(file)])).To(Equal("content B"))
			})
		})

		Context("When setting a file with a key in double curly braces", func() {
			It("is available via GetFile with the key in double curly braces", func() {
				err := store.SetFile("{{keyC}}", []byte("content C"))
				Expect(err).ToNot(HaveOccurred())
				file, ok := store.GetFile("{{keyC}}")
				Expect(ok).To(Equal(true))
				Expect(string(file[:len(file)])).To(Equal("content C"))
			})
		})
	})

	Context("Setting a file that already existsr", func() {
		It("returns an error", func() {
			err := store.SetFile("keyD", []byte("content D"))
			Expect(err).ToNot(HaveOccurred())

			err = store.SetFile("keyD", []byte("content D"))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("file keyD in virtual filestore already exists")))
		})
	})
})
