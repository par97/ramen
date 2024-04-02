package deployers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/ramendr/ramen/e2e/workloads"
	ocmclusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	ocmclusterv1beta2 "open-cluster-management.io/api/cluster/v1beta2"
	placementrulev1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"
	subscriptionv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/v1"
)

func (s *Subscription) createNamespace(namespace string) error {
	s.Ctx.Log.Info("enter Deployment createNamespace " + namespace)

	objNs := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NameSpace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	_, err := s.Ctx.HubK8sClientSet().CoreV1().Namespaces().Create(context.Background(), objNs, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("namespace " + namespace + " already Exists")
	}

	return nil
}

func (s *Subscription) createSubscription(w workloads.Workload) error {
	s.Ctx.Log.Info("enter Deployment createSubscription")

	labels := make(map[string]string)
	labels["apps"] = w.GetName()

	annotations := make(map[string]string)
	annotations["apps.open-cluster-management.io/github-branch"] = w.GetRevision()
	annotations["apps.open-cluster-management.io/github-path"] = w.GetPath()

	// Construct PlacementRef
	objReplacementRef := corev1.ObjectReference{
		Kind: "Placement",
		Name: w.GetPlacementName(),
	}

	objPlacementRulePlacement := &placementrulev1.Placement{}
	objPlacementRulePlacement.PlacementRef = &objReplacementRef

	objSubscription := &subscriptionv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "subscription",
			Namespace:   w.GetNameSpace(),
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: subscriptionv1.SubscriptionSpec{
			Channel:   s.ChannelName + "/" + s.ChannelName,
			Placement: objPlacementRulePlacement,
		},
	}

	err := s.Ctx.HubCtrlClient().Create(context.Background(), objSubscription)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("placement " + objSubscription.ObjectMeta.Name + " already Exists")
	}

	return nil
}

func (s *Subscription) createPlacement(w workloads.Workload) error {
	s.Ctx.Log.Info("enter Deployment createPlacement")

	labels := make(map[string]string)
	labels["apps"] = w.GetName()

	arrayClusterSets := []string{"default"}
	var numClusters int32 = 1

	objPlacement := &ocmclusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      w.GetPlacementName(),
			Namespace: w.GetNameSpace(),
			Labels:    labels,
		},
		Spec: ocmclusterv1beta1.PlacementSpec{
			ClusterSets:      arrayClusterSets,
			NumberOfClusters: &numClusters,
		},
	}

	err := s.Ctx.HubCtrlClient().Create(context.Background(), objPlacement)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("placement " + objPlacement.ObjectMeta.Name + " already Exists")
	}

	return nil
}

func (s *Subscription) createManagedClusterSetBinding(w workloads.Workload) error {
	s.Ctx.Log.Info("enter Deployment createManagedClusterSetBinding")

	labels := make(map[string]string)
	labels["apps"] = w.GetName()

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: w.GetNameSpace(),
			Labels:    labels,
		},
		Spec: ocmclusterv1beta2.ManagedClusterSetBindingSpec{
			ClusterSet: "default",
		},
	}

	err := s.Ctx.HubCtrlClient().Create(context.Background(), objMCSB)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("managedClusterSetBinding " + objMCSB.ObjectMeta.Name + " already Exists")
	}

	return nil
}

func (s *Subscription) deleteNamespace(namespace string) error {
	s.Ctx.Log.Info("enter Deployment deleteNamespace " + namespace)

	err := s.Ctx.HubK8sClientSet().CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("namespace " + namespace + " not found")
	}
	return nil
}

func (s *Subscription) deleteSubscription(w workloads.Workload) error {
	s.Ctx.Log.Info("enter Deployment deleteSubscription")

	objSubscription := &subscriptionv1.Subscription{}

	key := types.NamespacedName{
		Name:      w.GetSubscriptionName(),
		Namespace: w.GetNameSpace(),
	}

	err := s.Ctx.HubCtrlClient().Get(context.Background(), key, objSubscription)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("subscription " + w.GetSubscriptionName() + " not found")
		return nil
	}

	err = s.Ctx.HubCtrlClient().Delete(context.Background(), objSubscription)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (s *Subscription) deletePlacement(w workloads.Workload) error {
	s.Ctx.Log.Info("enter Deployment deletePlacement")

	objPlacement := &ocmclusterv1beta1.Placement{}

	key := types.NamespacedName{
		Name:      w.GetPlacementName(),
		Namespace: w.GetNameSpace(),
	}

	err := s.Ctx.HubCtrlClient().Get(context.Background(), key, objPlacement)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("placement " + w.GetPlacementName() + " not found")
		return nil
	}

	err = s.Ctx.HubCtrlClient().Delete(context.Background(), objPlacement)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (s *Subscription) deleteManagedClusterSetBinding(w workloads.Workload) error {
	s.Ctx.Log.Info("enter Deployment deleteManagedClusterSetBinding")

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{}

	key := types.NamespacedName{
		Name:      "default",
		Namespace: w.GetNameSpace(),
	}

	err := s.Ctx.HubCtrlClient().Get(context.Background(), key, objMCSB)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("managedClusterSetBinding " + "default" + " not found")
		return nil
	}

	err = s.Ctx.HubCtrlClient().Delete(context.Background(), objMCSB)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

/*
func (w *Deployment) createChannel() error {
	s.Ctx.Log.Info("enter Deployment createChannel")

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

	err := s.Ctx.HubCtrlClient().Create(context.Background(), objChannel)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (s *Subscription) deleteChannel() error {
	s.Ctx.Log.Info("enter Deployment deleteChannel")

	objChannel := &channelv1.Channel{}

	key := types.NamespacedName{
		Name:      w.ChannelName,
		Namespace: w.ChannelNamespace,
	}

	err := s.Ctx.HubCtrlClient().Get(context.Background(), key, objChannel)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("channel " + w.ChannelName + " not found")
		return nil
	}

	err = s.Ctx.HubCtrlClient().Delete(context.Background(), objChannel)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}
*/
