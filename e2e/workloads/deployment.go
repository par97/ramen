package workloads

import (
	"github.com/ramendr/ramen/e2e/util"
)

type Deployment struct {
	RepoURL          string // Possibly all this is part of Workload than each implementation of the interfaces?
	Path             string
	Revision         string
	Name             string // deployment-rbd
	Namespace        string // deployment-rbd
	PVCLabel         string // busybox
	PlacementName    string
	ChannelNamespace string

	Ctx *util.TestContext
}

func (w *Deployment) Init() {
	w.RepoURL = "https://github.com/ramendr/ocm-ramen-samples.git"
	w.Path = "subscription/deployment-k8s-regional-rbd"
	w.Revision = "main"
	w.Name = "deployment-rbd"
	w.Namespace = "deployment-rbd"
	w.PVCLabel = "busybox"
	w.PlacementName = "placement"
}

func (w *Deployment) GetName() string {
	return w.Name
}

func (w *Deployment) GetNameSpace() string {
	return w.Namespace
}

func (w *Deployment) GetPVCLabel() string {
	return w.PVCLabel
}

func (w *Deployment) GetPlacementName() string {
	return w.PlacementName
}

func (w *Deployment) GetResourceURL() string {
	//by default the timeout is 27s, could fail sometimes
	return w.RepoURL + "/" + w.Path + "?ref=" + w.Revision + "&timeout=90s"
}

func (w *Deployment) Kustomize() error {
	w.Ctx.Log.Info("enter Deployment Kustomize")

	return nil

}

func (w *Deployment) GetResources() error {
	// this would be a common function given the vars? But we need the resources Kustomized
	w.Ctx.Log.Info("enter Deployment GetResources")
	return nil
}

func (w *Deployment) Health() error {
	// Check the workload health on a targetCluster
	w.Ctx.Log.Info("enter Deployment Health")
	return nil
}

func (w *Deployment) Deploy() error {
	err := w.createNamespace(w.ChannelNamespace)
	if err != nil {
		return err
	}
	err = w.createChannel()
	if err != nil {
		return err
	}
	err = w.createNamespace(w.Namespace)
	if err != nil {
		return err
	}
	err = w.createSubscription()
	if err != nil {
		return err
	}
	err = w.createPlacement()
	if err != nil {
		return err
	}
	err = w.createManagedClusterSetBinding()
	if err != nil {
		return err
	}

	return nil
}

func (w *Deployment) Undeploy() error {
	err := w.deleteManagedClusterSetBinding()
	if err != nil {
		return err
	}
	err = w.deletePlacement()
	if err != nil {
		return err
	}
	err = w.deleteSubscription()
	if err != nil {
		return err
	}
	err = w.deleteNamespace(w.Namespace)
	if err != nil {
		return err
	}
	err = w.deleteChannel()
	if err != nil {
		return err
	}
	err = w.deleteNamespace(w.ChannelNamespace)
	if err != nil {
		return err
	}

	return nil
}
