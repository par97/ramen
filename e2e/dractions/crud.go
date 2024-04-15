package dractions

import (
	"context"
	"fmt"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getPlacement(ctrlClient client.Client, namespace, name string) (*clusterv1beta1.Placement, error) {
	placement := &clusterv1beta1.Placement{}
	key := types.NamespacedName{Namespace: namespace, Name: name}

	err := ctrlClient.Get(context.Background(), key, placement)
	if err != nil {
		return nil, err
	}

	return placement, nil
}

func updatePlacement(ctrlClient client.Client, placement *clusterv1beta1.Placement) error {
	err := ctrlClient.Update(context.Background(), placement)
	if err != nil {
		return fmt.Errorf("could not update placment")
	}

	return nil
}

func getPlacementDecision(ctrlClient client.Client, namespace, name string) (*clusterv1beta1.PlacementDecision, error) {
	placementDecision := &clusterv1beta1.PlacementDecision{}
	key := types.NamespacedName{Namespace: namespace, Name: name}

	err := ctrlClient.Get(context.Background(), key, placementDecision)
	if err != nil {
		return nil, err
	}

	return placementDecision, nil
}

func getDRPC(ctrlClient client.Client, namespace, name string) (*ramen.DRPlacementControl, error) {
	drpc := &ramen.DRPlacementControl{}
	key := types.NamespacedName{Namespace: namespace, Name: name}

	err := ctrlClient.Get(context.Background(), key, drpc)
	if err != nil {
		return nil, err
	}

	return drpc, nil
}

func (r DRActions) createDRPC(ctrlClient client.Client, drpc *ramen.DRPlacementControl) error {
	err := ctrlClient.Create(context.Background(), drpc)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return fmt.Errorf("could not create drpc " + drpc.Name)
		}

		r.Ctx.Log.Info("drpc " + drpc.Name + " already Exists")
	}

	return nil
}

func updateDRPC(ctrlClient client.Client, drpc *ramen.DRPlacementControl) error {
	err := ctrlClient.Update(context.Background(), drpc)
	if err != nil {
		return fmt.Errorf("could not update placement")
	}

	return nil
}

func getDRPolicy(ctrlClient client.Client, name string) (*ramen.DRPolicy, error) {
	drpolicy := &ramen.DRPolicy{}
	key := types.NamespacedName{Name: name}

	err := ctrlClient.Get(context.Background(), key, drpolicy)
	if err != nil {
		return nil, err
	}

	return drpolicy, nil
}

func deleteDRPC(ctrlClient client.Client, namespace, name string) error {
	objDrpc := &ramen.DRPlacementControl{}
	key := types.NamespacedName{Namespace: namespace, Name: name}

	err := ctrlClient.Get(context.Background(), key, objDrpc)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}

	err = ctrlClient.Delete(context.Background(), objDrpc)
	if err != nil {
		return err
	}

	return nil
}
