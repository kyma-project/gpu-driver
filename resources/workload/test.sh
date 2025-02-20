#!/usr/bin/env bash

kubectl apply -f k8s.yaml

kubectl wait --for=jsonpath='{.status.phase}'=Succeeded pod/busybox1