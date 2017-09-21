package fake_command_runner_test

import (
	"code.cloudfoundry.org/commandrunner/fake_command_runner"

	"os/exec"

	"code.cloudfoundry.org/commandrunner/fake_command_runner/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FakeCommandRunner", func() {

	var (
		runner    *fake_command_runner.FakeCommandRunner
		cmd, cmd2 *exec.Cmd
	)

	BeforeEach(func() {
		runner = fake_command_runner.New()
		cmd = &exec.Cmd{Path: "p", Args: []string{"p", "arg"}}
		cmd2 = &exec.Cmd{Path: "q", Args: []string{"q", "arg2"}}
	})

	Describe("Kill", func() {
		It("should record Kill commands", func() {
			runner.Kill(cmd)
			Expect(runner.KilledCommands()).To(Equal([]*exec.Cmd{cmd}))
		})

		// This may seem like an odd test, but it exposed a bug.
		It("should not confuse Kill and Wait", func() {
			runner.Kill(cmd)
			runner.Wait(cmd2)
			Expect(runner.KilledCommands()).To(Equal([]*exec.Cmd{cmd}))
		})
	})

	Describe("Wait", func() {
		It("should record Wait commands", func() {
			runner.Wait(cmd)
			Expect(runner.WaitedCommands()).To(Equal([]*exec.Cmd{cmd}))
		})

		It("should not confuse Wait and Kill", func() {
			runner.Wait(cmd)
			runner.Kill(cmd2)
			Expect(runner.WaitedCommands()).To(Equal([]*exec.Cmd{cmd}))
		})
	})

	Describe("Matchers", func() {
		Describe("HaveExecuted", func() {
			It("should match commands in any order", func() {
				runner.Run(cmd)
				runner.Run(cmd2)
				Expect(runner).To(fake_command_runner_matchers.HaveExecuted(fake_command_runner.CommandSpec{Path: "p", Args: []string{"arg"}},
					fake_command_runner.CommandSpec{Path: "q", Args: []string{"arg2"}}))
				Expect(runner).To(fake_command_runner_matchers.HaveExecuted(fake_command_runner.CommandSpec{Path: "q", Args: []string{"arg2"}},
					fake_command_runner.CommandSpec{Path: "p", Args: []string{"arg"}}))
			})

			It("should match a subset of the commands", func() {
				runner.Run(cmd)
				runner.Run(cmd2)
				Expect(runner).To(fake_command_runner_matchers.HaveExecuted(fake_command_runner.CommandSpec{Path: "p", Args: []string{"arg"}}))
			})
		})

		Describe("HaveExecutedSerially", func() {
			It("should match commands in order", func() {
				runner.Run(cmd)
				runner.Run(cmd2)
				Expect(runner).To(fake_command_runner_matchers.HaveExecutedSerially(fake_command_runner.CommandSpec{Path: "p", Args: []string{"arg"}},
					fake_command_runner.CommandSpec{Path: "q", Args: []string{"arg2"}}))
				Expect(runner).NotTo(fake_command_runner_matchers.HaveExecutedSerially(fake_command_runner.CommandSpec{Path: "q", Args: []string{"arg2"}},
					fake_command_runner.CommandSpec{Path: "p", Args: []string{"arg"}}))
			})

			It("should match a subset of the commands", func() {
				runner.Run(cmd)
				runner.Run(cmd2)
				Expect(runner).To(fake_command_runner_matchers.HaveExecutedSerially(fake_command_runner.CommandSpec{Path: "p", Args: []string{"arg"}}))
			})
		})
	})
})
