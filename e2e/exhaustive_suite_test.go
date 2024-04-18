// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"testing"
)

// var Deployers = []string{"Subscription", "AppSet", "Imperative"}
// var Workloads = []string{"Deployment", "STS", "DaemonSet"}
// var Classes = []string{"rbd", "cephfs"}

var currentDeployer Deployer
var currentWorkload Workload

func Exhaustive(t *testing.T) {
	t.Helper()

	ctx.Log.Info(t.Name())

	deployment := &Deployment{}
	deployment.Init()

	var Workloads = []Workload{deployment}

	subscrition := Subscription{}
	subscrition.Init()

	applicationSet := ApplicationSet{}
	applicationSet.Init()

	var Deployers = []Deployer{&subscrition, &applicationSet}

	for _, w := range Workloads {
		for _, d := range Deployers {

			currentWorkload = w
			currentDeployer = d

			t.Run(w.GetID(), func(t *testing.T) {
				ctx.Log.Info(t.Name())

				t.Run(d.GetID(), func(t *testing.T) {
					ctx.Log.Info(t.Name())

					addTestContext(t.Name(), w, d)

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
