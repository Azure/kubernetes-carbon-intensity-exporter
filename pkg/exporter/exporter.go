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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
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

func (e *Exporter) Run(ctx context.Context, configMapName, region string, patrolInterval time.Duration, stopChan <-chan struct{}) {
	// create configMap first time
	err := e.RefreshData(ctx, configMapName, region, stopChan)
	if err != nil {
		return
	}
	configMapWatch := e.GetGonfigMapWatch(ctx, configMapName)

	// configMapWatch is reassigned when the config map is deleted,
	// so we want to make sure to stop the last instance.
	defer func() {
		configMapWatch.Stop()
	}()

	refreshPatrol := time.NewTicker(patrolInterval)
	defer refreshPatrol.Stop()

	for {
		select {
		// if the configMap got deleted by user
		case event := <-configMapWatch.ResultChan():
			if event.Type == watch.Deleted {
				err := e.Patrol(ctx, configMapName, region, stopChan)
				if err != nil {
					return
				}

				configMapWatch.Stop()
				// refresh watch after deletion
				configMapWatch = e.GetGonfigMapWatch(ctx, configMapName)

				e.recorder.Eventf(&corev1.ObjectReference{
					Kind:      "Pod",
					Namespace: namespace,
					Name:      podName,
				}, corev1.EventTypeWarning, "Configmap Deleted", "Configmap got deleted")
			}

		// if refresh time elapsed
		case <-refreshPatrol.C:
			var err error
			retry.OnError(retry.DefaultBackoff, func(err error) bool {
				return true
			}, func() error {
				err = e.RefreshData(ctx, configMapName, region, stopChan)
				return err
			})
			if err != nil {
				return
			}

			configMapWatch.Stop()
			// refresh watch after deletion
			configMapWatch = e.GetGonfigMapWatch(ctx, configMapName)

			e.recorder.Eventf(&corev1.ObjectReference{
				Kind:      "Pod",
				Namespace: namespace,
				Name:      podName,
			}, corev1.EventTypeNormal, "Configmap updated", "Configmap gets updated")

			// context got canceled or done
		case <-ctx.Done():
		case <-stopChan:
			return
		}
	}
}

func (e *Exporter) RefreshData(ctx context.Context, configMapName string, region string, stopChan <-chan struct{}) error {
	err := e.DeleteConfigmap(ctx, configMapName)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	err = e.Patrol(ctx, configMapName, region, stopChan)
	if err != nil {
		return err
	}
	return nil
}

func (e *Exporter) Patrol(ctx context.Context, configMapName, region string, stopChan <-chan struct{}) error {
	forecast, err := e.getCurrentForecastData(ctx, region, stopChan)
	if err != nil {
		return err
	}

	err = e.CreateOrUpdateConfigMap(ctx, configMapName, forecast)
	if err != nil {
		e.recorder.Eventf(&corev1.ObjectReference{
			Kind:      "Pod",
			Namespace: namespace,
			Name:      podName,
		}, corev1.EventTypeWarning, "Configmap Create", "Error while creating configMap")
		klog.Errorf("an error has occurred while creating %s configMap", configMapName)
		return err
	}
	e.recorder.Eventf(&corev1.ObjectReference{
		Kind:      "Pod",
		Namespace: client.Namespace,
		Name:      client.PodName,
	}, corev1.EventTypeNormal, "Exporter results", "Done retrieve data")
	return nil
}

func (e *Exporter) getCurrentForecastData(ctx context.Context, region string, stopChan <-chan struct{}) ([]client.EmissionsForecastDto, error) {
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
