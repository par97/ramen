package dractions

import (
	"context"
	"fmt"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
)

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
		fmt.Printf("err: %v\n", err)
		return nil, fmt.Errorf("could not FromUnstructured in func getPlacment")
	}

	return &placement, nil
}

func updatePlacement(client *dynamic.DynamicClient, placement *clusterv1beta1.Placement) error {

	tempMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(placement)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fmt.Errorf("could not ToUnstructured")
	}

	unstr := &unstructured.Unstructured{Object: tempMap}
	resource := schema.GroupVersionResource{Group: "cluster.open-cluster-management.io", Version: "v1beta1", Resource: "placements"}
	_, err = client.Resource(resource).Namespace(placement.GetNamespace()).Update(context.Background(), unstr, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
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
		fmt.Printf("err: %v\n", err)
		return nil, fmt.Errorf("could not FromUnstructured in func getPlacementDecision")
	}

	return &placementDecision, nil
}

func getDRPC(client *dynamic.DynamicClient, namespace, name string) (*ramen.DRPlacementControl, error) {
	resource := schema.GroupVersionResource{Group: "ramendr.openshift.io", Version: "v1alpha1", Resource: "drplacementcontrols"}
	unstr, err := client.Resource(resource).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, fmt.Errorf("could not get drpc CR")
	}

	drpc := ramen.DRPlacementControl{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstr.UnstructuredContent(), &drpc)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, fmt.Errorf("could not FromUnstructured in func getDRPlacementControl")
	}

	return &drpc, nil
}

func (r DRActions) createDRPC(client *dynamic.DynamicClient, drpc *ramen.DRPlacementControl) error {

	tempMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(drpc)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fmt.Errorf("could not ToUnstructured")
	}

	unstr := &unstructured.Unstructured{Object: tempMap}
	resource := schema.GroupVersionResource{Group: "ramendr.openshift.io", Version: "v1alpha1", Resource: "drplacementcontrols"}
	_, err = client.Resource(resource).Namespace(drpc.GetNamespace()).Create(context.Background(), unstr, metav1.CreateOptions{})

	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return fmt.Errorf("could not create drpc " + drpc.Name)
		}
		r.Ctx.Log.Info("drpc " + drpc.Name + " already Exists")
	}

	return nil
}

func updateDRPC(client *dynamic.DynamicClient, drpc *ramen.DRPlacementControl) error {

	tempMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(drpc)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fmt.Errorf("could not ToUnstructured")
	}

	unstr := &unstructured.Unstructured{Object: tempMap}
	resource := schema.GroupVersionResource{Group: "ramendr.openshift.io", Version: "v1alpha1", Resource: "drplacementcontrols"}
	_, err = client.Resource(resource).Namespace(drpc.GetNamespace()).Update(context.Background(), unstr, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fmt.Errorf("could not update drpc CR")
	}

	return nil
}

func getDRPolicy(client *dynamic.DynamicClient, name string) (*ramen.DRPolicy, error) {
	resource := schema.GroupVersionResource{Group: "ramendr.openshift.io", Version: "v1alpha1", Resource: "drpolicies"}
	unstr, err := client.Resource(resource).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, fmt.Errorf("could not get drpolicies")
	}

	drpolicy := ramen.DRPolicy{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstr.UnstructuredContent(), &drpolicy)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, fmt.Errorf("could not FromUnstructured in func getDRPolicy")
	}

	return &drpolicy, nil
}

func deleteDRPC(client *dynamic.DynamicClient, namespace, name string) error {

	resource := schema.GroupVersionResource{Group: "ramendr.openshift.io", Version: "v1alpha1", Resource: "drplacementcontrols"}
	err := client.Resource(resource).Namespace(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}
