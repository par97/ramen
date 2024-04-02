package workloads

import (
	"github.com/ramendr/ramen/e2e/util"
)

type Deployment struct {
	RepoURL  string // Possibly all this is part of Workload than each implementation of the interfaces?
	Path     string
	Revision string
	// Name      string // deployment-rbd
	// Namespace string // deployment-rbd
	AppName string // busybox

	Ctx *util.TestContext
}

func (w *Deployment) Init() {
	w.RepoURL = "https://github.com/ramendr/ocm-ramen-samples.git"
	w.Path = "workloads/deployment/k8s-regional-rbd"
	w.Revision = "main"
	// w.Name = "deployment-rbd"
	// w.Namespace = "deployment-rbd"
	w.AppName = "busybox"
}

// func (w Deployment) GetName() string {
// 	return w.Name
// }

// func (w Deployment) GetNameSpace() string {
// 	return w.Namespace
// }

func (w Deployment) GetAppName() string {
	return w.AppName
}

func (w Deployment) GetRepoURL() string {
	return w.RepoURL
}

func (w Deployment) GetPath() string {
	return w.Path
}

func (w Deployment) GetRevision() string {
	return w.Revision
}

func (w Deployment) GetResourceURL() string {
	//by default the timeout is 27s, could fail sometimes
	return w.RepoURL + "/" + w.Path + "?ref=" + w.Revision + "&timeout=90s"
}

func (w Deployment) Kustomize() error {
	w.Ctx.Log.Info("enter Deployment Kustomize")

	return nil
}

func (w Deployment) GetResources() error {
	// this would be a common function given the vars? But we need the resources Kustomized
	w.Ctx.Log.Info("enter Deployment GetResources")
	return nil
}

func (w Deployment) Health() error {
	// Check the workload health on a targetCluster
	w.Ctx.Log.Info("enter Deployment Health")
	return nil
}
