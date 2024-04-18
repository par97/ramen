package e2e_test

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	argocdv1alpha1hack "github.com/ramendr/ramen/e2e/argocd"
	ocmclusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	ocmclusterv1beta2 "open-cluster-management.io/api/cluster/v1beta2"
	placementrulev1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"
	subscriptionv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/v1"
)

func createApplicationSet(a ApplicationSet, w Workload) error {
	ctx.Log.Info("enter createApplicationSet")
	var requeueSeconds int64 = 180

	name := a.NamePrefix + w.GetAppName()
	appset := &argocdv1alpha1hack.ApplicationSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: a.Namespace,
		},
		Spec: argocdv1alpha1hack.ApplicationSetSpec{
			Generators: []argocdv1alpha1hack.ApplicationSetGenerator{
				{
					ClusterDecisionResource: &argocdv1alpha1hack.DuckTypeGenerator{
						ConfigMapRef: name,
						LabelSelector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"cluster.open-cluster-management.io/placement": a.PlacementName,
							},
						},
						RequeueAfterSeconds: &requeueSeconds,
					},
				},
			},
			Template: argocdv1alpha1hack.ApplicationSetTemplate{
				ApplicationSetTemplateMeta: argocdv1alpha1hack.ApplicationSetTemplateMeta{
					Name: "rbd-{{name}}",
				},
				Spec: argocdv1alpha1hack.ApplicationSpec{
					Source: &argocdv1alpha1hack.ApplicationSource{
						RepoURL:        w.GetRepoURL(),
						Path:           w.GetPath(),
						TargetRevision: w.GetRevision(),
					},
					Destination: argocdv1alpha1hack.ApplicationDestination{
						Server:    "{{server}}",
						Namespace: name,
					},
					Project: "default",
					SyncPolicy: &argocdv1alpha1hack.SyncPolicy{
						Automated: &argocdv1alpha1hack.SyncPolicyAutomated{
							Prune:    true,
							SelfHeal: true,
						},
						SyncOptions: []string{
							"CreateNamespace=true",
							"PruneLast=true",
						},
					},
				},
			},
		},
	}

	err := ctx.Hub.CtrlClient.Create(context.Background(), appset)
	if err != nil {
		if !errors.IsAlreadyExists(err) {

			return err
		}
		ctx.Log.Info("applicationset " + appset.Name + " already Exists")
	}

	return nil
}

func deleteApplicationSet(a ApplicationSet, w Workload) error {
	ctx.Log.Info("enter deleteApplicationSet")
	name := a.NamePrefix + w.GetAppName()
	appset := &argocdv1alpha1hack.ApplicationSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: a.Namespace,
		},
	}

	err := ctx.Hub.CtrlClient.Delete(context.Background(), appset)
	if err != nil {
		if !errors.IsNotFound(err) {

			return err
		}
		ctx.Log.Info("applicationset " + appset.Name + " not found")
	}
	return nil
}

func createPlacementDecisionConfigMap(cmName string, cmNamespace string) error {

	object := metav1.ObjectMeta{Name: cmName, Namespace: cmNamespace}

	data := map[string]string{
		"apiVersion":    "cluster.open-cluster-management.io/v1beta1",
		"kind":          "placementdecisions",
		"statusListKey": "decisions",
		"matchKey":      "clusterName",
	}

	configMap := &corev1.ConfigMap{ObjectMeta: object, Data: data}

	err := ctx.Hub.CtrlClient.Create(context.Background(), configMap)
	if err != nil {
		if !errors.IsAlreadyExists(err) {

			return fmt.Errorf("could not create configMap " + cmName)
		}
		ctx.Log.Info("configMap " + cmName + " already Exists")
	}
	return nil
}

func deleteConfigMap(cmName string, cmNamespace string) error {

	object := metav1.ObjectMeta{Name: cmName, Namespace: cmNamespace}

	configMap := &corev1.ConfigMap{
		ObjectMeta: object,
	}

	err := ctx.Hub.CtrlClient.Delete(context.Background(), configMap)
	if err != nil {
		if !errors.IsNotFound(err) {

			return fmt.Errorf("could not delete configMap " + cmName)
		}
		ctx.Log.Info("configMap " + cmName + " not found")
	}

	return nil
}

func createNamespace(namespace string) error {

	objNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	err := ctx.Hub.CtrlClient.Create(context.Background(), objNs)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		ctx.Log.Info("namespace " + namespace + " already Exists")
	}
	return nil
}

