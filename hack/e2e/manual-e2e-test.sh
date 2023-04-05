#!/usr/bin/env bash

set -e -u -x

if ! type helm > /dev/null; then
  exit 1
fi

if ! type kind > /dev/null; then
  exit 1
fi


: "${SERVER_USERNAME:=${SERVER_USERNAME}}"
: "${SERVER_PASSWORD:=${SERVER_PASSWORD}}"
: "${REGISTRY:=${REGISTRY}}"
: "${SERVER_IMG_TAG:=${SERVER_IMG_TAG}}"
: "${EXPORTER_IMG_TAG:=${EXPORTER_IMG_TAG}}"


error() {
    echo "$@" >&2
    exit 1
}

TMPDIR=""

cleanup() {
 kind delete cluster -n carbon-e2e
 if [ -n "$TMPDIR" ]; then
     rm -rf "$TMPDIR"
 fi
}
trap 'cleanup' EXIT

TMPDIR="$(mktemp -d)"

kind create cluster -n carbon-e2e
make docker-build-server-image docker-build-exporter-image
kind load docker-image -n carbon-e2e $REGISTRY/server:$SERVER_IMG_TAG
kind load docker-image -n carbon-e2e $REGISTRY/exporter:$EXPORTER_IMG_TAG

helm install carbon-e2e \
   --set apiServer.image.repository=$REGISTRY/server \
   --set carbonDataExporter.image.repository=$REGISTRY/exporter \
   --set carbonDataExporter.patrolInterval=15s \
   --set apiServer.username=$SERVER_USERNAME \
   --set apiServer.password=$SERVER_PASSWORD \
   ./charts/carbon-intensity-exporter

kubectl wait --for=condition=available deploy carbon-e2e-carbon-intensity-exporter -n kube-system --timeout=300s

sleep 15

kubectl get configmap carbon-intensity -n kube-system
kubectl describe configmap carbon-intensity -n kube-system

sleep 10

kubectl delete configmap carbon-intensity -n kube-system

sleep 15

kubectl get configmap carbon-intensity -n kube-system
kubectl describe configmap carbon-intensity -n kube-system

sleep 5

kubectl logs $(echo $(kubectl get pods -n kube-system -o=name | grep carbon-e2e-carbon-intensity-exporter | sed "s/^.\{4\}//")) -c carbon-data-exporter -n kube-system
