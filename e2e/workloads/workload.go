package workloads

type Workload interface {
	Kustomize() error    // Can differ based on the workload, hence part of the Workload interface
	GetResources() error // Get the actual workload resources

	GetName() string
	GetNameSpace() string
	GetPVCLabel() string
	GetPlacementName() string

	GetRepoURL() string // Possibly all this is part of Workload than each implementation of the interfaces?
	GetPath() string
	GetRevision() string

	// GetChannelNamespace() string
	// GetChannelName() string
	// GetChannelType() channelv1.ChannelType

	GetSubscriptionName() string

	Init()
	// Deploy() error
	// Undeploy() error
}
