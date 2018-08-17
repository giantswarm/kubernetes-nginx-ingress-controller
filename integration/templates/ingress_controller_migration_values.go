// +build k8srequired

package templates

// NginxIngressControllerMigrationValues values required by kubernetes-nginx-ingress-controller-chart.
const NginxIngressControllerMigrationValues = `namespace: kube-system

controller:
  name: nginx-ingress-controller
  k8sAppLabel: nginx-ingress-controller
  metricsPort: 10254

  replicas: 1

  configmap:
    name: ingress-nginx

  image:
    repository: giantswarm/nginx-ingress-controller
    tag: 0.12.0

  service:
    enabled: false
    nodePorts:
      http: 30010
      https: 30011

  resources:
    limits:
      cpu: 500m
      memory: 350Mi
    requests:
      cpu: 500m
      memory: 350Mi

defaultBackend:
  name: default-http-backend
  k8sAppLabel: default-http-backend
  port: 8080

  replicas: 1

  image:
    repository: giantswarm/defaultbackend
    tag: 1.2

  resources:
    limits:
      cpu: 10m
      memory: 20Mi
    requests:
      cpu: 10m
      memory: 20Mi

global:
  controller:
    replicas: 1
    useProxyProtocol: true
  migration:
    enabled: true

image:
  registry: quay.io

initContainer:
  image:
    registry: quay.io
    repository: giantswarm/alpine
    tag: 3.7

test:
  image:
    registry: quay.io
    repository: giantswarm/alpine-testing
    tag: 0.1.0
`
