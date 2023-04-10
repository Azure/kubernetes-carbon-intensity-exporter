/*
MIT License
Copyright (c) Microsoft Corporation.
*/
package exporter

import (
	"fmt"
	"time"

	"github.com/Azure/kubernetes-carbon-intensity-exporter/pkg/sdk/client"
	"github.com/antihax/optional"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
)

var (
	constantBackoff = wait.Backoff{
		Duration: 3 * time.Second,
		Steps:    10,
	}
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
	err := e.RefreshData(ctx, configMapName, region)
	if err != nil {
		return
	}

	informerFactory := informers.NewSharedInformerFactory(e.clusterClient, time.Hour*1)
	configMapInformer := informerFactory.Core().V1().ConfigMaps().Informer()

	// Create a channel to receive events from the informer
	eventChan := make(chan interface{})
	defer close(eventChan)

	configMapInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: func(obj interface{}) {
				cm, ok := obj.(*corev1.ConfigMap)
				if !ok {
					return
				}
				if cm.ObjectMeta.Name != configMapName {
					return
				}
				eventChan <- obj
			},
		})

	go configMapInformer.Run(stopChan)
	// Wait for the informer to sync
	if !cache.WaitForCacheSync(stopChan, configMapInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for cache to sync"))
		return
	}

	refreshPatrol := time.NewTicker(patrolInterval)
	defer refreshPatrol.Stop()

	for {
		select {
		case <-eventChan:
			err := e.RefreshData(ctx, configMapName, region)
			if err != nil {
				return
			}

			e.recorder.Eventf(&corev1.ObjectReference{
				Kind:      "Pod",
				Namespace: client.Namespace,
				Name:      client.PodName,
			}, corev1.EventTypeWarning, "Configmap Deleted", "Configmap got deleted")

		// if refresh time elapsed
		case <-refreshPatrol.C:
			err := e.DeleteConfigMap(ctx, configMapName)
			if err != nil && !apierrors.IsNotFound(err) {
				break
			}
			e.recorder.Eventf(&corev1.ObjectReference{
				Kind:      "Pod",
				Namespace: client.Namespace,
				Name:      client.PodName,
			}, corev1.EventTypeNormal, "Configmap updated", "Configmap got updated")

			// context got canceled or done
		case <-ctx.Done():
		case <-stopChan:
			return
		}
	}
}

func (e *Exporter) RefreshData(ctx context.Context, configMapName string, region string) error {
	// get current object (if any) in case we could not update the data.
	currentConfigMap, err := e.GetConfigMap(ctx, configMapName)
	if err != nil {
		return err
	}

	//delete old configmap if any
	err = e.DeleteConfigMap(ctx, configMapName)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	var forecast []client.EmissionsForecastDto
	err = retry.OnError(constantBackoff, func(err error) bool {
		return true
	}, func() error {
		forecast, err = e.getCurrentForecastData(ctx, region)
		return err
	})
	if err != nil {
		if currentConfigMap != nil {
			// return old data with failed message
			return e.UseCurrentConfigMap(ctx, err.Error(), currentConfigMap)
		} else {
			e.recorder.Eventf(&corev1.ObjectReference{
				Kind:      "Pod",
				Namespace: client.Namespace,
				Name:      client.PodName,
			}, corev1.EventTypeWarning, "Cannot retrieve updated forecast data", "Error while retrieving updated forecast data")
			klog.Errorf("an error has occurred while retrieving updated forecast data")
			return err
		}
	}

	err = retry.OnError(constantBackoff, func(err error) bool {
		return true
	}, func() error {
		return e.CreateConfigMapFromEmissionForecast(ctx, configMapName, forecast)
	})
	if err != nil {
		e.recorder.Eventf(&corev1.ObjectReference{
			Kind:      "Pod",
			Namespace: client.Namespace,
			Name:      client.PodName,
		}, corev1.EventTypeWarning, "Configmap Create", "Error while creating configMap")
		klog.Errorf("an error has occurred while creating %s configMap, err: %s", configMapName, err.Error())
		return err
	}
	e.recorder.Eventf(&corev1.ObjectReference{
		Kind:      "Pod",
		Namespace: client.Namespace,
		Name:      client.PodName,
	}, corev1.EventTypeNormal, "Exporter results", "Done retrieve data")
	return nil
}

func (e *Exporter) UseCurrentConfigMap(ctx context.Context, message string, currentConfigMap *corev1.ConfigMap) error {
	if currentConfigMap.Data != nil {
		currentConfigMap.Data[ConfigMapLastHeartbeatTime] = time.Now().String()
		currentConfigMap.Data[ConfigMapMessage] = message
	} else {
		currentConfigMap.Data = map[string]string{
			ConfigMapLastHeartbeatTime: time.Now().String(),
			ConfigMapMessage:           message,
		}
	}
	if currentConfigMap.BinaryData == nil {
		currentConfigMap.BinaryData = map[string][]byte{
			BinaryData: {},
		}
	}
	return e.CreateConfigMapFromProperties(ctx, currentConfigMap.Name,
		currentConfigMap.Data, currentConfigMap.BinaryData[BinaryData])
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
