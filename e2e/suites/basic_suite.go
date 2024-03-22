package suites

import (
	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/dractions"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

type BasicSuite struct {
	w   workloads.Workload
	d   deployers.Deployer
	r   dractions.DRActions
	Ctx *util.TestContext
}

func (s *BasicSuite) SetContext(ctx *util.TestContext) {
	ctx.Log.Info("enter SetContext")
	s.Ctx = ctx
}

func (s *BasicSuite) SetupSuite() error {
	s.Ctx.Log.Info("enter SetupSuite")
	s.w = workloads.Deployment{
		RepoURL:   "https://github.com/ramendr/ocm-ramen-samples.git",
		Path:      "subscription/deployment-k8s-regional-rbd",
		Revision:  "main",
		Ctx:       s.Ctx,
		Name:      "deployment-rbd",
		Namespace: "deployment-rbd",
		Dr_policy: "dr-policy",
		Pvc_label: "busybox",
	}
	s.d = deployers.Subscription{Ctx: s.Ctx}
	s.r = dractions.DRActions{Ctx: s.Ctx}
	return nil
}

func (s *BasicSuite) TeardownSuite() error {
	s.Ctx.Log.Info("enter TeardownSuite")
	return nil
}

func (s *BasicSuite) Tests() []Test {
	s.Ctx.Log.Info("enter Tests")
	return []Test{
		// s.TestWorkloadDeployment,
		s.TestEnableProtection,
		// s.TestWorkloadFailover,
		// s.TestWorkloadRelocation,
		// s.TestDisableProtection,
		// s.TestWorkloadUndeployment,
	}
}

func (s *BasicSuite) TestWorkloadDeployment() error {
	s.Ctx.Log.Info("enter TestWorkloadDeployment")
	return s.d.Deploy(s.w)
}

func (s *BasicSuite) TestEnableProtection() error {
	s.Ctx.Log.Info("enter TestEnableProtection")
	return s.r.EnableProtection(s.w, s.d)
}

func (s *BasicSuite) TestWorkloadFailover() error {
	s.Ctx.Log.Info("enter TestWorkloadFailover")
	return s.r.Failover(s.w, s.d)
}

func (s *BasicSuite) TestWorkloadRelocation() error {
	s.Ctx.Log.Info("enter TestWorkloadRelocation")
	return s.r.Relocate(s.w, s.d)
}

func (s *BasicSuite) TestDisableProtection() error {
	s.Ctx.Log.Info("enter TestDisableProtection")
	return nil
}

func (s *BasicSuite) TestWorkloadUndeployment() error {
	s.Ctx.Log.Info("enter TestWorkloadUndeployment")
	return s.d.Undeploy(s.w)
}
