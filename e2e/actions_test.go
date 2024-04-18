package e2e

import (
	"testing"

	"github.com/ramendr/ramen/e2e/dractions"
	"github.com/ramendr/ramen/e2e/testcontext"
	"github.com/ramendr/ramen/e2e/util"
)

func DeployAction(t *testing.T) {
	util.Ctx.Log.Info(t.Name())
	testCtx, err := testcontext.GetTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := testCtx.Deployer.Deploy(testCtx.Workload); err != nil {
		t.Error(err)
	}
}

func EnableAction(t *testing.T) {
	util.Ctx.Log.Info(t.Name())
	testCtx, err := testcontext.GetTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := dractions.EnableProtection(testCtx.Workload, testCtx.Deployer); err != nil {
		t.Error(err)
	}
}

func FailoverAction(t *testing.T) {
	util.Ctx.Log.Info(t.Name())
	testCtx, err := testcontext.GetTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := dractions.Failover(testCtx.Workload, testCtx.Deployer); err != nil {
		t.Error(err)
	}
}

func RelocateAction(t *testing.T) {
	util.Ctx.Log.Info(t.Name())
	testCtx, err := testcontext.GetTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := dractions.Relocate(testCtx.Workload, testCtx.Deployer); err != nil {
		t.Error(err)
	}
}

func DisableAction(t *testing.T) {
	util.Ctx.Log.Info(t.Name())
	testCtx, err := testcontext.GetTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := dractions.DisableProtection(testCtx.Workload, testCtx.Deployer); err != nil {
		t.Error(err)
	}
}

func UndeployAction(t *testing.T) {
	util.Ctx.Log.Info(t.Name())
	testCtx, err := testcontext.GetTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := testCtx.Deployer.Undeploy(testCtx.Workload); err != nil {
		t.Error(err)
	}
}
