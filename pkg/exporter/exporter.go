/*
MIT License
Copyright (c) Microsoft Corporation.
*/
package exporter

import (
	"os"
	"time"

	"github.com/Azure/kubernetes-carbon-intensity-exporter/pkg/sdk/client"
	"github.com/antihax/optional"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
)

var (
	podName   = os.Getenv("POD_NAME")
	namespace = os.Getenv("NAMESPACE")
)

type Exporter struct {
	clusterClient clientset.Interface
	apiClient     *client.APIClient
	recorder      record.EventRecorder
}

func New(clusterClient clientset.Interface, apiClient *client.APIClient, recorder record.EventRecorder) (*Exporter, error) {
	b := &Exporter{
		clusterClient: clusterClient,
		apiClient:     apiClient,
		recorder:      recorder,
	}
	return b, nil
}

func (e *Exporter) Run(ctx context.Context, configmapName, region string, patrolInterval time.Duration, stopChan <-chan struct{}) {
	go wait.Until(func() {
		e.Patrol(ctx, configmapName, region)
	}, patrolInterval, stopChan)
}

func (e *Exporter) Patrol(ctx context.Context, configmapName, region string) {
	forecast, err := e.getCurrentForecastData(ctx, region)
	if err != nil {
		return
	}
	err = e.CreateOrUpdateConfigMap(ctx, configmapName, forecast)
	if err != nil {
		e.recorder.Eventf(&corev1.ObjectReference{
			Kind:      "Pod",
			Namespace: namespace,
			Name:      podName,
		}, corev1.EventTypeWarning, "Configmap Create", "Error while creating configmap")
		klog.Errorf("an error has occurred while creating %s configmap, %s", configmapName, err.Error())
		return
	}
	e.recorder.Eventf(&corev1.ObjectReference{
		Kind:      "Pod",
		Namespace: namespace,
		Name:      podName,
	}, corev1.EventTypeNormal, "Exporter results", "Done retrieve data")

}

func (e *Exporter) getCurrentForecastData(ctx context.Context, region string) ([]client.EmissionsForecastDto, error) {
	opt := &client.CarbonAwareApiGetCurrentForecastDataOpts{
		DataStartAt: optional.EmptyTime(),
		DataEndAt:   optional.EmptyTime(),
	}
	forecast, _, err := e.apiClient.CarbonAwareApi.GetCurrentForecastData(ctx,
		[]string{region}, opt)
	if err != nil {
		klog.ErrorS(err, "error while getting current forecast data")
		return nil, err
	}

	return forecast, nil
}
