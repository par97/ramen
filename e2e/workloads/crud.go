package workloads

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (w *Deployment) createNamespace(namespace string) error {
	w.Ctx.Log.Info("enter Deployment createNamespace " + namespace)

	objNs := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NameSpace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	_, err := w.Ctx.HubK8sClientSet().CoreV1().Namespaces().Create(context.Background(), objNs, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (w *Deployment) createChannel() error {
	w.Ctx.Log.Info("enter Deployment createChannel")
	return nil
}

func (w *Deployment) createSubscription() error {
	w.Ctx.Log.Info("enter Deployment createSubscription")
	return nil
}

func (w *Deployment) createPlacement() error {
	w.Ctx.Log.Info("enter Deployment createPlacement")
	return nil
}

func (w *Deployment) createManagedClusterSetBinding() error {
	w.Ctx.Log.Info("enter Deployment createManagedClusterSetBinding")
	return nil
}

func (w *Deployment) deleteNamespace(namespace string) error {
	w.Ctx.Log.Info("enter Deployment deleteNamespace " + namespace)

	err := w.Ctx.HubK8sClientSet().CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	return nil
}

func (w *Deployment) deleteChannel() error {
	w.Ctx.Log.Info("enter Deployment deleteChannel")
	return nil
}

func (w *Deployment) deleteSubscription() error {
	w.Ctx.Log.Info("enter Deployment deleteSubscription")
	return nil
}

func (w *Deployment) deletePlacement() error {
	w.Ctx.Log.Info("enter Deployment deletePlacement")
	return nil
}

func (w *Deployment) deleteManagedClusterSetBinding() error {
	w.Ctx.Log.Info("enter Deployment deleteManagedClusterSetBinding")
	return nil
}
