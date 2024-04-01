package util

import (
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Config struct {
	Clusters map[string]struct {
		KubeconfigPath string `mapstructure:"kubeconfigpath" required:"true"`
	} `mapstructure:"clusters" required:"true"`
}

type Cluster struct {
	K8sClientSet  *kubernetes.Clientset
	DynamicClient *dynamic.DynamicClient
	CtrlClient    client.Client
}

type Clusters map[string]*Cluster

type TestContext struct {
	Config          *Config
	Clusters        Clusters
	Log             logr.Logger
	ManagedClusters map[string]string
}

func GetClientSetFromKubeConfigPath(kubeconfigPath string) (*kubernetes.Clientset, *dynamic.DynamicClient, client.Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, nil, nil, err
	}

	k8sClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	ctrlClient, err := client.New(config, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, nil, nil, err
	}

	return k8sClientSet, dynamicClient, ctrlClient, nil
}

func (ctx *TestContext) HubK8sClientSet() *kubernetes.Clientset {
	return ctx.Clusters["hub"].K8sClientSet
}

func (ctx *TestContext) HubDynamicClient() *dynamic.DynamicClient {
	return ctx.Clusters["hub"].DynamicClient
}

func (ctx *TestContext) HubCtrlClient() client.Client {
	return ctx.Clusters["hub"].CtrlClient
}

func (ctx *TestContext) HubKubeconfig() string {
	return ctx.Config.Clusters["hub"].KubeconfigPath
}

func (ctx *TestContext) C1K8sClientSet() *kubernetes.Clientset {
	return ctx.Clusters["c1"].K8sClientSet
}

func (ctx *TestContext) C1DynamicClient() *dynamic.DynamicClient {
	return ctx.Clusters["c1"].DynamicClient
}

func (ctx *TestContext) C1CtrlClient() client.Client {
	return ctx.Clusters["c1"].CtrlClient
}

func (ctx *TestContext) C1Kubeconfig() string {
	return ctx.Config.Clusters["c1"].KubeconfigPath
}

func (ctx *TestContext) C2K8sClientSet() *kubernetes.Clientset {
	return ctx.Clusters["c2"].K8sClientSet
}

func (ctx *TestContext) C2DynamicClient() *dynamic.DynamicClient {
	return ctx.Clusters["c2"].DynamicClient
}

func (ctx *TestContext) C2CtrlClient() client.Client {
	return ctx.Clusters["c2"].CtrlClient
}

func (ctx *TestContext) C2Kubeconfig() string {
	return ctx.Config.Clusters["c2"].KubeconfigPath
}

func (ctx *TestContext) GetClusters() Clusters {
	return ctx.Clusters
}

func (ctx *TestContext) GetHubClusters() Clusters {
	hubClusters := make(Clusters)

	for clusterName, cluster := range ctx.Clusters {
		if strings.Contains(clusterName, "hub") {
			hubClusters[clusterName] = cluster
		}
	}

	return hubClusters
}

func (ctx *TestContext) GetManagedClusters() Clusters {
	managedClusters := make(Clusters)

	for clusterName, cluster := range ctx.Clusters {
		if !strings.Contains(clusterName, "hub") {
			managedClusters[clusterName] = cluster
		}
	}

	return managedClusters
}
