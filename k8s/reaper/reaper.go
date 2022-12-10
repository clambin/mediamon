package reaper

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Reaper struct {
	Connector func() kubernetes.Interface
}

func (r *Reaper) connect() (*k8sClient, error) {
	if r.Connector != nil {
		client := k8sClient{client: r.Connector()}
		return &client, nil
	}

	return connect()
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
	log.Debugf("found %d pods", len(pods))

	for _, pod := range pods {
		log.Debugf("checking pod %s ...", pod.GetName())
		var found bool
		var status coreV1.ConditionStatus
		for _, condition := range pod.Status.Conditions {
			if condition.Type == "Ready" {
				status = condition.Status
				found = true
				break
			}
		}

		log.Debugf("%s - Ready: %s", pod.GetName(), status)

		if !found {
			log.Debugf("%s doesn't appear to be running", pod.GetName())
			continue
		}
		if status == "True" {
			log.Debugf("%s is ready", pod.GetName())
			continue
		}

		log.Infof("pod %s is not ready. deleting ...", pod.GetName())
		err = client.client.CoreV1().Pods(namespace).Delete(ctx, pod.GetName(), metaV1.DeleteOptions{})
		if err != nil {
			break
		}
		deleted++
		log.Info("pod deleted")
	}

	return deleted, err
}
