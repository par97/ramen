package suites

import (
	"fmt"

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
	ctx.Log.Info("enter BasicSuite SetContext")
	s.Ctx = ctx
}

func (s *BasicSuite) SetupSuite() error {
	s.Ctx.Log.Info("enter BasicSuite SetupSuite")

	deployment := &workloads.Deployment{}
	deployment.Init()
	deployment.Ctx = s.Ctx
	s.w = deployment
	s.d = deployers.Subscription{Ctx: s.Ctx}
	s.r = dractions.DRActions{Ctx: s.Ctx}
	return nil
}

func (s *BasicSuite) TeardownSuite() error {
	s.Ctx.Log.Info("enter BasicSuite TeardownSuite")
	return nil
}

func (s *BasicSuite) Tests() []Test {
	s.Ctx.Log.Info("enter BasicSuite Tests")
	return []Test{
		s.TestWorkloadDeployment,
		//s.TestEnableProtection,
		//s.TestWorkloadFailover,
		//s.TestWorkloadRelocation,
		//s.TestDisableProtection,
		s.TestWorkloadUndeployment,
	}
}

func (s *BasicSuite) TestWorkloadDeployment() error {
	s.Ctx.Log.Info("enter BasicSuite TestWorkloadDeployment")
	return s.d.Deploy(s.w)
}

func (s *BasicSuite) TestEnableProtection() error {
	s.Ctx.Log.Info("enter BasicSuite TestEnableProtection")
	err := s.r.EnableProtection(s.w, s.d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	s.Ctx.Log.Info("TestEnableProtection: Pass")
	return nil
}

func (s *BasicSuite) TestWorkloadFailover() error {
	s.Ctx.Log.Info("enter BasicSuite TestWorkloadFailover")
	return s.r.Failover(s.w, s.d)
}

func (s *BasicSuite) TestWorkloadRelocation() error {
	s.Ctx.Log.Info("enter BasicSuite TestWorkloadRelocation")
	return s.r.Relocate(s.w, s.d)
}

func (s *BasicSuite) TestDisableProtection() error {
	s.Ctx.Log.Info("enter BasicSuite TestDisableProtection")
	err := s.r.DisableProtection(s.w, s.d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	s.Ctx.Log.Info("TestDisableProtection: Pass")
	return nil
}

func (s *BasicSuite) TestWorkloadUndeployment() error {
	s.Ctx.Log.Info("enter BasicSuite TestWorkloadUndeployment")
	return s.d.Undeploy(s.w)
}
