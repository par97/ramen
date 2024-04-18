package e2e_test

type ApplicationSet struct {
	// Name      string
	Namespace string

	NamePrefix    string
	PlacementName string
	McsbName      string

	// ClusterDecisionConfigMapName    string
	// ApplicationDestinationNamespace string
}

func (a *ApplicationSet) Init() {
	a.NamePrefix = "appset-"
	// a.Name = "appset-" + w.GetAppName()
	// appset need be created in argocd ns by default
	a.Namespace = "argocd"
	// a.PlacementName = a.Name
	a.McsbName = "default"
	// a.ClusterDecisionConfigMapName = a.Name
	// a.ApplicationDestinationNamespace = a.Name
}

func (a ApplicationSet) GetID() string {
	return "ApplicationSet"
}

func (a ApplicationSet) GetNamePrefix() string {
	return a.NamePrefix
}

func (a ApplicationSet) Deploy(w Workload) error {
	// Generate a Placement for the Workload
	// Generate a Binding for the namespace?
	// Generate an ApplicationSet for the Workload
	// - Kustomize the Workload; call Workload.Kustomize(StorageType)
	// Address namespace/label/suffix as needed for various resources
	// w.Kustomize()
	ctx.Log.Info("enter ApplicationSet Deploy")

	name := a.NamePrefix + w.GetAppName()
	namespace := name

	err := createManagedClusterSetBinding(a.McsbName, namespace, name)
	if err != nil {
		return err
	}

	err = createPlacement(name, namespace)
	if err != nil {
		return err
	}

	err = createPlacementDecisionConfigMap(name, namespace)
	if err != nil {
		return err
	}

	err = createApplicationSet(a, w)
	if err != nil {
		return err
	}

	return err
}

func (a ApplicationSet) Undeploy(w Workload) error {
	// Delete Placement, Binding, ApplicationSet
	ctx.Log.Info("enter ApplicationSet Undeploy")

	name := a.NamePrefix + w.GetAppName()
	namespace := name

	err := deleteApplicationSet(a, w)
	if err != nil {
		return err
	}

	err = deleteConfigMap(name, namespace)
	if err != nil {
		return err
	}

	err = deletePlacement(name, namespace)
	if err != nil {
		return err
	}

	err = deleteManagedClusterSetBinding(a.McsbName, namespace)
	if err != nil {
		return err
	}

	// don't use, this function is problematic
	// err := a.deleteArgoCDClusters()

	return nil
}

// func (a ApplicationSet) GetName() string {
// 	return a.Name
// }

// func (a ApplicationSet) GetNameSpace() string {
// 	return a.Namespace
// }

func (a ApplicationSet) Health(w Workload) error {
	ctx.Log.Info("enter ApplicationSet Health")
	w.GetResources()
	// Check health using reflection to known types of the workload on the targetCluster
	// Again if using reflection can be a common function outside of deployer as such
	return nil
}
