// +build k8srequired

package templates

// NginxIngressControllerBasicValues values required by kubernetes-nginx-ingress-controller-chart.
const NginxIngressControllerBasicValues = `namespace: kube-system

controller:
  name: nginx-ingress-controller
  k8sAppLabel: nginx-ingress-controller
  metricsPort: 10254

  replicas: 1

  configmap:
    name: ingress-nginx

  image:
    repository: giantswarm/nginx-ingress-controller
    tag: 0.23.0

  service:
    enabled: true
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

  replicas: 2

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
    useProxyProtocol: false
  migration:
    enabled: false

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
