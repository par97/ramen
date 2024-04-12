package deployers

import (
	"context"
	"fmt"
	"os/exec"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
	ocmclusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	ocmclusterv1beta2 "open-cluster-management.io/api/cluster/v1beta2"
	placementrulev1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"
	subscriptionv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/v1"

	argocdv1alpha1hack "github.com/ramendr/ramen/e2e/argocd"
)

func (a *ApplicationSet) createApplicationSet(w workloads.Workload) error {
	a.Ctx.Log.Info("enter createApplicationSet")
	var requeueSeconds int64 = 180

	appset := &argocdv1alpha1hack.ApplicationSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
		},
		Spec: argocdv1alpha1hack.ApplicationSetSpec{
			Generators: []argocdv1alpha1hack.ApplicationSetGenerator{
				{
					ClusterDecisionResource: &argocdv1alpha1hack.DuckTypeGenerator{
						ConfigMapRef: a.ClusterDecisionConfigMapName,
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
						Namespace: a.ApplicationDestinationNamespace,
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

	err := a.Ctx.HubCtrlClient().Create(context.Background(), appset)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		a.Ctx.Log.Info("applicationset " + appset.Name + " already Exists")
	}

	return nil
}

func (a *ApplicationSet) deleteApplicationSet() error {
	a.Ctx.Log.Info("enter deleteApplicationSet")

	appset := &argocdv1alpha1hack.ApplicationSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
		},
	}

	err := a.Ctx.HubCtrlClient().Delete(context.Background(), appset)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		a.Ctx.Log.Info("applicationset " + appset.Name + " not found")
	}
	return nil
}

func createPlacementDecisionConfigMap(ctx *util.TestContext, cmName string, cmNamespace string) error {

	object := metav1.ObjectMeta{Name: cmName, Namespace: cmNamespace}

	data := map[string]string{
		"apiVersion":    "cluster.open-cluster-management.io/v1beta1",
		"kind":          "placementdecisions",
		"statusListKey": "decisions",
		"matchKey":      "clusterName",
	}

	configMap := &corev1.ConfigMap{ObjectMeta: object, Data: data}

	err := ctx.HubCtrlClient().Create(context.Background(), configMap)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return fmt.Errorf("could not create configMap " + cmName)
		}
		ctx.Log.Info("configMap " + cmName + " already Exists")
	}
	return nil
}

func deleteConfigMap(ctx *util.TestContext, cmName string, cmNamespace string) error {

	object := metav1.ObjectMeta{Name: cmName, Namespace: cmNamespace}

	configMap := &corev1.ConfigMap{
		ObjectMeta: object,
	}

	err := ctx.HubCtrlClient().Delete(context.Background(), configMap)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return fmt.Errorf("could not delete configMap " + cmName)
		}
		ctx.Log.Info("configMap " + cmName + " not found")
	}

	return nil
}

// argocd command currently has bug to honor the --namespace param
// see https://github.com/argoproj/argo-cd/issues/14167
func (a *ApplicationSet) addArgoCDClusters() error {
	//TODO: clusternames better to be dynamically got from config
	for _, c := range util.ClusterNames {
		a.Ctx.Log.Info("add cluster " + c + " into ArgoCD")
		cmd := exec.Command("argocd", "cluster", "add", c, " -y --namespace argocd --kubeconfig "+a.Ctx.HubKubeconfig())
		out, err := util.RunCommand(cmd)
		if err != nil {
			return err
		}
		a.Ctx.Log.Info(out)
	}
	return nil
}

// argocd command currently has bug to honor the --namespace param
// see https://github.com/argoproj/argo-cd/issues/14167
func (a *ApplicationSet) deleteArgoCDClusters() error {
	for _, c := range util.ClusterNames {
		a.Ctx.Log.Info("delete cluster " + c + " from ArgoCD")
		cmd := exec.Command("argocd", "cluster", "rm", c, " -y")
		out, err := util.RunCommand(cmd)
		if err != nil {
			return err
		}
		a.Ctx.Log.Info(out)
	}
	return nil
}

