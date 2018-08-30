package evaluator_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/JulzDiverse/aviator/evaluator"
)

var _ = Describe("Evaluate", func() {
	Context("When evaluating an aviator file variable", func() {

		var (
			aviatorYaml  string
			expectedYaml string
			evaluated    []byte
			err          error
			vars         map[string]string
		)

		JustBeforeEach(func() {
			evaluated, err = Evaluate([]byte(aviatorYaml), vars)
		})

		Context("When the varialbe key is an existing key", func() {
			Context("When the value is a usual value", func() {
				BeforeEach(func() {
					vars = map[string]string{
						"base_yaml": "base.yml",
					}

					aviatorYaml = `---
spruce:
- base: (( base_yaml )).old`

					expectedYaml = `---
spruce:
- base: base.yml.old`
				})

				It("should replace the key with it's provided value", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(strings.TrimSpace(string(evaluated))).To(Equal(expectedYaml))
				})
			})

			Context("When the value is a mulit-line value", func() {
				BeforeEach(func() {
					vars = map[string]string{
						"multi_line": `hello
world`,
					}

					aviatorYaml = `---
fly:
  vars:
    key: (( multi_line ))`

					expectedYaml = `---
fly:
  vars:
    key: "hello\nworld"`

				})

				It("should replace the key with it's provided value", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(strings.TrimSpace(string(evaluated))).To(Equal(expectedYaml))
				})
			})
		})

		Context("When the variable key is not existing", func() {
			BeforeEach(func() {
				aviatorYaml = `---
key: (( not_provided ))`
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should return a meaningful error message", func() {
				Expect(err).To(MatchError(ContainSubstring("Variable (( not_provided )) not provided")))
			})
		})
	})
})
