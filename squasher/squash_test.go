package squasher_test

import (
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/JulzDiverse/aviator/squasher"
)

var _ = Describe("Squash", func() {

	var (
		yamls    [][]byte
		squashed []byte
		err      error
	)

	JustBeforeEach(func() {
		squashed, err = Squash(yamls)
	})

	Context("When providing multiple yaml files", func() {
		BeforeEach(func() {
			yaml1 := `---
i_am_yaml: 1`

			yaml2 := `---
i_am_yaml: 2`

			yamls = [][]byte{[]byte(yaml1), []byte(yaml2)}
		})

		It("should not fail", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("squashes it to a single yaml file", func() {
			expected := `---
i_am_yaml: 1
---
i_am_yaml: 2`

			Expect(strings.Trim(string(squashed), "\n")).To(Equal(strings.Trim(string(expected), "\n")))

			ioutil.WriteFile("check", squashed, 0644)
		})
	})

	Context("When providing less than two yaml files", func() {
		Context("When providing a single yaml file", func() {
			BeforeEach(func() {
				yamls = [][]byte{{}}
				squashed, err = Squash(yamls)
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should not try to squash anything", func() {
				Expect(squashed).To(BeNil())
			})
		})

		Context("When providing no yaml file", func() {
			BeforeEach(func() {
				yamls = [][]byte{}
				squashed, err = Squash(yamls)
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should not try to squash anything", func() {
				Expect(squashed).To(BeNil())
			})
		})
	})
})