func createNamespace(ctx *util.TestContext, namespace string) error {

	objNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	err := ctx.HubCtrlClient().Create(context.Background(), objNs)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		ctx.Log.Info("namespace " + namespace + " already Exists")
	}
	return nil
}

func (s *Subscription) createSubscription(w workloads.Workload) error {

	labels := make(map[string]string)
	labels["apps"] = s.Name

	annotations := make(map[string]string)
	annotations["apps.open-cluster-management.io/github-branch"] = w.GetRevision()
	annotations["apps.open-cluster-management.io/github-path"] = w.GetPath()

	// Construct PlacementRef
	objReplacementRef := corev1.ObjectReference{
		Kind: "Placement",
		Name: s.PlacementName,
	}

	objPlacementRulePlacement := &placementrulev1.Placement{}
	objPlacementRulePlacement.PlacementRef = &objReplacementRef

	objSubscription := &subscriptionv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:        s.Name,
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
		s.Ctx.Log.Info("placement " + objSubscription.Name + " already Exists")
	}

	return nil
}

func createPlacement(ctx *util.TestContext, plName string, plNamespace string, appName string) error {

	labels := make(map[string]string)
	labels["apps"] = appName

	arrayClusterSets := []string{"default"}
	var numClusters int32 = 1

	objPlacement := &ocmclusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      plName,      // util.DefaultPlacement,
			Namespace: plNamespace, // s.Namespace,
			Labels:    labels,
		},
		Spec: ocmclusterv1beta1.PlacementSpec{
			ClusterSets:      arrayClusterSets,
			NumberOfClusters: &numClusters,
		},
	}

	err := ctx.HubCtrlClient().Create(context.Background(), objPlacement)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		ctx.Log.Info("placement " + objPlacement.Name + " already Exists")
	}
	return nil
}

func createManagedClusterSetBinding(ctx *util.TestContext, mcsbName string, mcsbNamespace string, appName string) error {

	labels := make(map[string]string)
	labels["apps"] = appName // s.AppName

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mcsbName,      // "default",
			Namespace: mcsbNamespace, // s.Namespace,
			Labels:    labels,
		},
		Spec: ocmclusterv1beta2.ManagedClusterSetBindingSpec{
			ClusterSet: "default",
		},
	}

	err := ctx.HubCtrlClient().Create(context.Background(), objMCSB)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		ctx.Log.Info("managedClusterSetBinding " + objMCSB.Name + " already Exists")
	}
	return nil
}

func deleteNamespace(ctx *util.TestContext, namespace string) error {

	objNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	err := ctx.HubCtrlClient().Delete(context.Background(), objNs)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		ctx.Log.Info("namespace " + namespace + " not found")
	}

	return nil
}

func (s *Subscription) deleteSubscription() error {

	objSubscription := &subscriptionv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
	}

	err := s.Ctx.HubCtrlClient().Delete(context.Background(), objSubscription)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		s.Ctx.Log.Info("subscription " + s.Name + " not found")
	}

	return nil
}

func deletePlacement(ctx *util.TestContext, plName string, plNamespace string) error {

	objPlacement := &ocmclusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      plName,      // util.DefaultPlacement,
			Namespace: plNamespace, //s.Namespace,
		},
	}

	err := ctx.HubCtrlClient().Delete(context.Background(), objPlacement)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
			return err
		}
		ctx.Log.Info("placement " + plName + " not found")
	}
	return nil
}

func deleteManagedClusterSetBinding(ctx *util.TestContext, mcsbName string, mcsbNamespace string) error {

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mcsbName,      // "default",
			Namespace: mcsbNamespace, // s.Namespace,
		},
	}

	err := ctx.HubCtrlClient().Delete(context.Background(), objMCSB)
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Printf("err: %v\n", err)
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
