package reaper

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestReaper_Reap(t *testing.T) {
	deployment := &appsV1.DeploymentList{
		TypeMeta: metaV1.TypeMeta{Kind: "deployments", APIVersion: "v1"},
		Items:    []appsV1.Deployment{{ObjectMeta: metaV1.ObjectMeta{Namespace: "media", Name: "transmission", UID: "deployment-transmission"}}},
	}
	replicasets1 := &appsV1.ReplicaSetList{
		TypeMeta: metaV1.TypeMeta{Kind: "replicasets", APIVersion: "v1"},
		Items: []appsV1.ReplicaSet{
			{ObjectMeta: metaV1.ObjectMeta{Namespace: "media", Name: "transmission-1", UID: "rs-transmission-1", OwnerReferences: []metaV1.OwnerReference{{UID: "deployment-transmission"}}}},
		},
	}
	replicasets2 := &appsV1.ReplicaSetList{
		TypeMeta: metaV1.TypeMeta{Kind: "replicasets", APIVersion: "v1"},
		Items: []appsV1.ReplicaSet{
			{ObjectMeta: metaV1.ObjectMeta{Namespace: "media", Name: "transmission-1", UID: "rs-transmission-1", OwnerReferences: []metaV1.OwnerReference{{UID: "deployment-transmission"}}}},
			{ObjectMeta: metaV1.ObjectMeta{Namespace: "media", Name: "transmission-2", UID: "rs-transmission-2", OwnerReferences: []metaV1.OwnerReference{{UID: "deployment-transmission"}}}},
		},
	}
	pods1 := &coreV1.PodList{
		TypeMeta: metaV1.TypeMeta{Kind: "pods", APIVersion: "v1"},
		Items: []coreV1.Pod{
			{
				ObjectMeta: metaV1.ObjectMeta{Namespace: "media", Name: "transmission-1", UID: "pod-transmission-1", OwnerReferences: []metaV1.OwnerReference{{UID: "rs-transmission-1"}}},
				Status:     coreV1.PodStatus{Phase: "Running", Conditions: []coreV1.PodCondition{{Type: "Ready", Status: "True"}}},
			},
		},
	}
	pods2 := &coreV1.PodList{
		TypeMeta: metaV1.TypeMeta{Kind: "pods", APIVersion: "v1"},
		Items: []coreV1.Pod{
			{
				ObjectMeta: metaV1.ObjectMeta{Namespace: "media", Name: "transmission-1", UID: "pod-transmission-1", OwnerReferences: []metaV1.OwnerReference{{UID: "rs-transmission-1"}}},
				Status:     coreV1.PodStatus{Phase: "Running", Conditions: []coreV1.PodCondition{{Type: "Ready", Status: "False"}}},
			},
		},
	}
	pods3 := &coreV1.PodList{
		TypeMeta: metaV1.TypeMeta{Kind: "pods", APIVersion: "v1"},
		Items: []coreV1.Pod{
			{
				ObjectMeta: metaV1.ObjectMeta{Namespace: "media", Name: "transmission-1", UID: "pod-transmission-1", OwnerReferences: []metaV1.OwnerReference{{UID: "rs-transmission-1"}}},
				Status:     coreV1.PodStatus{Phase: "Running", Conditions: []coreV1.PodCondition{{Type: "Ready", Status: "True"}}},
			},
			{
				ObjectMeta: metaV1.ObjectMeta{Namespace: "media", Name: "transmission-2", UID: "pod-transmission-2", OwnerReferences: []metaV1.OwnerReference{{UID: "rs-transmission-1"}}},
				Status:     coreV1.PodStatus{Phase: "Running", Conditions: []coreV1.PodCondition{{Type: "Ready", Status: "False"}}},
			},
			{
				ObjectMeta: metaV1.ObjectMeta{Namespace: "media", Name: "transmission-3", UID: "pod-transmission-3", OwnerReferences: []metaV1.OwnerReference{{UID: "rs-transmission-1"}}},
				Status:     coreV1.PodStatus{Phase: "Pending", Conditions: []coreV1.PodCondition{{}}},
			},
		},
	}

	tests := []struct {
		name    string
		objects []runtime.Object
		pass    bool
		count   int
	}{
		{
			name:    "no deployments",
			objects: []runtime.Object{},
			pass:    false,
		},
		{
			name:    "no replica sets",
			objects: []runtime.Object{deployment},
			pass:    false,
		},
		{
			name:    "no pods",
			objects: []runtime.Object{deployment, replicasets1},
			pass:    true,
			count:   0,
		},
		{
			name:    "one running pod",
			objects: []runtime.Object{deployment, replicasets2, pods1},
			pass:    true,
			count:   0,
		},
		{
			name:    "one failing pod",
			objects: []runtime.Object{deployment, replicasets2, pods2},
			pass:    true,
			count:   1,
		},
		{
			name:    "one of many pods failing",
			objects: []runtime.Object{deployment, replicasets2, pods3},
			pass:    true,
			count:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Reaper{Connector: func() (kubernetes.Interface, error) {
				return fake.NewSimpleClientset(tt.objects...), nil
			}}
			count, err := r.Reap(context.Background(), "media", "transmission")
			if tt.pass {
				require.NoError(t, err)
				assert.Equal(t, tt.count, count)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
