package workloads

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ramendr/ramen/e2e/util"
)

type Deployment struct {
	RepoURL       string // Possibly all this is part of Workload than each implementation of the interfaces?
	Path          string
	Revision      string
	Ctx           *util.TestContext
	Name          string // deployment-rbd
	Namespace     string // deployment-rbd
	PVCLabel      string // busybox
	PlacementName string
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
	//run: kustomize build "url"
	cmd := exec.Command("kustomize", "build", w.GetResourceURL())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command failed")
	}

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
