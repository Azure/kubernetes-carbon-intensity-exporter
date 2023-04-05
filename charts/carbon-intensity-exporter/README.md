# Kubernetes Carbon Intensity Exporter Helm Chart

## Installation

Quick start instructions for the setup  using Helm.

### Prerequisites

- [Helm](https://helm.sh/docs/intro/quickstart/#install-helm)
- [AKS](https://docs.microsoft.com/en-us/azure/aks/learn/quick-kubernetes-deploy-cli)

### Installing the chart

1. Clone project

```shell

$ git https://github.com/Azure/kubernetes-carbon-intensity-exporter.git
$ cd kubernetes-carbon-intensity-exporter

```

2. Install chart using Helm v3.0+

```shell
$ export CHART_NAME=carbon-intensity-exporter

$ helm install "$CHART_NAME" ./charts/carbon-intensity-exporter
```

3. Verify that carbon intensity exporter pod is running properly

```shell
$ kubectl get pods -n kube-system | grep carbon-intensity-exporter 
```
<details>
<summary>Result</summary>

```shell
NAME                                                 READY   STATUS    RESTARTS      AGE       VERSION
expoter-carbon-intensity-exporter-766885c789-6pvfv   2/2     Running   0             15m

```
</details><br/>

### Configuration

The following table lists the configurable parameters of the kubernetes carbon intensity exporter chart and the default values.

| Parameter                                      | Description                                                                                                           | Default                               |
|------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------|---------------------------------------|
