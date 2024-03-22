package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/ramendr/ramen/e2e/util"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
)

func validateConfig(config *util.Config) error {
	if config.Clusters["hub"].KubeconfigPath == "" {
		return fmt.Errorf("failed to find hub cluster in configuration")
	}

	if config.Clusters["c1"].KubeconfigPath == "" {
		return fmt.Errorf("failed to find c1 cluster in configuration")
	}

	if config.Clusters["c2"].KubeconfigPath == "" {
		return fmt.Errorf("failed to find c2 cluster in configuration")
	}

	return nil
}

func readConfig() (*util.Config, error) {
	config := &util.Config{}

	viper.SetConfigFile("config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("failed to find configuration file: %v", err)
		}

		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	if err := viper.UnmarshalExact(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %v", err)
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %v", err)
	}

	return config, nil
}

func configContext(ctx *util.TestContext, config *util.Config) error {

	ctx.Config = config
	ctx.Clusters = make(util.Clusters)

	for clusterName, cluster := range config.Clusters {
		k8sClientSet, dynamicClient, err := util.GetClientSetFromKubeConfigPath(cluster.KubeconfigPath)
		if err != nil {
			return err
		}

		ctx.Clusters[clusterName] = &util.Cluster{
			K8sClientSet:  k8sClientSet,
			DynamicClient: dynamicClient,
		}
	}

	return setupManagedClustersMapping(ctx)
}

func setupManagedClustersMapping(ctx *util.TestContext) error {

	cmd := exec.Command("kubectl", "config", "view", "--minify", "--kubeconfig="+ctx.C1Kubeconfig(), "-o=jsonpath={.clusters[0].name}")
	err, c1name := util.RunCommand(cmd)
	if err != nil {
		return err
	}
	// fmt.Printf("out c1name: %v\n", c1name)

	cmd = exec.Command("kubectl", "config", "view", "--minify", "--kubeconfig="+ctx.C2Kubeconfig(), "-o=jsonpath={.clusters[0].name}")
	err, c2name := util.RunCommand(cmd)
	if err != nil {
		return err
	}
	// fmt.Printf("out c2name: %v\n", c2name)

	client := ctx.HubDynamicClient()

	resource := schema.GroupVersionResource{Group: "cluster.open-cluster-management.io", Version: "v1", Resource: "managedclusters"}
	resp, err := client.Resource(resource).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fmt.Errorf("could not list managedcluster")
	}
	respJson, err := resp.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not marshaljson")
	}
	list := clusterv1.ManagedClusterList{}
	err = json.Unmarshal(respJson, &list)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not unmarshaljson")
	}
	number := len(list.Items)
	if number != 2 {
		return fmt.Errorf("found %v managedclusters, not equal to 2", number)
	}

	mc1name := list.Items[0].ObjectMeta.Name
	mc2name := list.Items[1].ObjectMeta.Name

	// fmt.Printf("mc1name: %v\n", mc1name)
	// fmt.Printf("mc2name: %v\n", mc2name)
	// fmt.Printf("c1name: %v\n", c1name)
	// fmt.Printf("c2name: %v\n", c2name)

	ctx.ManagedClusters = make(map[string]string)

	if mc1name == c1name && mc2name == c2name {
		ctx.ManagedClusters[mc1name] = "c1"
		ctx.ManagedClusters[mc2name] = "c2"
	} else if mc1name == c2name && mc2name == c1name {
		ctx.ManagedClusters[mc1name] = "c2"
		ctx.ManagedClusters[mc2name] = "c1"
	} else {
		err := fmt.Errorf("could not find matched cluster name in kubeconfig and managedcluster")
		return err
	}

	return nil
}
