package filemanager_test

import (
	. "github.com/JulzDiverse/aviator/filemanager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filemanager", func() {

	var store *FileManager
	var allowCurlyBraces bool

	BeforeEach(func() {
		allowCurlyBraces = true
	})

	JustBeforeEach(func() {
		store = Store(allowCurlyBraces)
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

	Context("overwrting a file", func() {
		It("is successful", func() {
			err := store.WriteFile("{{keyD}}", []byte("content D"))
			Expect(err).ToNot(HaveOccurred())

			err = store.WriteFile("{{keyD}}", []byte("content D"))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("ReadFile", func() {
		It("reads a existing file from filesystem", func() {
			file, ok := store.ReadFile("integration/fake.yml")
			Expect(ok).To(Equal(true))
			Expect(string(file)).To(ContainSubstring("test:"))
		})
	})

	//Context("WriteFile", func() {
	//It("create non existing dirs", func() {
	//err := store.WriteFile("integration/non/existing/fake.yml", []byte("file"))
	//Expect(err).ToNot(HaveOccurred())
	//})
	//})

	Context("When double curly braces are not allowed", func() {
		BeforeEach(func() {
			allowCurlyBraces = false
		})

		It("doesn't unquote curly braces on write", func() {
			err := store.WriteFile("{{keyE}}", []byte("{{content E}}"))
			Expect(err).ToNot(HaveOccurred())
			file, ok := store.ReadFile("{{keyE}}")
			Expect(ok).To(Equal(true))
			Expect(string(file)).To(ContainSubstring("{{content E}}"))
		})
	})
})