func createSubscription(s Subscription, w Workload) error {

	name := s.NamePrefix + w.GetAppName()
	namespace := name

	labels := make(map[string]string)
	labels["apps"] = name

	annotations := make(map[string]string)
	annotations["apps.open-cluster-management.io/github-branch"] = w.GetRevision()
	annotations["apps.open-cluster-management.io/github-path"] = w.GetPath()

	// Construct PlacementRef
	objReplacementRef := corev1.ObjectReference{
		Kind: "Placement",
		Name: name,
	}

	objPlacementRulePlacement := &placementrulev1.Placement{}
	objPlacementRulePlacement.PlacementRef = &objReplacementRef

	objSubscription := &subscriptionv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: subscriptionv1.SubscriptionSpec{
			Channel:   s.ChannelNamespace + "/" + s.ChannelName,
			Placement: objPlacementRulePlacement,
		},
	}

	err := ctx.Hub.CtrlClient.Create(context.Background(), objSubscription)
	if err != nil {
		if !errors.IsAlreadyExists(err) {

			return err
		}
		ctx.Log.Info("placement " + objSubscription.Name + " already Exists")
	}

	return nil
}

func createPlacement(plName string, plNamespace string) error {

	labels := make(map[string]string)
	labels["apps"] = plName

	arrayClusterSets := []string{"default"}
	var numClusters int32 = 1

	objPlacement := &ocmclusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      plName,      // util.DefaultPlacement,
			Namespace: plNamespace, // namespace,
			Labels:    labels,
		},
		Spec: ocmclusterv1beta1.PlacementSpec{
			ClusterSets:      arrayClusterSets,
			NumberOfClusters: &numClusters,
		},
	}

	err := ctx.Hub.CtrlClient.Create(context.Background(), objPlacement)
	if err != nil {
		if !errors.IsAlreadyExists(err) {

			return err
		}
		ctx.Log.Info("placement " + objPlacement.Name + " already Exists")
	}
	return nil
}

func createManagedClusterSetBinding(mcsbName string, mcsbNamespace string, appName string) error {

	labels := make(map[string]string)
	labels["apps"] = appName // s.AppName

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mcsbName,      // "default",
			Namespace: mcsbNamespace, // namespace,
			Labels:    labels,
		},
		Spec: ocmclusterv1beta2.ManagedClusterSetBindingSpec{
			ClusterSet: "default",
		},
	}

	err := ctx.Hub.CtrlClient.Create(context.Background(), objMCSB)
	if err != nil {
		if !errors.IsAlreadyExists(err) {

			return err
		}
		ctx.Log.Info("managedClusterSetBinding " + objMCSB.Name + " already Exists")
	}
	return nil
}

func deleteNamespace(namespace string) error {

	objNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	err := ctx.Hub.CtrlClient.Delete(context.Background(), objNs)
	if err != nil {
		if !errors.IsNotFound(err) {

			return err
		}
		ctx.Log.Info("namespace " + namespace + " not found")
	}

	return nil
}

func deleteSubscription(s Subscription, w Workload) error {

	name := s.NamePrefix + w.GetAppName()
	namespace := name

	objSubscription := &subscriptionv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	err := ctx.Hub.CtrlClient.Delete(context.Background(), objSubscription)
	if err != nil {
		if !errors.IsNotFound(err) {

			return err
		}
		ctx.Log.Info("subscription " + name + " not found")
	}

	return nil
}

func deletePlacement(plName string, plNamespace string) error {

	objPlacement := &ocmclusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      plName,      // util.DefaultPlacement,
			Namespace: plNamespace, //namespace,
		},
	}

	err := ctx.Hub.CtrlClient.Delete(context.Background(), objPlacement)
	if err != nil {
		if !errors.IsNotFound(err) {

			return err
		}
		ctx.Log.Info("placement " + plName + " not found")
	}
	return nil
}

func deleteManagedClusterSetBinding(mcsbName string, mcsbNamespace string) error {

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mcsbName,      // "default",
			Namespace: mcsbNamespace, // namespace,
		},
	}

	err := ctx.Hub.CtrlClient.Delete(context.Background(), objMCSB)
	if err != nil {
		if !errors.IsNotFound(err) {

			return err
		}
		ctx.Log.Info("managedClusterSetBinding default not found")
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

	err := s.Ctx.Hub.CtrlClient.Create(context.Background(), objChannel)
	if err != nil {

		return err
	}

	return nil
}

func (s *Subscription) deleteChannel() error {
	s.Ctx.Log.Info("enter Deployment deleteChannel")

	objChannel := &channelv1.Channel{}

	key := typenamespacedName{
		Name:      w.ChannelName,
		Namespace: w.ChannelNamespace,
	}

	err := s.Ctx.Hub.CtrlClient.Get(context.Background(), key, objChannel)
	if err != nil {
		if !errors.IsNotFound(err) {

			return err
		}
		s.Ctx.Log.Info("channel " + w.ChannelName + " not found")
		return nil
	}

	err = s.Ctx.Hub.CtrlClient.Delete(context.Background(), objChannel)
	if err != nil {

		return err
	}

	return nil
}
*/
