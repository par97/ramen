package e2e_test

import (
	"testing"

	"github.com/ramendr/ramen/e2e/util"
)

func Basic(t *testing.T) {
	t.Helper()

	util.Ctx.Log.Info(t.Name())
	// t.Run("Deploy", DeployAction)
	// t.Run("Enable", EnableAction)
	// t.Run("Failover", FailoverAction)
	// t.Run("Relocate", RelocateAction)
	// t.Run("Disable", DisableAction)
	// t.Run("Undeploy", UndeployAction)
}
