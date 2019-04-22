package executor_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/JulzDiverse/aviator"
	. "github.com/JulzDiverse/aviator/executor"
)

var _ = Describe("Genericexecutor", func() {

	var (
		genExec GenericExecutor
		err     error
		cfg     []aviator.Executable
		cmds    []*exec.Cmd
	)

	JustBeforeEach(func() {
		cmds, err = genExec.Command(cfg)
	})

	Context("When only arguments are provided", func() {
		BeforeEach(func() {
			cfg = []aviator.Executable{
				{
					Executable: "cp",
					Args: []string{
						"file", "destination/",
					},
				},
			}
		})

		It("shouldn't fail", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not contain any other flags or commands", func() {
			Expect(stringifyCmd(cmds[0])).To(Equal("cp file destination/"))
		})
	})

	Context("When global options are provided", func() {
		Context("and the option is a bool", func() {

			BeforeEach(func() {
				cfg = []aviator.Executable{
					{
						Executable: "cp",
						GlobalOptions: []aviator.Option{
							{
								Name: "-R",
							},
						},
						Args: []string{
							"file", "destination/",
						},
					},
				}
			})

			Context("When a single option is provided", func() {
				It("shouldn't fail", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should include the gloabl option", func() {
					Expect(stringifyCmd(cmds[0])).To(Equal("cp -R file destination/"))
				})
			})

			Context("When multiple options are provided", func() {
				BeforeEach(func() {
					cfg[0].GlobalOptions = append(cfg[0].GlobalOptions, aviator.Option{
						Name: "-H",
					})
				})

				It("shouldn't fail", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should include multiple gloabl options", func() {
					Expect(stringifyCmd(cmds[0])).To(Equal("cp -R -H file destination/"))
				})
			})
		})

		Context("and the option is not a bool", func() {

			BeforeEach(func() {
				cfg = []aviator.Executable{
					{

						Executable: "exec",
						GlobalOptions: []aviator.Option{
							{
								Name:  "--global-option",
								Value: "glob-value",
							},
						},
						Args: []string{
							"arg",
						},
					},
				}
			})

			Context("When a single option is provided", func() {
				It("shouldn't fail", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should include the gloabl option including the value", func() {
					Expect(stringifyCmd(cmds[0])).To(Equal("exec --global-option glob-value arg"))
				})
			})

			Context("When multiple options are provided", func() {
				BeforeEach(func() {
					cfg[0].GlobalOptions = append(cfg[0].GlobalOptions, aviator.Option{
						Name:  "--another-global-option",
						Value: "another-global-value",
					})
				})

				It("shouldn't fail", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should include multiple gloabl options", func() {
					Expect(stringifyCmd(cmds[0])).To(Equal("exec --global-option glob-value --another-global-option another-global-value arg"))
				})
			})
		})
	})

	Context("When a command is provided", func() {
		BeforeEach(func() {
			cfg = []aviator.Executable{
				{
					Executable: "executable",
					Command: aviator.Command{
						Name: "command",
						Options: []aviator.Option{
							{
								Name:  "--command-option",
								Value: "option",
							},
						},
					},
					Args: []string{
						"arg",
					},
				},
			}
		})

		It("should include the command", func() {
			Expect(stringifyCmd(cmds[0])).To(Equal("executable command --command-option option arg"))
		})

		Context("and a command-option is a bool", func() {
			BeforeEach(func() {
				cfg[0].Command.Options = append(cfg[0].Command.Options, aviator.Option{
					Name: "--bool-option",
				})
			})

			It("shouldn't fail", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should include the bool option without value", func() {
				Expect(stringifyCmd(cmds[0])).To(Equal("executable command --command-option option --bool-option arg"))
			})
		})
	})
})

func stringifyCmd(cmd *exec.Cmd) string {
	result := ""
	result = fmt.Sprintf("%s", cmd.Args[0])
	for i := 1; i < len(cmd.Args); i++ {
		result = fmt.Sprintf("%s %s", result, cmd.Args[i])
	}
	return result
}
