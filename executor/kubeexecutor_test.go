package executor_test

import (
	fakesrunner "code.cloudfoundry.org/commandrunner/fake_command_runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/JulzDiverse/aviator"
	. "github.com/JulzDiverse/aviator/executor"
)

var _ = Describe("Kubeexecutor", func() {

	var (
		kubeExec *KubeExecutor
		kubeCtl  aviator.Kube
		runner   *fakesrunner.FakeCommandRunner
		args     []string
		err      error
	)

	Context("For a given kubectl apply config", func() {

		JustBeforeEach(func() {
			runner = new(fakesrunner.FakeCommandRunner)
			kubeExec = NewKubeExecutorWithCustomRunner(runner)
			err = kubeExec.Execute(kubeCtl)
			cmds := runner.ExecutedCommands()
			args = cmds[0].Args
		})

		Context("with only a file to apply", func() {

			BeforeEach(func() {
				kubeCtl = aviator.Kube{
					aviator.KubeApply{
						File: "kube.yaml",
					},
				}
			})

			It("shouldn't error when executing the command", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should apply the given file", func() {
				Expect(args).To(ContainElement("kube.yaml"))
			})

			It("should call the runnter with no additional commands", func() {
				Expect(args).To(HaveLen(4))
			})
		})

		Context("When 'force' is set to true", func() {

			BeforeEach(func() {
				kubeCtl = aviator.Kube{
					aviator.KubeApply{
						File:  "kube.yaml",
						Force: true,
					},
				}
			})

			It("should add the 'force' flag to the exec call", func() {
				Expect(args).To(ContainElement("--force"))
			})
		})

		Context("When 'dry_run' is set to true", func() {

			BeforeEach(func() {
				kubeCtl = aviator.Kube{
					aviator.KubeApply{
						File:   "kube.yaml",
						DryRun: true,
					},
				}
			})

			It("should add the 'dry-run' flag to the exec call", func() {
				Expect(args).To(ContainElement("--dry-run"))
			})
		})

		Context("When 'recursive' is set to true", func() {

			BeforeEach(func() {
				kubeCtl = aviator.Kube{
					aviator.KubeApply{
						File:      "kube.yaml",
						Recursive: true,
					},
				}
			})

			It("should add the '--recursive' flag to the kubectl call", func() {
				Expect(args).To(ContainElement("--recursive"))
			})
		})

		Context("When 'overwrite' is set to true", func() {

			BeforeEach(func() {
				kubeCtl = aviator.Kube{
					aviator.KubeApply{
						File:      "kube.yaml",
						Overwrite: true,
					},
				}
			})

			It("should add the '--overwrite' flag to the kubectl call", func() {
				Expect(args).To(ContainElement("--overwrite"))
			})
		})
	})
})
