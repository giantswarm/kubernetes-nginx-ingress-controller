[![CircleCI](https://circleci.com/gh/giantswarm/kubernetes-nginx-ingress-controller.svg?style=svg&circle-token=beb2d4248f0f227ce6618f71b2c35e841b903c04)](https://circleci.com/gh/giantswarm/kubernetes-nginx-ingress-controller)

# kubernetes-nginx-ingress-controller
Helm chart for nginx ingress controller running in Guest Clusters


* Installs the [nginx-ingress-controller](https://github.com/nginxinc/kubernetes-ingress)

## Installing the Chart

To install the chart locally:

```bash
$ git clone https://github.com/giantswarm/kubernetes-nginx-ingress-controller.git
$ cd kubernetes-nginx-ingress-controller
$ helm install kubernetes-nginx-ingress-controller/helm/kubernetes-nginx-ingress-controller-chart
```

Provide a custom `Charts.yaml`:

```bash
$ helm install kubernetes-nginx-ingress-controller-chart -f values.yaml
```

Deployment to Guest Clusters will be handled by [chart-operator](https://github.com/giantswarm/chart-operator)

## Configuration

| Parameter                            | Description                                             | Default                                       |
|--------------------------------------|---------------------------------------------------------|-----------------------------------------------|
| `controller.image.repository`        | The controller container image repository               | `quay.io/giantswarm/nginx-ingress-controller` |
| `controller.image.tag`               | The controller container image tag                      | `0.11.0`                                      |
| `controller.replicaCount`            | The desired number of controller pods                   | `3`                                           |
| `controller.resources`               | The controller pod resource requests & limits           | `cpu:500m memory:350Mi`                       |
| `controller.metricsPort`             | Sets the metricsport used for metrics and health checks | `10254`                                       |
| `controller.service.nodePorts.http`  | Sets the nodePort that maps to the Ingress' port 80     | `30010`                                       |
| `controller.service.nodePorts.https` | Sets the nodePort that maps to the Ingress' port 443    | `30011`                                       |
| `defaultBackend.name`                | The name of the default backend component               | `default-http-backend`                        |
| `defaultBackend.image.repository`    | The default backend container image repository          | `quay.io/giantswarm/defaultbackend`           |
| `defaultBackend.image.tag`           | The default backend container image tag                 | `1.2`                                         |
| `defaultBackend.replicaCount`        | The desired number of default backend pods              | `2`                                           |
| `defaultBackend.resources`           | The default backend pod resource requests & limits      | `cpu:10m memory:20Mi`                         |
