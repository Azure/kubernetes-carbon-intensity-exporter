package provider

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

const (
	//TODO: make it a configurable option.
	patrolInterval = time.Second * 1200
)

type Provider struct {
	clusterClient clientset.Interface

	recorder record.EventRecorder
}

func New(clusterClient clientset.Interface, recorder record.EventRecorder) (*Provider, error) {
	b := &Provider{
		clusterClient: clusterClient,
		recorder:      recorder,
	}
	return b, nil
}

func (b *Provider) Run(stopChan <-chan struct{}) {
	go wait.Until(b.Patrol, patrolInterval, stopChan)
}

func (b *Provider) Patrol() {

	/*
	   Calling SDK to get 24 hours foreceast data,
	*/

	b.recorder.Eventf(&corev1.ObjectReference{
		Kind:      "Pod",
		Namespace: "kube-system",
		Name:      "carbon-data-provider", // TODO: replace this with the actual Pod name, passed through the downward API.
	}, corev1.EventTypeNormal, "Provider results", "Done retrieve data")

}
