package e2e_test

type Workload interface {
	Kustomize() error    // Can differ based on the workload, hence part of the Workload interface
	GetResources() error // Get the actual workload resources

	GetID() string
	GetAppName() string

	GetRepoURL() string // Possibly all this is part of Workload than each implementation of the interfaces?
	GetPath() string
	GetRevision() string

	Init()
}
