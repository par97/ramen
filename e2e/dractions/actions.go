package dractions

import (
	"fmt"
	"time"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	OcmSchedulingDisable = "cluster.open-cluster-management.io/experimental-scheduling-disable"
	DefaultDRPolicyName  = "dr-policy"
	FiveSecondsDuration  = 5 * time.Second
)

func EnableProtection(w workloads.Workload, d deployers.Deployer) error {
	// If AppSet/Subscription, find Placement
	// Determine DRPolicy
	// Determine preferredCluster
	// Determine PVC label selector
	// Determine KubeObjectProtection requirements if Imperative (?)
	// Create DRPC, in desired namespace
	util.Ctx.Log.Info("enter DRActions EnableProtection")

	_, isSub := d.(*deployers.Subscription)
	_, isAppSet := d.(*deployers.ApplicationSet)

	if !isSub && !isAppSet {
		return fmt.Errorf("deployers.Deployer not known")
	}

	name := d.GetNamePrefix() + w.GetAppName()
	namespace := name

	if isAppSet {
		namespace = util.ArgocdNamespace
	}

	drPolicyName := DefaultDRPolicyName
	appname := w.GetAppName()

	placementName := name
	drpcName := name
	client := util.Ctx.Hub.CtrlClient

	placement, placementDecisionName, err := waitPlacementDecision(client, namespace, placementName)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("get placementdecision " + placementDecisionName)

	placementDecision, err := getPlacementDecision(client, namespace, placementDecisionName)
	if err != nil {
		return err
	}

	clusterName := placementDecision.Status.Decisions[0].ClusterName
	util.Ctx.Log.Info("placementdecision clusterName: " + clusterName)

	// move update placement annotation after placement has been handled
	// otherwise if we first add ocm disable annotation then it might not
	// yet be handled by ocm and thus PlacementSatisfied=false

	if placement.Annotations == nil {
		placement.Annotations = make(map[string]string)
	}

	placement.Annotations[OcmSchedulingDisable] = "true"

	util.Ctx.Log.Info("update placement " + placementName + " annotation")

	err = updatePlacement(client, placement)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("create drpc " + drpcName)
	drpc := &ramen.DRPlacementControl{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DRPlacementControl",
			APIVersion: "ramendr.openshift.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      drpcName,
			Namespace: namespace,
			Labels:    map[string]string{"app": name},
		},
		Spec: ramen.DRPlacementControlSpec{
			PreferredCluster: clusterName,
			DRPolicyRef: v1.ObjectReference{
				Name: drPolicyName,
			},
			PlacementRef: v1.ObjectReference{
				Kind: "placement",
				Name: placementName,
			},
			PVCSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{"appname": appname},
			},
		},
	}

	err = createDRPC(client, drpc)
	if err != nil {
		return err
	}

	err = waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	return nil
}

func DisableProtection(w workloads.Workload, d deployers.Deployer) error {
	// remove DRPC
	// update placement annotation
	util.Ctx.Log.Info("enter DRActions DisableProtection")

	_, isSub := d.(*deployers.Subscription)
	_, isAppSet := d.(*deployers.ApplicationSet)

	if !isSub && !isAppSet {
		return fmt.Errorf("deployers.Deployer not known")
	}

	name := d.GetNamePrefix() + w.GetAppName()

	namespace := name
	if isAppSet {
		namespace = util.ArgocdNamespace
	}

	placementName := name

	drpcName := name
	client := util.Ctx.Hub.CtrlClient

	util.Ctx.Log.Info("delete drpc " + drpcName)

	err := deleteDRPC(client, namespace, drpcName)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("get placement " + placementName)

	placement, err := getPlacement(client, namespace, placementName)
	if err != nil {
		return err
	}

	delete(placement.Annotations, OcmSchedulingDisable)

	util.Ctx.Log.Info("update placement " + placementName + " annotation")

	err = updatePlacement(client, placement)
	if err != nil {
		return err
	}

	return nil
}

func Failover(w workloads.Workload, d deployers.Deployer) error {
	util.Ctx.Log.Info("enter DRActions Failover")

	name := d.GetNamePrefix() + w.GetAppName()
	namespace := name

	_, isAppSet := d.(*deployers.ApplicationSet)
	if isAppSet {
		namespace = util.ArgocdNamespace
	}

	drPolicyName := DefaultDRPolicyName
	drpcName := name
	client := util.Ctx.Hub.CtrlClient

	// here we expect drpc should be ready before failover
	err := waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("get drpc " + drpcName)

	drpc, err := getDRPC(client, namespace, drpcName)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("get drpolicy " + drPolicyName)

	drpolicy, err := getDRPolicy(client, drPolicyName)
	if err != nil {
		return err
	}

	targetCluster, err := getTargetCluster(client, namespace, name, drpolicy)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("failover to cluster: " + targetCluster)

	drpc.Spec.Action = "Failover"
	drpc.Spec.FailoverCluster = targetCluster

	util.Ctx.Log.Info("update drpc " + drpcName)

	err = updateDRPC(client, drpc)
	if err != nil {
		return err
	}

	// sleep to wait for DRPC is processed
	time.Sleep(FiveSecondsDuration)

	// check Phase
	err = waitDRPCPhase(client, namespace, drpcName, "FailedOver")
	if err != nil {
		return err
	}
	// then check Conditions
	err = waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	return nil
}

func Relocate(w workloads.Workload, d deployers.Deployer) error {
	// Determine DRPC
	// Check Placement
	// Relocate to Primary in DRPolicy as the PrimaryCluster
	// Update DRPC
	util.Ctx.Log.Info("enter DRActions Relocate")

	name := d.GetNamePrefix() + w.GetAppName()
	namespace := name

	_, isAppSet := d.(*deployers.ApplicationSet)
	if isAppSet {
		namespace = util.ArgocdNamespace
	}
	// placementName := w.GetPlacementName()
	drPolicyName := DefaultDRPolicyName
	drpcName := name
	client := util.Ctx.Hub.CtrlClient

	// here we expect drpc should be ready before relocate
	err := waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("get drpc " + drpcName)

	drpc, err := getDRPC(client, namespace, drpcName)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("get drpolicy " + drPolicyName)

	drpolicy, err := getDRPolicy(client, drPolicyName)
	if err != nil {
		return err
	}

	targetCluster, err := getTargetCluster(client, namespace, name, drpolicy)
	if err != nil {
		return err
	}

	util.Ctx.Log.Info("relocate to cluster: " + targetCluster)

	drpc.Spec.Action = "Relocate"
	drpc.Spec.PreferredCluster = targetCluster

	util.Ctx.Log.Info("update drpc " + drpcName)

	err = updateDRPC(client, drpc)
	if err != nil {
		return err
	}

	// sleep to wait for DRPC is processed
	time.Sleep(FiveSecondsDuration)

	// check Phase
	err = waitDRPCPhase(client, namespace, drpcName, "Relocated")
	if err != nil {
		return err
	}
	// then check Conditions
	err = waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	return nil
}
