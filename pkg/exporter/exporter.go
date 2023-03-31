/*
MIT License
Copyright (c) Microsoft Corporation.
*/
package exporter

import (
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

const (
	//TODO: make it a configurable option.
	patrolInterval = time.Second * 1200
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

func (e *Exporter) Run(stopChan <-chan struct{}) {
	go wait.Until(e.Patrol, patrolInterval, stopChan)
}

func (e *Exporter) Patrol() {
	ctx := context.Background()
	//e.getEmissionData(ctx, "eastus")
	//e.getCarbonIntensity(ctx, "eastus")
	forecast, err := e.getCurrentForecastData(ctx, []string{"eastus"})
	if err != nil {
		return
	}
	err = e.CreateOrUpdateConfigMap(ctx, forecast)
	if err != nil {
		e.recorder.Eventf(&corev1.ObjectReference{
			Kind:      "Pod",
			Namespace: "",
			Name:      "carbon-data-exporter", // TODO: replace this with the actual Pod name, passed through the downward API.
		}, corev1.EventTypeWarning, "Configmap Create", "Error while creating configmap")
		klog.Errorf("an error has occurred while creating %s configmap, %s", configMapName, err.Error())
		return
	}
	e.recorder.Eventf(&corev1.ObjectReference{
		Kind:      "Pod",
		Namespace: "",
		Name:      "carbon-data-exporter", // TODO: replace this with the actual Pod name, passed through the downward API.
	}, corev1.EventTypeNormal, "Exporter results", "Done retrieve data")

}

func (e *Exporter) getCurrentForecastData(ctx context.Context, region []string) ([]client.EmissionsForecastDto, error) {
	opt := &client.CarbonAwareApiGetCurrentForecastDataOpts{
		DataStartAt: optional.EmptyTime(),
		DataEndAt:   optional.EmptyTime(),
	}
	forecast, _, err := e.apiClient.CarbonAwareApi.GetCurrentForecastData(ctx,
		region, opt)
	if err != nil {
		klog.ErrorS(err, "error while getting current forecast data")
		return nil, err
	}

	//klog.Infof("current forecast data for %s region is: \n", region)
	//for i := range forecast {
	//	klog.Infof("%d. Location: %s {DataStartAt: %s, DataEndAt: %s, ForecastData: %v, OptimalDataPoints: %v}\n",
	//		i+1, forecast[i].Location, forecast[i].DataStartAt.String(), forecast[i].DataEndAt.String(), forecast[i].ForecastData, forecast[i].OptimalDataPoints)
	//}
	return forecast, nil
}
