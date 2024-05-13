// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package workloads

type DeploymentFS struct {
	RepoURL  string
	Path     string
	Revision string
	AppName  string
}

func (w *DeploymentFS) Init() {
	w.RepoURL = ""
	w.Path = "workloads/deployment/k8s-regional-cephfs"
	w.Revision = "main"
	w.AppName = "busybox"
}

func (w DeploymentFS) GetAppName() string {
	return w.AppName
}

func (w DeploymentFS) GetName() string {
	return "DeploymentFS"
}

func (w DeploymentFS) GetRepoURL() string {
	return w.RepoURL
}

func (w DeploymentFS) GetPath() string {
	return w.Path
}

func (w DeploymentFS) GetRevision() string {
	return w.Revision
}

func (w DeploymentFS) Kustomize() error {
	return nil
}

func (w DeploymentFS) GetResources() error {
	// this would be a common function given the vars? But we need the resources Kustomized
	return nil
}

func (w DeploymentFS) Health() error {
	// Check the workload health on a targetCluster
	return nil
}
