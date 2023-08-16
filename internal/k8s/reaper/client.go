package reaper

import (
	"context"
	"fmt"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"log/slog"
)

type k8sClient struct {
	client kubernetes.Interface
}

func (c *k8sClient) getPods(ctx context.Context, namespace, name string) ([]coreV1.Pod, error) {
	deployment, err := c.getDeployment(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	replicaSets, err := c.getReplicaSetsForDeployment(ctx, namespace, deployment)
	if err != nil {
		return nil, err
	}

	return c.getPodsForReplicaSets(ctx, namespace, replicaSets)
}

func (c *k8sClient) getDeployment(ctx context.Context, namespace, name string) (*appsV1.Deployment, error) {
	deployments, err := c.client.
		AppsV1().
		Deployments(namespace).
		List(ctx, metaV1.ListOptions{
			FieldSelector: "metadata.name=" + name,
		})
	if err != nil {
		return nil, fmt.Errorf("list deployments: %w", err)
	}
	if len(deployments.Items) != 1 {
		return nil, fmt.Errorf("no deployments found for %s", name)
	}
	return &deployments.Items[0], nil
}

func (c *k8sClient) getReplicaSetsForDeployment(ctx context.Context, namespace string, deployment *appsV1.Deployment) (map[types.UID]struct{}, error) {
	sets := make(map[types.UID]struct{})
	replicaSets, err := c.client.
		AppsV1().
		ReplicaSets(namespace).
		List(ctx, metaV1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list replicasets: %w", err)
	}
	for _, replicaset := range replicaSets.Items {
		for _, ownerReference := range replicaset.GetOwnerReferences() {
			if ownerReference.UID == deployment.GetUID() {
				slog.Debug("replica set found", "name", replicaset.GetName())
				sets[replicaset.GetUID()] = struct{}{}
			}
		}
	}
	if len(sets) == 0 {
		return nil, fmt.Errorf("no replicatesets found for %s", deployment.GetName())
	}
	return sets, nil
}

func (c *k8sClient) getPodsForReplicaSets(ctx context.Context, namespace string, replicaSets map[types.UID]struct{}) ([]coreV1.Pod, error) {
	var targets []coreV1.Pod
	pods, err := c.client.
		CoreV1().
		Pods(namespace).
		List(ctx, metaV1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}
	for _, pod := range pods.Items {
		for _, ownerReference := range pod.OwnerReferences {
			if _, ok := replicaSets[ownerReference.UID]; ok {
				slog.Debug("pod found", "name", pod.GetName())
				targets = append(targets, pod)
			}
		}
	}
	return targets, nil
}
