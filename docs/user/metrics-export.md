# GPU Metrics Export

Besides installing the relevant drivers on the GPU-based nodes, typically you need to observe the GPU for troubleshooting and operations. For NVIDIA based GPUs, this can be achieved by running an additional component provided by NVIDIA on every relevant node as a DaemonSet, the [DCGM-exporter](https://github.com/NVIDIA/dcgm-exporter). This component will expose metrics in a prometheus-compatible format and with that can be integrated nicely into existing setups like by leveraging the [telemetry module](https://kyma-project.io/#/telemetry-manager/user/README)

This guide gives you a heads up on how to get started.

## Instructions

1. Define your custom values.yaml
   The exporter is packaged via Helm. Available settings can be found in the [values.yaml](https://github.com/NVIDIA/dcgm-exporter/blob/main/deployment/values.yaml). To get started, create your own custom values.yaml file with following settings:
   ```yaml
   # Select the proper nodes where a GPU is available
   nodeSelector:
     worker.gardener.cloud/pool: myWorkerPool

   # Disable the ServiceMonitor relying on the Prometheus Operator
   serviceMonitor:
     enabled: false

   # Enable a Service for the metrics endpoint and place prometheus scrape annotations
   service: 
     enable: true
     annotations:
       prometheus.io/scrape: "true"
   ```

2. Install the helm chart using your local values.yaml file
   ```sh
   helm repo add gpu-helm-charts https://nvidia.github.io/dcgm-exporter/helm-charts
   helm repo update

   helm upgrade --install gpu-metric-exporter gpu-helm-charts/dcgm-exporter -f values.yaml
   ```

3. Verify the installation
   Check that pods are being running on all the relevant pods:
   ```sh
   kubectl get pods -l "app.kubernetes.io/name=dcgm-exporter" -owide
   ```

4. See the metrics coming in
   Get a monitoring setup in place based on scraping of workloads being annotated with prometheus annotations, for example using a `MetricPipeline` of the [Kyma Telemetry module](https://kyma-project.io/#/telemetry-manager/user/04-metrics?id=_4-activate-prometheus-based-metrics), maybe in combination with [Prometheus integration](https://kyma-project.io/#/telemetry-manager/user/integration/prometheus/README).
   Then go to the related UI and check that metrics with names `xxx` are being ingested.
