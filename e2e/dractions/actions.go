package dractions

import (
	"context"
	"fmt"

	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
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
	r.Ctx.Log.Info("enter dractions EnableProtection")

	_, ok := d.(deployers.Subscription)
	if ok {
		client := r.Ctx.HubDynamicClient()

		resource := schema.GroupVersionResource{Group: "cluster.open-cluster-management.io", Version: "v1beta1", Resource: "placements"}
		resp, err := client.Resource(resource).Namespace("deployment-rbd").Get(context.TODO(), "placement", metav1.GetOptions{})
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return fmt.Errorf("could not get placement")
		}

		placement := clusterv1beta1.Placement{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(resp.UnstructuredContent(), &placement)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("could not FromUnstructured")
		}

		placement.Annotations[OCM_SCHEDULING_DISABLE] = "true"

		mapCR, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&placement)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("could not ToUnstructured")
		}
		unstructuredCR := &unstructured.Unstructured{Object: mapCR}
		_, err = client.Resource(resource).Namespace("deployment-rbd").Update(context.TODO(), unstructuredCR, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("could not update placment")
		}
		r.Ctx.Log.Info("updated placement annotation")
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
