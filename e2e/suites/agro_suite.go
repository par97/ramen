package suites

import (
	"fmt"

	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/dractions"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

type ArogSuite struct {
	w   workloads.Workload
	d   deployers.Deployer
	r   dractions.DRActions
	Ctx *util.TestContext
}

func (s *ArogSuite) SetContext(ctx *util.TestContext) {
	ctx.Log.Info("enter BasicSuite SetContext")
	s.Ctx = ctx
}

func (s *ArogSuite) SetupSuite() error {
	s.Ctx.Log.Info("enter BasicSuite SetupSuite")

	deployment := &workloads.Deployment{Ctx: s.Ctx}
	deployment.Init()
	s.w = deployment

	appset := &deployers.ApplicationSet{Ctx: s.Ctx}
	appset.Init()
	s.d = appset

	s.r = dractions.DRActions{Ctx: s.Ctx}
	return nil
}

func (s *ArogSuite) TeardownSuite() error {
	s.Ctx.Log.Info("enter BasicSuite TeardownSuite")
	return nil
}

func (s *ArogSuite) Tests() []Test {
	s.Ctx.Log.Info("enter BasicSuite Tests")
	return []Test{
		s.TestWorkloadDeployment,
		// s.TestEnableProtection,
		// s.TestWorkloadFailover,
		// s.TestWorkloadRelocation,
		// s.TestDisableProtection,
		// s.TestWorkloadUndeployment,
	}
}

func (s *ArogSuite) TestWorkloadDeployment() error {
	s.Ctx.Log.Info("enter BasicSuite TestWorkloadDeployment")
	err := s.d.Deploy(s.w)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestWorkloadDeployment: Pass")
	return nil
}

func (s *ArogSuite) TestEnableProtection() error {
	s.Ctx.Log.Info("enter BasicSuite TestEnableProtection")
	err := s.r.EnableProtection(s.w, s.d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestEnableProtection: Pass")
	return nil
}

func (s *ArogSuite) TestWorkloadFailover() error {
	s.Ctx.Log.Info("enter BasicSuite TestWorkloadFailover")
	err := s.r.Failover(s.w, s.d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestWorkloadFailover: Pass")
	return nil
}

func (s *ArogSuite) TestWorkloadRelocation() error {
	s.Ctx.Log.Info("enter BasicSuite TestWorkloadRelocation")
	err := s.r.Relocate(s.w, s.d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestWorkloadRelocation: Pass")
	return nil
}

func (s *ArogSuite) TestDisableProtection() error {
	s.Ctx.Log.Info("enter BasicSuite TestDisableProtection")
	err := s.r.DisableProtection(s.w, s.d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestDisableProtection: Pass")
	return nil
}

func (s *ArogSuite) TestWorkloadUndeployment() error {
	s.Ctx.Log.Info("enter BasicSuite TestWorkloadUndeployment")
	err := s.d.Undeploy(s.w)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.Ctx.Log.Info("TestWorkloadUndeployment: Pass")
	return nil
}
