package dractions

import (
	"fmt"
	"time"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	"github.com/ramendr/ramen/e2e/util"
	"open-cluster-management.io/api/cluster/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// return placement object, placementDecisionName, error
func waitPlacementDecision(client client.Client, namespace string, placementName string) (*v1beta1.Placement, string, error) {

	timeout := 300 //seconds
	interval := 30 //seconds
	startTime := time.Now()
	placementDecisionName := ""

	for {
		placement, err := getPlacement(client, namespace, placementName)
		if err != nil {

			return nil, "", err
		}
		for _, cond := range placement.Status.Conditions {
			if cond.Type == "PlacementSatisfied" && cond.Status == "True" {
				placementDecisionName = placement.Status.DecisionGroups[0].Decisions[0]
				if placementDecisionName != "" {
					util.Ctx.Log.Info("got placementdecision name " + placementDecisionName)
					return placement, placementDecisionName, nil
				}
			}
		}
		if time.Since(startTime) > time.Second*time.Duration(timeout) {
			fmt.Println("could not get placement decision before timeout")
			return nil, "", fmt.Errorf("could not get placement decision before timeout")
		}
		util.Ctx.Log.Info(fmt.Sprintf("could not get placement decision, retry in %v seconds", interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func waitDRPCReady(client client.Client, namespace string, drpcName string) error {

	timeout := 300 //seconds
	interval := 30 //seconds
	startTime := time.Now()
	for {
		ready := true
		drpc, err := getDRPC(client, namespace, drpcName)
		if err != nil {

			return err
		}

		for _, cond := range drpc.Status.Conditions {
			if cond.Type == "Available" && cond.Status != "True" {
				util.Ctx.Log.Info("drpc status Available is not True")
				ready = false
				break
			}
			if cond.Type == "PeerReady" && cond.Status != "True" {
				util.Ctx.Log.Info("drpc status PeerReady is not True")
				ready = false
				break
			}
		}
		if ready {
			if drpc.Status.LastGroupSyncTime == nil {
				util.Ctx.Log.Info("drpc status LastGroupSyncTime is nil")
				ready = false
			}
		}
		if ready {
			util.Ctx.Log.Info("drpc status is ready")
			return nil
		}
		if time.Since(startTime) > time.Second*time.Duration(timeout) {
			fmt.Printf("drpc status is not ready yet before timeout of %v\n", timeout)
			return fmt.Errorf(fmt.Sprintf("drpc status is not ready yet before timeout of %v", timeout))
		}
		util.Ctx.Log.Info(fmt.Sprintf("drpc status is not ready yet, retry in %v seconds", interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func waitDRPCPhase(client client.Client, namespace string, drpcName string, phase string) error {

	timeout := 600 //seconds
	interval := 30 //seconds
	startTime := time.Now()
	for {
		drpc, err := getDRPC(client, namespace, drpcName)
		if err != nil {

			return err
		}
		currentPhase := string(drpc.Status.Phase)
		if currentPhase == phase {
			util.Ctx.Log.Info("drpc phase is " + phase)
			return nil
		}
		if time.Since(startTime) > time.Second*time.Duration(timeout) {
			fmt.Printf("drpc phase is not %s yet before timeout of %v\n", phase, timeout)
			return fmt.Errorf(fmt.Sprintf("drpc status is not %s yet before timeout of %v", phase, timeout))
		}
		util.Ctx.Log.Info(fmt.Sprintf("current drpc phase is %s, expecting %s, retry in %v seconds", currentPhase, phase, interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func getCurrentCluster(client client.Client, namespace string, placementName string) (string, error) {

	_, placementDecisionName, err := waitPlacementDecision(client, namespace, placementName)
	if err != nil {
		return "", err
	}

	util.Ctx.Log.Info("get placementdecision " + placementDecisionName)
	placementDecision, err := getPlacementDecision(client, namespace, placementDecisionName)
	if err != nil {
		return "", err
	}

	clusterName := placementDecision.Status.Decisions[0].ClusterName
	util.Ctx.Log.Info("placementdecision clusterName: " + clusterName)

	return clusterName, nil
}

func getTargetCluster(client client.Client, namespace string, placementName string, drpolicy *ramen.DRPolicy) (string, error) {
	currentCluster, err := getCurrentCluster(client, namespace, placementName)
	if err != nil {
		return "", err
	}

	targetCluster := ""

	if currentCluster == drpolicy.Spec.DRClusters[0] {
		targetCluster = drpolicy.Spec.DRClusters[1]
	} else {
		targetCluster = drpolicy.Spec.DRClusters[0]
	}

	return targetCluster, nil
}
