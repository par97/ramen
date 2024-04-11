package deployers

import (
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

type ApplicationSet struct {
	Ctx *util.TestContext

	Name      string
	Namespace string

	// ArgoCDNamespace                 string
	PlacementName                   string
	McsbName                        string
	ClusterDecisionConfigMapName    string
	ApplicationDestinationNamespace string
}

func (a *ApplicationSet) Init() {
	a.Name = "appset"
	a.Namespace = "argocd"
	// a.ArgoCDNamespace = "argocd"
	a.PlacementName = a.Name + "-placement"
	a.McsbName = "default"
	a.ClusterDecisionConfigMapName = a.Name + "-configmap"
	a.ApplicationDestinationNamespace = a.Name

}

func (a ApplicationSet) Deploy(w workloads.Workload) error {
	// Generate a Placement for the Workload
	// Generate a Binding for the namespace?
	// Generate an ApplicationSet for the Workload
	// - Kustomize the Workload; call Workload.Kustomize(StorageType)
	// Address namespace/label/suffix as needed for various resources
	// w.Kustomize()
	a.Ctx.Log.Info("enter ApplicationSet Deploy")

	err := a.addArgoCDClusters()
	if err != nil {
		return err
	}

	err = createManagedClusterSetBinding(a.Ctx, a.McsbName, a.Namespace, a.Name)
	if err != nil {
		return err
	}

	err = createPlacement(a.Ctx, a.PlacementName, a.Namespace, a.Name)
	if err != nil {
		return err
	}

	err = createPlacementDecisionConfigMap(a.Ctx, a.ClusterDecisionConfigMapName, a.Namespace)
	if err != nil {
		return err
	}

	err = a.createApplicationSet(w)
	if err != nil {
		return err
	}

	return err
}

func (a ApplicationSet) Undeploy(w workloads.Workload) error {
	// Delete Placement, Binding, ApplicationSet
	a.Ctx.Log.Info("enter ApplicationSet Undeploy")

	err := a.deleteApplicationSet()
	if err != nil {
		return err
	}

	err = deleteConfigMap(a.Ctx, a.ClusterDecisionConfigMapName, a.Namespace)
	if err != nil {
		return err
	}

	err = deletePlacement(a.Ctx, a.PlacementName, a.Namespace)
	if err != nil {
		return err
	}

	err = deleteManagedClusterSetBinding(a.Ctx, a.McsbName, a.Namespace)
	if err != nil {
		return err
	}

	// don't use, this function is problematic
	// err := a.deleteArgoCDClusters()

	return nil
}

func (a ApplicationSet) GetName() string {
	return a.Name
}

func (a ApplicationSet) GetNameSpace() string {
	return a.Namespace
}

func (a ApplicationSet) Health(w workloads.Workload) error {
	a.Ctx.Log.Info("enter ApplicationSet Health")
	w.GetResources()
	// Check health using reflection to known types of the workload on the targetCluster
	// Again if using reflection can be a common function outside of deployer as such
	return nil
}
