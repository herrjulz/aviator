package executor_test

import (
	"strings"

	fakesrunner "code.cloudfoundry.org/commandrunner/fake_command_runner"
	"github.com/JulzDiverse/aviator"
	. "github.com/JulzDiverse/aviator/executor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flyexecutor", func() {
	var flyExecutor *FlyExecutor
	var runner *fakesrunner.FakeCommandRunner
	var fly aviator.Fly

	BeforeEach(func() {
		runner = fakesrunner.New()
		flyExecutor = NewFlyExecutorWithCustomRunner(runner)
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
			It("calls the runner with the right commands", func() {
				err := flyExecutor.ExecuteWithCustomRunner(fly)
				Expect(err).ToNot(HaveOccurred())

				cmds := runner.ExecutedCommands()
				args := cmds[0].Args
				argsExpose := cmds[1].Args
				argsString := strings.Join(args, " ")
				Expect(argsString).To(Equal(
					"fly -t target-name set-pipeline -p pipeline-name -c pipeline.yml -l credentials.yml -l props.yml",
				))
				argsExposeString := strings.Join(argsExpose, " ")
				Expect(argsExposeString).To(Equal("fly -t target-name expose-pipeline -p pipeline-name"))
			})
		})

	})
})
