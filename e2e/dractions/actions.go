package dractions

import (
	"context"
	"fmt"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
)

type DRActions struct {
	Ctx *util.TestContext
}

const OCM_SCHEDULING_DISABLE = "cluster.open-cluster-management.io/experimental-scheduling-disable"

func getPlacement(client *dynamic.DynamicClient, namespace, name string) (*clusterv1beta1.Placement, error) {

	resource := schema.GroupVersionResource{Group: "cluster.open-cluster-management.io", Version: "v1beta1", Resource: "placements"}
	unstr, err := client.Resource(resource).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, fmt.Errorf("could not get placement CR")
	}

	placement := clusterv1beta1.Placement{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstr.UnstructuredContent(), &placement)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("could not FromUnstructured in func getPlacment")
	}

	return &placement, nil
}

func updatePlacement(client *dynamic.DynamicClient, placement *clusterv1beta1.Placement) error {

	tempMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(placement)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not ToUnstructured")
	}

	unstr := &unstructured.Unstructured{Object: tempMap}
	resource := schema.GroupVersionResource{Group: "cluster.open-cluster-management.io", Version: "v1beta1", Resource: "placements"}
	_, err = client.Resource(resource).Namespace(placement.GetNamespace()).Update(context.TODO(), unstr, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not update placment")
	}

	return nil
}

func getPlacementDecision(client *dynamic.DynamicClient, namespace, name string) (*clusterv1beta1.PlacementDecision, error) {
	resource := schema.GroupVersionResource{Group: "cluster.open-cluster-management.io", Version: "v1beta1", Resource: "placementdecisions"}
	unstr, err := client.Resource(resource).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, fmt.Errorf("could not get placementDecision CR")
	}

	placementDecision := clusterv1beta1.PlacementDecision{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstr.UnstructuredContent(), &placementDecision)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("could not FromUnstructured in func getPlacementDecision")
	}

	return &placementDecision, nil
}

func (r DRActions) EnableProtection(w workloads.Workload, d deployers.Deployer) error {
	// If AppSet/Subscription, find Placement
	// Determine DRPolicy
	// Determine preferredCluster
	// Determine PVC label selector
	// Determine KubeObjectProtection requirements if Imperative (?)
	// Create DRPC, in desired namespace
	r.Ctx.Log.Info("enter dractions EnableProtection")

	_, ok := d.(deployers.Subscription)
	if ok {

		name := "deployment-rbd"
		namespace := "deployment-rbd"
		drPolicyName := "dr-policy"
		pvcLabel := "busybox"
		placementName := "placement"
		placementKind := "placement"
		client := r.Ctx.HubDynamicClient()

		placement, err := getPlacement(client, namespace, placementName)
		if err != nil {
			return err
		}

		placement.Annotations[OCM_SCHEDULING_DISABLE] = "true"

		err = updatePlacement(client, placement)
		if err != nil {
			return err
		}

		// L1
		placement, err = getPlacement(client, namespace, placementName)
		if err != nil {
			return err
		}
		placementDecisionName := ""
		for _, cond := range placement.Status.Conditions {
			if cond.Type == "PlacementSatisfied" && cond.Status == "True" {
				placementDecisionName = placement.Status.DecisionGroups[0].Decisions[0]
				// kubectl.get(
				// 	"placementdecisions",
				// 	f"--selector=cluster.open-cluster-management.io/placement={placement_name}",
				// 	f"--namespace={config['namespace']}",
				// 	"--output=jsonpath={.items[0].status.decisions[0].clusterName}",
				// 	context=env["hub"],
			}
		}
		if placementDecisionName == "" {
			fmt.Println("can not find placement decision")
			// if not timeout, go to L1
			// else return error
		}

		placementDecision, err := getPlacementDecision(client, namespace, placementDecisionName)
		if err != nil {
			return err
		}

		clusterName := placementDecision.Status.Decisions[0].ClusterName
		fmt.Printf("clusterName: %v\n", clusterName)

		drpc := &ramen.DRPlacementControl{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name + "-drpc",
				Namespace: namespace,
				Labels:    map[string]string{"app": name},
			},
			Spec: ramen.DRPlacementControlSpec{
				PreferredCluster: clusterName,
				DRPolicyRef: v1.ObjectReference{
					Name: drPolicyName,
				},
				PlacementRef: v1.ObjectReference{
					Kind: placementKind,
					Name: placementName,
				},
				PVCSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"appname": pvcLabel},
				},
			},
		}

		tempMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(drpc)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("could not ToUnstructured")
		}

		unstr := &unstructured.Unstructured{Object: tempMap}
		resource := schema.GroupVersionResource{Group: "cluster.open-cluster-management.io", Version: "v1beta1", Resource: "drplacementcontrols"}
		_, err = client.Resource(resource).Namespace(drpc.GetNamespace()).Update(context.TODO(), unstr, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("could not update drplacementcontrol")
		}

		return nil

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
