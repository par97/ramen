package workloads

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ramendr/ramen/e2e/util"
)

type Deployment struct {
	RepoURL  string // Possibly all this is part of Workload than each implementation of the interfaces?
	Path     string
	Revision string
	Ctx      *util.TestContext
}

func (w Deployment) GetKustomizeURL() string {
	return w.RepoURL + "/" + w.Path + "?ref=" + w.Revision
}

func (w Deployment) Kustomize() error {
	w.Ctx.Log.Info("enter Deployment Kustomize")
	//run: kustomize build "url"
	cmd := exec.Command("kustomize", "build", w.GetKustomizeURL())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command failed")
	}

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
