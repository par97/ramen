package suites

import (
	"fmt"

	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/dractions"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

type AppSetSuite struct {
	w   workloads.Workload
	d   deployers.Deployer
	r   dractions.DRActions
	Ctx *util.TestContext
}

func (s *AppSetSuite) SetContext(ctx *util.TestContext) {
	ctx.Log.Info("enter AppSetSuite SetContext")
	s.Ctx = ctx
}

func (s *AppSetSuite) SetupSuite() error {
	s.Ctx.Log.Info("enter AppSetSuite SetupSuite")

	deployment := &workloads.Deployment{Ctx: s.Ctx}
	deployment.Init()
	s.w = deployment

	sub := &deployers.ApplicationSet{Ctx: s.Ctx}
	sub.Init(deployment)
	s.d = sub

	s.r = dractions.DRActions{Ctx: s.Ctx}
	return nil
}

func (s *AppSetSuite) TeardownSuite() error {
	s.Ctx.Log.Info("enter AppSetSuite TeardownSuite")
	return nil
}

func (s *AppSetSuite) Tests() []Test {
	s.Ctx.Log.Info("enter AppSetSuite Tests")
	return []Test{
		s.TestWorkloadDeployment,
		s.TestEnableProtection,
		s.TestWorkloadFailover,
		s.TestWorkloadRelocation,
		s.TestDisableProtection,
		s.TestWorkloadUndeployment,
	}
}

func (s *AppSetSuite) TestWorkloadDeployment() error {
	s.Ctx.Log.Info("enter AppSetSuite TestWorkloadDeployment")
	err := s.d.Deploy(s.w)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestWorkloadDeployment: Pass")
	return nil
}

func (s *AppSetSuite) TestEnableProtection() error {
	s.Ctx.Log.Info("enter AppSetSuite TestEnableProtection")
	err := s.r.EnableProtection(s.w, s.d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestEnableProtection: Pass")
	return nil
}

func (s *AppSetSuite) TestWorkloadFailover() error {
	s.Ctx.Log.Info("enter AppSetSuite TestWorkloadFailover")
	err := s.r.Failover(s.w, s.d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestWorkloadFailover: Pass")
	return nil
}

func (s *AppSetSuite) TestWorkloadRelocation() error {
	s.Ctx.Log.Info("enter AppSetSuite TestWorkloadRelocation")
	err := s.r.Relocate(s.w, s.d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestWorkloadRelocation: Pass")
	return nil
}

func (s *AppSetSuite) TestDisableProtection() error {
	s.Ctx.Log.Info("enter AppSetSuite TestDisableProtection")
	err := s.r.DisableProtection(s.w, s.d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestDisableProtection: Pass")
	return nil
}

func (s *AppSetSuite) TestWorkloadUndeployment() error {
	s.Ctx.Log.Info("enter AppSetSuite TestWorkloadUndeployment")
	err := s.d.Undeploy(s.w)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestWorkloadUndeployment: Pass")
	return nil
}
