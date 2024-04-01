package workloads

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	ocmclusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	ocmclusterv1beta2 "open-cluster-management.io/api/cluster/v1beta2"
	channelv1 "open-cluster-management.io/multicloud-operators-channel/pkg/apis/apps/v1"
	placementrulev1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"
	subscriptionv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/v1"
)

func (w *Deployment) createNamespace(namespace string) error {
	w.Ctx.Log.Info("enter Deployment createNamespace " + namespace)

	objNs := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NameSpace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	_, err := w.Ctx.HubK8sClientSet().CoreV1().Namespaces().Create(context.Background(), objNs, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

/*
	func (w *Deployment) createChannel() error {
		w.Ctx.Log.Info("enter Deployment createChannel")

		objChannel := &channelv1.Channel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      w.ChannelName,
				Namespace: w.ChannelNamespace,
			},
			Spec: channelv1.ChannelSpec{
				Pathname: w.RepoURL,
				Type:     w.ChannelType,
			},
		}

		err := w.Ctx.HubCtrlClient().Create(context.Background(), objChannel)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}

		return nil
	}
*/
func (w *Deployment) createSubscription() error {
	w.Ctx.Log.Info("enter Deployment createSubscription")

	labels := make(map[string]string)
	labels["apps"] = "deployment-rbd"

	annotations := make(map[string]string)
	annotations["apps.open-cluster-management.io/github-branch"] = w.Revision
	annotations["apps.open-cluster-management.io/github-path"] = w.Path

	// Construct PlacementRef
	objReplacementRef := corev1.ObjectReference{
		Kind: "Placement",
		Name: w.PlacementName,
	}

	objPlacementRulePlacement := &placementrulev1.Placement{}
	objPlacementRulePlacement.PlacementRef = &objReplacementRef

	objSubscription := &subscriptionv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "subscription",
			Namespace:   w.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: subscriptionv1.SubscriptionSpec{
			Channel:   w.ChannelNamespace + "/" + w.ChannelName,
			Placement: objPlacementRulePlacement,
		},
	}

	err := w.Ctx.HubCtrlClient().Create(context.Background(), objSubscription)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (w *Deployment) createPlacement() error {
	w.Ctx.Log.Info("enter Deployment createPlacement")

	labels := make(map[string]string)
	labels["apps"] = "deployment-rbd"

	arrayClusterSets := []string{"default"}
	var numClusters int32 = 1

	objPlacement := &ocmclusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      w.PlacementName,
			Namespace: w.Namespace,
			Labels:    labels,
		},
		Spec: ocmclusterv1beta1.PlacementSpec{
			ClusterSets:      arrayClusterSets,
			NumberOfClusters: &numClusters,
		},
	}

	err := w.Ctx.HubCtrlClient().Create(context.Background(), objPlacement)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (w *Deployment) createManagedClusterSetBinding() error {
	w.Ctx.Log.Info("enter Deployment createManagedClusterSetBinding")

	labels := make(map[string]string)
	labels["apps"] = "deployment-rbd"

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: w.Namespace,
			Labels:    labels,
		},
		Spec: ocmclusterv1beta2.ManagedClusterSetBindingSpec{
			ClusterSet: "default",
		},
	}

	err := w.Ctx.HubCtrlClient().Create(context.Background(), objMCSB)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (w *Deployment) deleteNamespace(namespace string) error {
	w.Ctx.Log.Info("enter Deployment deleteNamespace " + namespace)

	err := w.Ctx.HubK8sClientSet().CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (w *Deployment) deleteChannel() error {
	w.Ctx.Log.Info("enter Deployment deleteChannel")

	objChannel := &channelv1.Channel{}

	key := types.NamespacedName{
		Name:      w.ChannelName,
		Namespace: w.ChannelNamespace,
	}

	err := w.Ctx.HubCtrlClient().Get(context.Background(), key, objChannel)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	err = w.Ctx.HubCtrlClient().Delete(context.Background(), objChannel)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (w *Deployment) deleteSubscription() error {
	w.Ctx.Log.Info("enter Deployment deleteSubscription")

	objSubscription := &subscriptionv1.Subscription{}

	key := types.NamespacedName{
		Name:      "subscription",
		Namespace: w.Namespace,
	}

	err := w.Ctx.HubCtrlClient().Get(context.Background(), key, objSubscription)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	err = w.Ctx.HubCtrlClient().Delete(context.Background(), objSubscription)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (w *Deployment) deletePlacement() error {
	w.Ctx.Log.Info("enter Deployment deletePlacement")

	objPlacement := &ocmclusterv1beta1.Placement{}

	key := types.NamespacedName{
		Name:      w.PlacementName,
		Namespace: w.Namespace,
	}

	err := w.Ctx.HubCtrlClient().Get(context.Background(), key, objPlacement)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	err = w.Ctx.HubCtrlClient().Delete(context.Background(), objPlacement)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (w *Deployment) deleteManagedClusterSetBinding() error {
	w.Ctx.Log.Info("enter Deployment deleteManagedClusterSetBinding")

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{}

	key := types.NamespacedName{
		Name:      "default",
		Namespace: w.Namespace,
	}

	err := w.Ctx.HubCtrlClient().Get(context.Background(), key, objMCSB)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	err = w.Ctx.HubCtrlClient().Delete(context.Background(), objMCSB)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}