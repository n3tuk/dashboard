# `dashboard` Helm Chart

![Type: application](https://img.shields.io/badge/Type-application-informational?style=for-the-badge)

This is the documentation for the `dashboard` Helm Chart,
providing information on how to use and deploy the application, and how to
configure it through the `values.yaml` file.

## Chart Maintainers

<!-- prettier-ignore-start -->

| Name | Email | Url |
| ---- | ------ | --- |
| Jonathan Wright | <jon@than.io> | <https://github.com/jonathanio> |

<!-- prettier-ignore-end -->

## Chart Values

<!-- prettier-ignore-start -->

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| nameOverride | string | `nil` | Override the name of the objects created by this chart |
| fullnameOverride | string | `nil` | Override the name of the objects created by this chart |
| podDisruptionBudget.create | bool | `true` | Set whether or not to create the podDisruptionBudget resource for the dashboard deployment |
| podDisruptionBudget.minAvailable | int | `1` | The number of pods in a RepicaSet which much be available after an eviction, either absolute or a percentage. Only one of `minAvailable` or `maxAvailable` can be configured, with `minAvailable` taking priority |
| podDisruptionBudget.maxUnavilable | int | `0` | The number of pods in a RepicaSet which can be unavailable after an eviction, either absolute or a percentage |
| podDisruptionBudget.annotations | object | `{}` | Set any additional annotations which should be added to the PodDisruptionBudget resource |
| podDisruptionBudget.labels | object | `{}` | Set any additional labels which should be added to the PodDisruptionBudget resource |
| service.type | string | `"ClusterIP"` | Set whether the Service should be a ClusterIP or NodeIP |
| service.port | int | `8080` | Set the TCP port the Service should be configured to listen on |
| service.annotations | object | `{}` | Set any additional annotations which should be added to the Ingress resource |
| service.labels | object | `{}` | Set any additional labels which should be added to the Ingress resource |
| ingress.create | bool | `false` | Set whether or not to create the Ingress for the dashboard Service |
| ingress.className | string | `nil` | Set the class name of the controller which will handle the Ingress |
| ingress.annotations | object | `{}` | Set any additional annotations which should be added to the Ingress resource |
| ingress.labels | object | `{}` | Set any additional labels which should be added to the Ingress resource |
| ingress.hosts | list | `[]` | Set the hostname and path mappings for this service on the Ingress |
| ingress.tls | list | `[]` | Set the TLS secret and hostnames for this service on the Ingress |
| serviceAccount.create | bool | `false` | Set whether or not to create a ServiceAccount resource for the dashboard service |
| serviceAccount.name | string | `nil` | Override the name of the ServiceAccount @default `.Chart.Name` |
| serviceAccount.annotations | object | `{}` | Set any additional annotations which should be added to the ServiceAccount resource |
| serviceAccount.labels | object | `{}` | Set any additional labels which should be added to the ServiceAccount resource |
| networkPolicy.create | bool | `false` | Set whether or not to create a NetworkPolicy resource for the dashboard service |
| networkPolicy.annotations | object | `{}` | Set any additional annotations which should be added to the NetworkPolicy resource |
| networkPolicy.labels | object | `{}` | Set any additional labels which should be added to the NetworkPolicy resource |
| deployment.replicaCount | int | `1` | Set the number of replicas of this deployment which should be deployed onto the kubernetes cluster |
| deployment.revisionHistoryLimit | int | `10` | Set the number of deployments which should be kept to enable a rollback of the deployment in the event of any issues or failures |
| deployment.annotations | object | `{}` | Set any additional annotations which should be added to the Deployment resource |
| deployment.labels | object | `{}` | Set any additional labels which should be added to the Deployment resource |
| persistentVolumeClaim.create | bool | `false` | Set whether or not to create a PersistentVolumeClaim resource for the dashboard service and attach it to the Pods |
| persistentVolumeClaim.storageClassName | string | `nil` | Set the name of the StorageClass to use for the volumes in the PersistentVolumeClaim |
| persistentVolumeClaim.size | string | `"32Gi"` | Set the size of each PersistentVolumeClaim to be created |
| persistentVolumeClaim.accessModes | list | `["ReadWriteOnce"]` | Configure the access modes to be set on the PersistentVolumeClaim |
| pod.image.repository | string | `"ghcr.io/n3tuk/dashboard"` | Set the URI for the container image to be deployed for the dashboard Deployment |
| pod.image.pullPolicy | string | `"IfNotPresent"` | Set the pull policy for the host running each Pod of the deployment |
| pod.annotations | object | `{}` | Set any additional annotations which should be added to the Ingress resource |
| pod.labels | object | `{}` | Set any additional labels which should be added to the Ingress resource |
| pod.nodeSelector | object | `{}` | Set any node labels which should be used for assigning Pods to Nodes |
| pod.tolerations | list | `[]` | Set the list of node taints to tolerate for deployment, if required |
| pod.affinity | object | `{}` | Set the Node/Pod affinities used to schedule and distribute the deployed Pods |
| pod.topologySpreadConstraints | list | `[]` | Set the constraints which should be applied to help schedule and run the Pods across the cluster |
| pod.resources | object | `{"requests":{"cpu":"50m","memory":"64Mi"}}` | Set the resources which should be requested and limited to when the Pod is run on each host |
| pod.extraArgs | list | `[]` | Set addtional command-line arguments to pass to the dashboard application at runtime |
| pod.logging.level | string | `"info"` | Set the logging level for the dashboard application, selecting from debug, info, warning, and error |
| pod.logging.json | bool | `true` | Set whether or not to output the structured logs in JSON forward, or in plain text format |
| pod.probes.startup.create | bool | `true` | Set whether or not to add the Startup Probe to the Pod |
| pod.probes.startup.periodSeconds | int | `2` | Set how often (in seconds) to trigger the probe |
| pod.probes.startup.initialDelaySeconds | int | `0` | Number of seconds after the container has started before the probe is initiated |
| pod.probes.startup.timeoutSeconds | int | `1` | Number of seconds after which the probe times out |
| pod.probes.startup.successThreshold | int | `1` | Minimum consecutive successes for the probe to be considered successful after having failed |
| pod.probes.startup.failureThreshold | int | `30` | Number of times in a row the probe should fail before  Kubernetes considers that the overall check has failed |
| pod.probes.liveness.create | bool | `true` | Set whether or not to add the Liveness Probe to the Pod |
| pod.probes.liveness.periodSeconds | int | `5` | Set how often (in seconds) to trigger the probe |
| pod.probes.liveness.initialDelaySeconds | int | `0` | Number of seconds after the container has started before the probe is initiated |
| pod.probes.liveness.timeoutSeconds | int | `3` | Number of seconds after which the probe times out |
| pod.probes.liveness.successThreshold | int | `1` | Minimum consecutive successes for the probe to be considered successful after having failed |
| pod.probes.liveness.failureThreshold | int | `2` | Number of times in a row the probe should fail before  Kubernetes considers that the overall check has failed |
| pod.probes.readiness.create | bool | `true` | Set whether or not to add the Readiness Probe to the Pod |
| pod.probes.readiness.periodSeconds | int | `15` | Set how often (in seconds) to trigger the probe |
| pod.probes.readiness.initialDelaySeconds | int | `0` | Number of seconds after the container has started before the probe is initiated |
| pod.probes.readiness.timeoutSeconds | int | `5` | Number of seconds after which the probe times out |
| pod.probes.readiness.successThreshold | int | `1` | Minimum consecutive successes for the probe to be considered successful after having failed |
| pod.probes.readiness.failureThreshold | int | `1` | Number of times in a row the probe should fail before  Kubernetes considers that the overall check has failed |
| serviceMonitor.create | bool | `false` | Set whether or not to create a ServiceMonitor resource for the dashboard service |
| serviceMonitor.interval | string | `"15s"` | Set the number of seconds Prometheus should monitor the metrics for on this service |

<!-- prettier-ignore-end -->
