package deployers

import (
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

type ApplicationSet struct {
	// repoURL  string // From the Workload?
	// path     string // From the Workload?
	// revision string // From the Workload?
	Ctx *util.TestContext

	AppName   string // deployment-rbd
	Namespace string // deployment-rbd
}

func (s *ApplicationSet) Init() {
	s.AppName = "deployment-rbd"
	s.Namespace = "deployment-rbd"
}

func (s ApplicationSet) GetAppName() string {
	return s.AppName
}

func (s ApplicationSet) GetNameSpace() string {
	return s.Namespace
}

func (s ApplicationSet) Deploy(w workloads.Workload) error {
	// Generate a Placement for the Workload
	// Generate a Binding for the namespace?
	// Generate an ApplicationSet for the Workload
	// - Kustomize the Workload; call Workload.Kustomize(StorageType)
	// Address namespace/label/suffix as needed for various resources
	w.Kustomize()

	// w.Kustomize()

	/*
		apiVersion: cluster.open-cluster-management.io/v1beta1
		kind: Placement
		metadata:
		  name: rbd-placement
		  namespace: openshift-gitops
		spec:
		  clusterSets:
		    - default
		  numberOfClusters: 1
	*/

	// err := s.createPlacement()
	// if err != nil {
	// 	return err
	// }

	/*
		apiVersion: argoproj.io/v1alpha1
		kind: ApplicationSet
		metadata:
		  name: rbd
		  namespace: openshift-gitops
		spec:
		  generators:
		    - clusterDecisionResource:
		        configMapRef: acm-placement
		        labelSelector:
		          matchLabels:
		            cluster.open-cluster-management.io/placement: rbd-placement
		        requeueAfterSeconds: 180
		  template:
		    metadata:
		      name: rbd-{{name}}
		      labels:
		        velero.io/exclude-from-backup: "true"
		    spec:
		      project: default
		      sources:
		        - repositoryType: git
		          repoURL: https://github.com/RamenDR/ocm-ramen-samples.git
		          targetRevision: main
		          path: workloads/deployment/k8s-regional-rbd
		      destination:
		        namespace: rbd
		        server: "{{server}}"
		      syncPolicy:
		        automated:
		          selfHeal: true
		          prune: true
		        syncOptions:
		          - CreateNamespace=true
		          - PruneLast=true
	*/
	err := s.createApplicationSet(w)
	if err != nil {
		return err
	}

	return nil
}

func (s ApplicationSet) Undeploy(w workloads.Workload) error {
	// Delete Placement, Binding, ApplicationSet
	return nil
}

func (s ApplicationSet) Health(w workloads.Workload) error {
	s.Ctx.Log.Info("enter Subscription Health")
	w.GetResources()
	// Check health using reflection to known types of the workload on the targetCluster
	// Again if using reflection can be a common function outside of deployer as such
	return nil
}
