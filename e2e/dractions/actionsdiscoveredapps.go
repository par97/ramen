// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package dractions

import (
	"strings"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

const RamenOpsNs = "ramen-ops"

func EnableProtectionDiscoveredApps(w workloads.Workload, d deployers.Deployer) error {
	util.Ctx.Log.Info("enter EnableProtectionDiscoveredApps " + w.GetName() + "/" + d.GetName())

	name := GetCombinedName(d, w)
	namespace := name // namespace of the app will be deployed into drclusters

	drPolicyName := DefaultDRPolicyName
	appname := w.GetAppName()
	placementName := name
	drpcName := name

	// create mcsb default in ramen-ops ns
	if err := deployers.CreateManagedClusterSetBinding(deployers.McsbName, RamenOpsNs); err != nil {
		return err
	}

	// create placement
	if err := createPlacementManagedByRamen(placementName, RamenOpsNs); err != nil {
		return err
	}

	// create drpc
	drpolicy, err := getDRPolicy(util.Ctx.Hub.CtrlClient, drPolicyName)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("create drpc " + drpcName)

	clusterName := drpolicy.Spec.DRClusters[0]

	drpc := generateDRPCDiscoveredApps(name, RamenOpsNs, clusterName, drPolicyName, placementName, appname, namespace)
	if err = createDRPC(util.Ctx.Hub.CtrlClient, drpc); err != nil {
		return err
	}

	// wait for drpc ready
	return waitDRPCReady(util.Ctx.Hub.CtrlClient, RamenOpsNs, drpcName)
}

// remove DRPC
// update placement annotation
func DisableProtectionDiscoveredApps(w workloads.Workload, d deployers.Deployer) error {
	util.Ctx.Log.Info("enter DisableProtectionDiscoveredApps DisableProtection")

	name := GetCombinedName(d, w)

	placementName := name
	drpcName := name

	client := util.Ctx.Hub.CtrlClient

	util.Ctx.Log.Info("delete drpc " + drpcName)

	if err := deleteDRPC(client, RamenOpsNs, drpcName); err != nil {
		return err
	}

	if err := waitDRPCDeleted(client, RamenOpsNs, drpcName); err != nil {
		return err
	}

	// delete placement
	if err := deployers.DeletePlacement(placementName, RamenOpsNs); err != nil {
		return err
	}

	return deployers.DeleteManagedClusterSetBinding(deployers.McsbName, RamenOpsNs)
}

func FailoverDiscoveredApps(w workloads.Workload, d deployers.Deployer) error {
	util.Ctx.Log.Info("enter DRActions FailoverDiscoveredApps")

	return failoverRelocateAction(w, d, "Failover")
}

func RelocateDiscoveredApps(w workloads.Workload, d deployers.Deployer) error {
	util.Ctx.Log.Info("enter DRActions RelocateDiscoveredApps")

	return failoverRelocateAction(w, d, "Relocate")
}

// nolint:funlen
func failoverRelocateAction(w workloads.Workload, d deployers.Deployer, action string) error {
	name := GetCombinedName(d, w)
	namespace := name

	drPolicyName := DefaultDRPolicyName
	drpcName := name
	hubClient := util.Ctx.Hub.CtrlClient

	// here we expect drpc should be ready before action
	if err := waitDRPCReady(hubClient, RamenOpsNs, drpcName); err != nil {
		return err
	}

	util.Ctx.Log.Info("get drpc " + drpcName)

	drpc, err := getDRPC(hubClient, RamenOpsNs, drpcName)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("get drpolicy " + drPolicyName)

	drpolicy, err := getDRPolicy(hubClient, drPolicyName)
	if err != nil {
		return err
	}

	currentCluster, err := getCurrentCluster(hubClient, RamenOpsNs, name)
	if err != nil {
		return err
	}

	drClient := getDRClusterClient(currentCluster, drpolicy)

	targetCluster, err := getTargetCluster(hubClient, RamenOpsNs, name, drpolicy)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info(strings.ToLower(action) + " to cluster: " + targetCluster)

	drpc.Spec.Action = ramen.DRAction(action)
	drpc.Spec.PreferredCluster = targetCluster

	util.Ctx.Log.Info("update drpc " + drpcName)

	if err = updateDRPC(hubClient, drpc); err != nil {
		return err
	}

	if err = waitDRPCProgression(hubClient, RamenOpsNs, name, "WaitOnUserToCleanUp"); err != nil {
		return err
	}

	// delete pvc and deployment from dr cluster
	util.Ctx.Log.Info("start to clean up discovered apps from " + currentCluster)

	if err = deployers.DeleteDiscoveredApps(drClient, namespace); err != nil {
		return err
	}

	if err = waitDRPCProgression(hubClient, RamenOpsNs, name, "Completed"); err != nil {
		return err
	}

	return waitDRPCReady(hubClient, RamenOpsNs, name)
}
