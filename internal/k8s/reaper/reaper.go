package reaper

import (
	"context"
	"fmt"
	"golang.org/x/exp/slog"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Reaper struct {
	Connector func() (kubernetes.Interface, error)
}

func (r *Reaper) connect() (*k8sClient, error) {
	if r.Connector == nil {
		return nil, fmt.Errorf("no connector specified")
	}

	var client k8sClient
	var err error
	if client.client, err = r.Connector(); err != nil {
		err = fmt.Errorf("could not connect to cluster: %w", err)
	}
	return &client, err
}

func (r *Reaper) Reap(ctx context.Context, namespace, name string) (int, error) {
	client, err := r.connect()
	if err != nil {
		return 0, fmt.Errorf("reap: %w", err)
	}

	var deleted int
	pods, err := client.getPods(ctx, namespace, name)
	if err != nil {
		return 0, err
	}
	slog.Debug("found pods", "count", len(pods))

	for _, pod := range pods {
		podLogger := slog.With("name", pod.GetName())

		podLogger.Debug("checking pod")
		var found bool
		var status coreV1.ConditionStatus
		for _, condition := range pod.Status.Conditions {
			if condition.Type == "Ready" {
				status = condition.Status
				found = true
				break
			}
		}

		podLogger.Debug("pod status", "status", string(status))

		if !found {
			podLogger.Debug("pod doesn't appear to be running")
			continue
		}

		if status == "True" {
			podLogger.Debug("pod is ready")
			continue
		}

		podLogger.Info("pod not ready. deleting ...")
		if err = client.client.CoreV1().Pods(namespace).Delete(ctx, pod.GetName(), metaV1.DeleteOptions{}); err != nil {
			break
		}

		deleted++
		podLogger.Info("pod deleted")
	}

	return deleted, err
}
