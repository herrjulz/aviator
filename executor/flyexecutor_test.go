package executor_test

import (
	"os/exec"

	"github.com/JulzDiverse/aviator"
	. "github.com/JulzDiverse/aviator/executor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flyexecutor", func() {
	var (
		flyExecutor *FlyExecutor
		fly         aviator.Fly
		args        []string
		exposeArgs  []string
		cmds        []*exec.Cmd
		err         error
	)

	Context("When generating commands", func() {
		JustBeforeEach(func() {
			flyExecutor = &FlyExecutor{}
			cmds, err = flyExecutor.Command(fly)
			args = cmds[0].Args
			exposeArgs = cmds[1].Args
		})

		Context("for a given fly config", func() {
			BeforeEach(func() {
				fly = aviator.Fly{
					Name:       "pipeline-name",
					Target:     "target-name",
					Config:     "pipeline.yml",
					CheckCreds: true,
					Expose:     true,
					TeamName:   "team-name",
					Vars:       []string{"credentials.yml", "props.yml"},
				}
			})

			It("shouldn't error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should generate two commands", func() {
				Expect(cmds).To(HaveLen(2))
			})

			It("generates the set-pipeline command with the 'target' flag and the right target", func() {
				Expect(args).To(ContainElement("--target"))
				Expect(args).To(ContainElement("target-name"))
			})

			It("generates the set-pipeline command with the 'team' flag and the right target", func() {
				Expect(args).To(ContainElement("--team"))
				Expect(args).To(ContainElement("team-name"))
			})

			It("generates the set-pipeline command with the 'pipeline' flag and the right pipeline", func() {
				Expect(args).To(ContainElement("--pipeline"))
				Expect(args).To(ContainElement("pipeline-name"))
			})

			It("generates the set-pipeline command with the 'config' flag and the right config file", func() {
				Expect(args).To(ContainElement("--config"))
				Expect(args).To(ContainElement("pipeline.yml"))
			})

			It("generates the set-pipeline command with the 'load-vars-from' flag and the right files", func() {
				Expect(args).To(ContainElement("--load-vars-from"))
				Expect(args).To(ContainElement("credentials.yml"))
				Expect(args).To(ContainElement("props.yml"))
			})

			It("generates the set-pipeline command including the '--check-creds' flag", func() {
				Expect(args).To(ContainElement("--check-creds"))
			})

			It("should create the expose command with pipeline name", func() {
				Expect(exposeArgs).To(ContainElement("expose-pipeline"))
				Expect(exposeArgs).To(ContainElement("--target"))
				Expect(exposeArgs).To(ContainElement("target-name"))
				Expect(exposeArgs).To(ContainElement("--pipeline"))
				Expect(exposeArgs).To(ContainElement("pipeline-name"))
			})
		})

		Context("When expose is not set (or false)", func() {
			BeforeEach(func() {
				fly = aviator.Fly{
					Name:   "pipeline-name",
					Target: "target-name",
					Config: "pipeline.yml",
				}
			})

			It("should generate two commands", func() {
				Expect(cmds).To(HaveLen(2))
			})

			It("should generate a hide-pipeline command", func() {
				Expect(exposeArgs).To(ContainElement("hide-pipeline"))
			})

			It("should add the --target flag following by the right target", func() {
				Expect(exposeArgs).To(ContainElement("--target"))
				Expect(exposeArgs).To(ContainElement("target-name"))
			})

			It("should add the pipeline flag following by the right pipeline name", func() {
				Expect(exposeArgs).To(ContainElement("--pipeline"))
				Expect(exposeArgs).To(ContainElement("pipeline-name"))
			})
		})
	})
})
