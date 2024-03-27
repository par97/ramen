package dractions

import (
	"context"
	"fmt"
	"time"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"open-cluster-management.io/api/cluster/v1beta1"
)

type DRActions struct {
	Ctx *util.TestContext
}

const OCM_SCHEDULING_DISABLE = "cluster.open-cluster-management.io/experimental-scheduling-disable"

func (r DRActions) EnableProtection(w workloads.Workload, d deployers.Deployer) error {
	// If AppSet/Subscription, find Placement
	// Determine DRPolicy
	// Determine preferredCluster
	// Determine PVC label selector
	// Determine KubeObjectProtection requirements if Imperative (?)
	// Create DRPC, in desired namespace
	r.Ctx.Log.Info("enter DRActions EnableProtection")

	_, ok := d.(deployers.Subscription)
	if ok {

		name := w.GetName()
		namespace := w.GetNameSpace()
		drPolicyName := util.DefaultDRPolicy
		pvcLabel := w.GetPVCLabel()
		placementName := w.GetPlacementName()
		drpcName := name + "-drpc"
		client := r.Ctx.HubDynamicClient()

		r.Ctx.Log.Info("get placement " + placementName + " and wait for PlacementSatisfied=True")

		var placement *v1beta1.Placement
		var err error
		placementDecisionName := ""

		timeout := 300 //seconds
		interval := 30 //seconds
		startTime := time.Now()

		for {
			placement, err = getPlacement(client, namespace, placementName)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return err
			}

			for _, cond := range placement.Status.Conditions {
				if cond.Type == "PlacementSatisfied" && cond.Status == "True" {
					placementDecisionName = placement.Status.DecisionGroups[0].Decisions[0]
				}
			}
			if placementDecisionName != "" {
				r.Ctx.Log.Info("got placementdecision name" + placementDecisionName)
				break
			}
			if time.Since(startTime) > time.Second*time.Duration(timeout) {
				fmt.Println("time out")
				return fmt.Errorf("could not find placement decision before timeout")
			}
			r.Ctx.Log.Info(fmt.Sprintf("placementSatisfied is not True, sleep %v seconds", interval))
			time.Sleep(time.Second * time.Duration(interval))
		}

		r.Ctx.Log.Info("get placementdecision " + placementDecisionName)
		placementDecision, err := getPlacementDecision(client, namespace, placementDecisionName)
		if err != nil {
			return err
		}

		clusterName := placementDecision.Status.Decisions[0].ClusterName
		r.Ctx.Log.Info("placementdecision clusterName: " + clusterName)

		// move update placement annotation after placement has been handled
		// otherwise if we first add ocm disable annotation then it might not
		// yet be handled by ocm and thus PlacementSatisfied=false

		placement.Annotations[OCM_SCHEDULING_DISABLE] = "true"

		r.Ctx.Log.Info("update placement " + placementName + " annotation")
		err = updatePlacement(client, placement)
		if err != nil {
			return err
		}

		r.Ctx.Log.Info("create drpc " + drpcName)
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
					MatchLabels: map[string]string{"appname": pvcLabel},
				},
			},
		}

		tempMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(drpc)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return fmt.Errorf("could not ToUnstructured")
		}

		unstr := &unstructured.Unstructured{Object: tempMap}
		resource := schema.GroupVersionResource{Group: "ramendr.openshift.io", Version: "v1alpha1", Resource: "drplacementcontrols"}

		_, err = client.Resource(resource).Namespace(namespace).Create(context.Background(), unstr, metav1.CreateOptions{})
		if err != nil {
			if !k8serrors.IsAlreadyExists(err) {
				fmt.Printf("err: %v\n", err)
				return fmt.Errorf("could not create drpc " + drpcName)
			}
			r.Ctx.Log.Info("drpc " + drpcName + " already Exists")

		}

		timeout = 300 //seconds
		interval = 30 //seconds
		startTime = time.Now()
		for {
			ready := true
			drpc, err = getDRPlacementControl(client, namespace, drpcName)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return err
			}

			for _, cond := range drpc.Status.Conditions {
				if cond.Type == "Available" && cond.Status != "True" {
					r.Ctx.Log.Info("drpc status Available is not True")
					ready = false
					break
				}
				if cond.Type == "PeerReady" && cond.Status != "True" {
					r.Ctx.Log.Info("drpc status PeerReady is not True")
					ready = false
					break
				}
			}
			if ready {
				if drpc.Status.LastGroupSyncTime == nil {
					r.Ctx.Log.Info("drpc status LastGroupSyncTime is nil")
					ready = false
				}
			}
			if ready {
				r.Ctx.Log.Info("drpc status is ready")
				break
			}
			if time.Since(startTime) > time.Second*time.Duration(timeout) {
				fmt.Println("time out")
				return fmt.Errorf("drpc status is not ready yet before timeout")
			}
			r.Ctx.Log.Info(fmt.Sprintf("drpc status is not ready yet, sleep %v seconds", interval))
			time.Sleep(time.Second * time.Duration(interval))
		}
	} else {
		return fmt.Errorf("deployer not known")
	}
	return nil
}

func (r DRActions) DisableProtection(w workloads.Workload, d deployers.Deployer) error {
	// remove DRPC
	// update placement annotation
	r.Ctx.Log.Info("enter DRActions DisableProtection")

	_, ok := d.(deployers.Subscription)
	if ok {

		name := w.GetName()
		namespace := w.GetNameSpace()
		placementName := w.GetPlacementName()
		drpcName := name + "-drpc"
		client := r.Ctx.HubDynamicClient()

		r.Ctx.Log.Info("delete drpc " + drpcName)
		err := deleteDRPlacementControl(client, namespace, drpcName)
		if err != nil {
			return err
		}

		r.Ctx.Log.Info("get placement " + placementName)
		placement, err := getPlacement(client, namespace, placementName)
		if err != nil {
			return err
		}

		delete(placement.Annotations, OCM_SCHEDULING_DISABLE)

		r.Ctx.Log.Info("update placement " + placementName + " annotation")
		err = updatePlacement(client, placement)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("deployer not known")
	}
	return nil
}

func (r DRActions) Failover(w workloads.Workload, d deployers.Deployer) error {
	// Determine DRPC
	// Check Placement
	// Failover to alternate in DRPolicy as the failoverCluster
	// Update DRPC
	r.Ctx.Log.Info("enter dractions Failover")
	return nil
}

func (r DRActions) Relocate(w workloads.Workload, d deployers.Deployer) error {
	// Determine DRPC
	// Check Placement
	// Relocate to Primary in DRPolicy as the PrimaryCluster
	// Update DRPC
	r.Ctx.Log.Info("enter dractions Relocate")
	return nil
}
