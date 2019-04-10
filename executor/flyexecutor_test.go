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
		cmd         *exec.Cmd
		err         error
	)

	JustBeforeEach(func() {
		flyExecutor = &FlyExecutor{}
		cmd, err = flyExecutor.Command(fly)
		args = cmd.Args
	})

	Context("Execute", func() {
		Context("for a given fly config", func() {
			BeforeEach(func() {
				fly = aviator.Fly{
					Name:   "pipeline-name",
					Target: "target-name",
					Config: "pipeline.yml",
					Expose: true,
					Vars:   []string{"credentials.yml", "props.yml"},
				}
			})

			It("shouldn't error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("calls the runner with the right commands", func() {
				Expect(args).To(ContainElement("target-name"))
				Expect(args).To(ContainElement("pipeline-name"))
				Expect(args).To(ContainElement("pipeline.yml"))
				Expect(args).To(ContainElement("credentials.yml"))
				Expect(args).To(ContainElement("props.yml"))
			})
		})
	})
})
