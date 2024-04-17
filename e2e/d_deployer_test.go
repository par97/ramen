package e2e_test

// Deployer interface has methods to deploy a workload to a cluster
type Deployer interface {
	Init()
	Deploy(Workload) error
	Undeploy(Workload) error
	// Scale(Workload) for adding/removing PVCs; in Deployer even though scaling is a Workload interface
	// as we can Kustomize the Workload and change the deployer to perform the right action
	// Resize(Workload) for changing PVC(s) size
	// Health(Workload) error
	GetNamePrefix() string
	GetID() string
}
