// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const ramenSystemNamespace = "ramen-system"

func Validate(t *testing.T) {
	t.Helper()

	ctx.Log.Info(t.Name())

	t.Run("CheckRamenHubOperatorStatus", CheckRamenHubOperatorStatus)
	// t.Run("RamenSpokes", RamenSpoke)
	// t.Run("Ceph", Ceph)
}

func CheckRamenHubOperatorStatus(t *testing.T) {
	ctx.Log.Info(t.Name())

	isRunning, podName, err := CheckRamenHubPodRunningStatus(ctx.Hub.K8sClientSet)
	if err != nil {
		t.Error(err)
	}

	if isRunning {
		ctx.Log.Info("Ramen Hub Operator is running", "pod", podName)
	} else {
		t.Error("no running Ramen Hub Operator pod")
	}

	ctx.Log.Info("TestRamenHubOperatorStatus: Pass")
}

func CheckRamenHubPodRunningStatus(k8sClient *kubernetes.Clientset) (bool, string, error) {
	labelSelector := "app=ramen-hub"
	podIdentifier := "ramen-hub-operator"

	ramenNameSpace, err := GetRamenNameSpace(k8sClient)
	if err != nil {
		return false, "", err
	}

	return CheckPodRunningStatus(k8sClient, ramenNameSpace, labelSelector, podIdentifier)
}

func GetRamenNameSpace(k8sClient *kubernetes.Clientset) (string, error) {
	isOpenShift, err := IsOpenShiftCluster(k8sClient)
	if err != nil {
		return "", err
	}

	if isOpenShift {
		return "openshift-operators", nil
	}

	return ramenSystemNamespace, nil
}

// IsOpenShiftCluster checks if the given Kubernetes cluster is an OpenShift cluster.
// It returns true if the cluster is OpenShift, false otherwise, along with any error encountered.
func IsOpenShiftCluster(k8sClient *kubernetes.Clientset) (bool, error) {
	discoveryClient := k8sClient.Discovery()

	apiGroups, err := discoveryClient.ServerGroups()
	if err != nil {
		return false, err
	}

	for _, group := range apiGroups.Groups {
		if group.Name == "route.openshift.io" {
			return true, nil
		}
	}

	return false, nil
}

// CheckPodRunningStatus checks if there is at least one pod matching the labelSelector
// in the given namespace that is in the "Running" phase and contains the podIdentifier in its name.
func CheckPodRunningStatus(client *kubernetes.Clientset, namespace, labelSelector, podIdentifier string) (
	bool, string, error,
) {
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return false, "", fmt.Errorf("failed to list pods in namespace %s: %v", namespace, err)
	}

	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, podIdentifier) && pod.Status.Phase == "Running" {
			return true, pod.Name, nil
		}
	}

	return false, "", nil
}

func RamenSpoke(t *testing.T) {
	ctx.Log.Info(t.Name())
}

func Ceph(t *testing.T) {
	ctx.Log.Info(t.Name())
}
