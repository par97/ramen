// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"testing"

	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/testcontext"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

// var Deployers = []string{"Subscription", "AppSet", "Imperative"}
// var Workloads = []string{"Deployment", "STS", "DaemonSet"}
// var Classes = []string{"rbd", "cephfs"}

var currentDeployer deployers.Deployer
var currentWorkload workloads.Workload

func Exhaustive(t *testing.T) {
	t.Parallel()
	t.Helper()

	util.Ctx.Log.Info(t.Name())

	deployment := &workloads.Deployment{}
	deployment.Init()

	var Workloads = []workloads.Workload{deployment}

	subscrition := deployers.Subscription{}
	subscrition.Init()

	applicationSet := deployers.ApplicationSet{}
	applicationSet.Init()

	var Deployers = []deployers.Deployer{&subscrition, &applicationSet}

	for _, w := range Workloads {
		// this is needed to avoid parallel test issue
		// see https://go.dev/wiki/CommonMistakes
		w := w
		for _, d := range Deployers {
			// this is needed to avoid parallel test issue
			// see https://go.dev/wiki/CommonMistakes
			d := d

			currentWorkload = w
			currentDeployer = d

			t.Run(w.GetID(), func(t *testing.T) {
				t.Parallel()
				util.Ctx.Log.Info(t.Name())

				t.Run(d.GetID(), func(t *testing.T) {
					t.Parallel()

					util.Ctx.Log.Info(t.Name())

					testcontext.AddTestContext(t.Name(), w, d)

					if !t.Run("Deploy", DeployAction) {
						t.Fatal("Deploy failed")
					}
					if !t.Run("Enable", EnableAction) {
						t.Fatal("Enable failed")
					}
					if !t.Run("Failover", FailoverAction) {
						t.Fatal("Failover failed")
					}
					if !t.Run("Relocate", RelocateAction) {
						t.Fatal("Relocate failed")
					}
					if !t.Run("Disable", DisableAction) {
						t.Fatal("Disable failed")
					}
					if !t.Run("Undeploy", UndeployAction) {
						t.Fatal("Undeploy failed")
					}
				})
			})
		}
	}
}
