// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"testing"
)

// var Deployers = []string{"Subscription", "AppSet", "Imperative"}
// var Workloads = []string{"Deployment", "STS", "DaemonSet"}
// var Classes = []string{"rbd", "cephfs"}

func Exhaustive(t *testing.T) {
	t.Helper()

	ctx.Log.Info(t.Name())

	deployment := &Deployment{}
	deployment.Init()
	var Workloads = []Workload{deployment}

	subscrition := Subscription{}
	subscrition.Init()
	var Deployers = []Deployer{&subscrition}

	for _, w := range Workloads {
		for _, d := range Deployers {

			t.Run(w.GetID(), func(t *testing.T) {
				ctx.Log.Info(t.Name())

				t.Run(d.GetID(), func(t *testing.T) {
					ctx.Log.Info(t.Name())

					if err := d.Deploy(w); err != nil {
						t.Error(err)
					}
					if err := EnableProtection(w, d); err != nil {
						t.Error(err)
					}
					if err := Failover(w, d); err != nil {
						t.Error(err)
					}
					if err := Relocate(w, d); err != nil {
						t.Error(err)
					}
					if err := DisableProtection(w, d); err != nil {
						t.Error(err)
					}
					if err := d.Undeploy(w); err != nil {
						t.Error(err)
					}

					// t.Run("Deploy", Deploy)
					// t.Run("Enable", Enable)
					// t.Run("Failover", Failover)
					// t.Run("Relocate", Relocate)
					// t.Run("Disable", Disable)
					// t.Run("Undeploy", Undeploy)
				})
			})
		}
	}
}
