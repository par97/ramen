// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	// "context"
	// "fmt"
	// "strings"
	"testing"

	"github.com/ramendr/ramen/e2e/util"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes"
)

func Validate(t *testing.T) {
	t.Helper()

	util.Ctx.Log.Info(t.Name())

	t.Run("CheckRamenHubOperatorStatus", CheckRamenHubOperatorStatus)
	t.Run("RamenSpokes", RamenSpoke)
	t.Run("Ceph", Ceph)
}

func CheckRamenHubOperatorStatus(t *testing.T) {
	util.Ctx.Log.Info(t.Name())

	isRunning, podName, err := util.CheckRamenHubPodRunningStatus(util.Ctx.Hub.K8sClientSet)
	if err != nil {
		t.Error(err)
	}

	if isRunning {
		util.Ctx.Log.Info("Ramen Hub Operator is running", "pod", podName)
	} else {
		t.Error("no running Ramen Hub Operator pod")
	}

	util.Ctx.Log.Info("TestRamenHubOperatorStatus: Pass")
}

func RamenSpoke(t *testing.T) {
	util.Ctx.Log.Info(t.Name())
}

func Ceph(t *testing.T) {
	util.Ctx.Log.Info(t.Name())
}
