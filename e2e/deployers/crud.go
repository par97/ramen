package deployers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
	ocmclusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	ocmclusterv1beta2 "open-cluster-management.io/api/cluster/v1beta2"
	placementrulev1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"
	subscriptionv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/v1"
)

func (s *Subscription) createNamespace() error {

	objNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Namespace,
		},
	}

	err := s.Ctx.HubCtrlClient().Create(context.Background(), objNs)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("namespace " + s.Namespace + " already Exists")
	}
	return nil
}

func (s *Subscription) createSubscription(w workloads.Workload) error {

	labels := make(map[string]string)
	labels["apps"] = s.AppName

	annotations := make(map[string]string)
	annotations["apps.open-cluster-management.io/github-branch"] = w.GetRevision()
	annotations["apps.open-cluster-management.io/github-path"] = w.GetPath()

	// Construct PlacementRef
	objReplacementRef := corev1.ObjectReference{
		Kind: "Placement",
		Name: util.DefaultPlacement,
	}

	objPlacementRulePlacement := &placementrulev1.Placement{}
	objPlacementRulePlacement.PlacementRef = &objReplacementRef

	objSubscription := &subscriptionv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:        s.SubscriptionName,
			Namespace:   s.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: subscriptionv1.SubscriptionSpec{
			Channel:   s.ChannelNamespace + "/" + s.ChannelName,
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

func (s *Subscription) createPlacement() error {

	labels := make(map[string]string)
	labels["apps"] = s.AppName

	arrayClusterSets := []string{"default"}
	var numClusters int32 = 1

	objPlacement := &ocmclusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.DefaultPlacement,
			Namespace: s.Namespace,
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

func (s *Subscription) createManagedClusterSetBinding() error {

	labels := make(map[string]string)
	labels["apps"] = s.AppName

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: s.Namespace,
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

func (s *Subscription) deleteNamespace() error {

	objNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Namespace,
		},
	}
	err := s.Ctx.HubCtrlClient().Delete(context.Background(), objNs)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("namespace " + s.Namespace + " not found")
	}

	return nil
}

func (s *Subscription) deleteSubscription() error {

	objSubscription := &subscriptionv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.SubscriptionName,
			Namespace: s.Namespace,
		},
	}

	err := s.Ctx.HubCtrlClient().Delete(context.Background(), objSubscription)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("subscription " + s.SubscriptionName + " not found")
	}

	return nil
}

func (s *Subscription) deletePlacement() error {

	objPlacement := &ocmclusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.DefaultPlacement,
			Namespace: s.Namespace,
		},
	}

	err := s.Ctx.HubCtrlClient().Delete(context.Background(), objPlacement)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("placement " + util.DefaultPlacement + " not found")
	}
	return nil
}

func (s *Subscription) deleteManagedClusterSetBinding() error {

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: s.Namespace,
		},
	}

	err := s.Ctx.HubCtrlClient().Delete(context.Background(), objMCSB)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("managedClusterSetBinding default not found")
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
