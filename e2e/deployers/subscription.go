package deployers

import (
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

type Subscription struct {
	NamePrefix string
	McsbName   string

	ChannelName      string
	ChannelNamespace string
}

func (s *Subscription) Init() {
	s.NamePrefix = "sub-"
	s.McsbName = "default"
	s.ChannelName = "ramen-gitops"
	s.ChannelNamespace = "ramen-samples"
}

func (s Subscription) GetID() string {
	return "Subscription"
}

func (s Subscription) GetNamePrefix() string {
	return s.NamePrefix
}

func (s Subscription) Deploy(w workloads.Workload) error {
	// Generate a Placement for the Workload
	// Use the global Channel
	// Generate a Binding for the namespace (does this need clusters?)
	// Generate a Subscription for the Workload
	// - Kustomize the Workload; call Workload.Kustomize(StorageType)
	// Address namespace/label/suffix as needed for various resources
	util.Ctx.Log.Info("enter Subscription Deploy")

	name := s.NamePrefix + w.GetAppName()
	namespace := name
	// w.Kustomize()
	err := createNamespace(namespace)
	if err != nil {
		return err
	}

	err = createManagedClusterSetBinding(s.McsbName, namespace, name)
	if err != nil {
		return err
	}

	err = createPlacement(name, namespace)
	if err != nil {
		return err
	}

	err = createSubscription(s, w)
	if err != nil {
		return err
	}

	return nil
}

func (s Subscription) Undeploy(w workloads.Workload) error {
	// Delete Subscription, Placement, Binding
	util.Ctx.Log.Info("enter Subscription Undeploy")

	name := s.NamePrefix + w.GetAppName()
	namespace := name

	err := deleteSubscription(s, w)
	if err != nil {
		return err
	}
	err = deletePlacement(name, namespace)
	if err != nil {
		return err
	}
	err = deleteManagedClusterSetBinding(s.McsbName, namespace)
	if err != nil {
		return err
	}
	err = deleteNamespace(namespace)
	if err != nil {
		return err
	}

	return nil
}
