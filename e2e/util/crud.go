package util

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNamespace(client client.Client, namespace string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	err := client.Create(context.Background(), ns)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}

		Ctx.Log.Info("namespace " + namespace + " already Exists")
	}

	return nil
}

func DeleteNamespace(client client.Client, namespace string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	err := client.Delete(context.Background(), ns)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		Ctx.Log.Info("namespace " + namespace + " not found")
	}

	return nil
}

// Problem: currently we must manually add an annotation to application’s namespace to make volsync work.
// See this link https://volsync.readthedocs.io/en/stable/usage/permissionmodel.html#controlling-mover-permissions
// Workaround: create ns in both drclusters and add annotation
func CreateNamespaceAndAddAnnotation(namespace string) error {
	if err := CreateNamespace(Ctx.C1.CtrlClient, namespace); err != nil {
		return err
	}

	if err := AddNamespaceAnnotationForVolSync(Ctx.C1.CtrlClient, namespace); err != nil {
		return err
	}

	if err := CreateNamespace(Ctx.C2.CtrlClient, namespace); err != nil {
		return err
	}

	return AddNamespaceAnnotationForVolSync(Ctx.C2.CtrlClient, namespace)
}

func AddNamespaceAnnotationForVolSync(client client.Client, namespace string) error {
	key := types.NamespacedName{Name: namespace}
	objNs := &corev1.Namespace{}

	if err := client.Get(context.Background(), key, objNs); err != nil {
		return err
	}

	annotations := objNs.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	annotations["volsync.backube/privileged-movers"] = "true"
	objNs.SetAnnotations(annotations)

	return client.Update(context.Background(), objNs)
}
