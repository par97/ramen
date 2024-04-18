// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"testing"
)

func Basic(t *testing.T) {
	t.Helper()

	ctx.Log.Info(t.Name())

	// t.Run("Deploy", DeployAction)
	// t.Run("Enable", EnableAction)
	// t.Run("Failover", FailoverAction)
	// t.Run("Relocate", RelocateAction)
	// t.Run("Disable", DisableAction)
	// t.Run("Undeploy", UndeployAction)
}

func DeployAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	if err := currentDeployer.Deploy(currentWorkload); err != nil {
		t.Error(err)
	}
}

func EnableAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	if err := EnableProtection(currentWorkload, currentDeployer); err != nil {
		t.Error(err)
	}
}

func FailoverAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	if err := Failover(currentWorkload, currentDeployer); err != nil {
		t.Error(err)
	}
}

func RelocateAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	if err := Relocate(currentWorkload, currentDeployer); err != nil {
		t.Error(err)
	}
}

func DisableAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	if err := DisableProtection(currentWorkload, currentDeployer); err != nil {
		t.Error(err)
	}
}

func UndeployAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	if err := currentDeployer.Undeploy(currentWorkload); err != nil {
		t.Error(err)
	}
}
