package e2e_test

import (
	"fmt"
	"testing"
	"time"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const OCM_SCHEDULING_DISABLE = "cluster.open-cluster-management.io/experimental-scheduling-disable"
const DefaultDRPolicyName = "dr-policy"

func EnableProtection(w Workload, d Deployer) error {
	// If AppSet/Subscription, find Placement
	// Determine DRPolicy
	// Determine preferredCluster
	// Determine PVC label selector
	// Determine KubeObjectProtection requirements if Imperative (?)
	// Create DRPC, in desired namespace
	ctx.Log.Info("enter DRActions EnableProtection")

	_, isSub := d.(*Subscription)
	if isSub {
		// _, isAppSet := d.(*ApplicationSet)
		// if isSub || isAppSet {

		name := d.GetNamePrefix() + w.GetAppName()
		namespace := name
		drPolicyName := DefaultDRPolicyName
		appname := w.GetAppName()

		placementName := name
		drpcName := name
		client := ctx.Hub.CtrlClient

		placement, placementDecisionName, err := waitPlacementDecision(client, namespace, placementName)
		if err != nil {
			return err
		}

		ctx.Log.Info("get placementdecision " + placementDecisionName)
		placementDecision, err := getPlacementDecision(client, namespace, placementDecisionName)
		if err != nil {
			return err
		}

		clusterName := placementDecision.Status.Decisions[0].ClusterName
		ctx.Log.Info("placementdecision clusterName: " + clusterName)

		// move update placement annotation after placement has been handled
		// otherwise if we first add ocm disable annotation then it might not
		// yet be handled by ocm and thus PlacementSatisfied=false

		if placement.Annotations == nil {
			placement.Annotations = make(map[string]string)
		}

		placement.Annotations[OCM_SCHEDULING_DISABLE] = "true"

		ctx.Log.Info("update placement " + placementName + " annotation")
		err = updatePlacement(client, placement)
		if err != nil {
			return err
		}

		ctx.Log.Info("create drpc " + drpcName)
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

	} else {
		return fmt.Errorf("deployer not known")
	}
	return nil
}

func DisableProtection(w Workload, d Deployer) error {
	// remove DRPC
	// update placement annotation
	ctx.Log.Info("enter DRActions DisableProtection")

	_, isSub := d.(*Subscription)
	if isSub {
		// _, isAppSet := d.(*ApplicationSet)
		// if isSub || isAppSet {

		name := d.GetNamePrefix() + w.GetAppName()
		namespace := name

		placementName := name

		drpcName := name
		client := ctx.Hub.CtrlClient

		ctx.Log.Info("delete drpc " + drpcName)
		err := deleteDRPC(client, namespace, drpcName)
		if err != nil {
			return err
		}

		ctx.Log.Info("get placement " + placementName)
		placement, err := getPlacement(client, namespace, placementName)
		if err != nil {
			return err
		}

		delete(placement.Annotations, OCM_SCHEDULING_DISABLE)

		ctx.Log.Info("update placement " + placementName + " annotation")
		err = updatePlacement(client, placement)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("deployer not known")
	}
	return nil
}

func Failover(w Workload, d Deployer) error {
	// Determine DRPC
	// Check Placement
	// Failover to alternate in DRPolicy as the failoverCluster
	// Update DRPC
	ctx.Log.Info("enter DRActions Failover")

	name := d.GetNamePrefix() + w.GetAppName()
	namespace := name
	//placementName := w.GetPlacementName()
	drPolicyName := DefaultDRPolicyName
	drpcName := name
	client := ctx.Hub.CtrlClient

	// here we expect drpc should be ready before failover
	err := waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	ctx.Log.Info("get drpc " + drpcName)
	drpc, err := getDRPC(client, namespace, drpcName)
	if err != nil {
		return err
	}

	ctx.Log.Info("get drpolicy " + drPolicyName)
	drpolicy, err := getDRPolicy(client, drPolicyName)
	if err != nil {
		return err
	}

	targetCluster, err := getTargetCluster(client, namespace, name, drpolicy)
	if err != nil {
		return err
	}

	ctx.Log.Info("failover to cluster: " + targetCluster)
	drpc.Spec.Action = "Failover"
	drpc.Spec.FailoverCluster = targetCluster

	ctx.Log.Info("update drpc " + drpcName)
	err = updateDRPC(client, drpc)
	if err != nil {
		return err
	}

	//sleep to wait for DRPC is processed
	time.Sleep(5 * time.Second)

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

func Relocate(w Workload, d Deployer) error {
	// Determine DRPC
	// Check Placement
	// Relocate to Primary in DRPolicy as the PrimaryCluster
	// Update DRPC
	ctx.Log.Info("enter DRActions Relocate")

	name := d.GetNamePrefix() + w.GetAppName()
	namespace := name
	//placementName := w.GetPlacementName()
	drPolicyName := DefaultDRPolicyName
	drpcName := name
	client := ctx.Hub.CtrlClient

	// here we expect drpc should be ready before relocate
	err := waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	ctx.Log.Info("get drpc " + drpcName)
	drpc, err := getDRPC(client, namespace, drpcName)
	if err != nil {
		return err
	}

	ctx.Log.Info("get drpolicy " + drPolicyName)
	drpolicy, err := getDRPolicy(client, drPolicyName)
	if err != nil {
		return err
	}

	targetCluster, err := getTargetCluster(client, namespace, name, drpolicy)
	if err != nil {
		return err
	}
	ctx.Log.Info("relocate to cluster: " + targetCluster)
	drpc.Spec.Action = "Relocate"
	drpc.Spec.PreferredCluster = targetCluster

	ctx.Log.Info("update drpc " + drpcName)
	err = updateDRPC(client, drpc)
	if err != nil {
		return err
	}

	//sleep to wait for DRPC is processed
	time.Sleep(5 * time.Second)

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

func DeployAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	testCtx, err := getTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := testCtx.d.Deploy(testCtx.w); err != nil {
		t.Error(err)
	}
}

func EnableAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	testCtx, err := getTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := EnableProtection(testCtx.w, testCtx.d); err != nil {
		t.Error(err)
	}
}

func FailoverAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	testCtx, err := getTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := Failover(testCtx.w, testCtx.d); err != nil {
		t.Error(err)
	}
}

func RelocateAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	testCtx, err := getTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := Relocate(testCtx.w, testCtx.d); err != nil {
		t.Error(err)
	}
}

func DisableAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	testCtx, err := getTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := DisableProtection(testCtx.w, testCtx.d); err != nil {
		t.Error(err)
	}
}

func UndeployAction(t *testing.T) {
	ctx.Log.Info(t.Name())
	testCtx, err := getTestContext(t.Name())
	if err != nil {
		t.Error(err)
	}
	if err := testCtx.d.Undeploy(testCtx.w); err != nil {
		t.Error(err)
	}
}
