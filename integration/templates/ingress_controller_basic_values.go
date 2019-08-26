// +build k8srequired

package templates

// NginxIngressControllerBasicValues values required by kubernetes-nginx-ingress-controller-chart.
const NginxIngressControllerBasicValues = `namespace: kube-system

controller:
  name: nginx-ingress-controller
  k8sAppLabel: nginx-ingress-controller
  metricsPort: 10254

  replicas: 3
  maxUnavailable: 0

  configmap:
    name: ingress-nginx

  image:
    repository: giantswarm/nginx-ingress-controller
    tag: 0.25.1

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
    tag: 3.9-giantswarm

test:
  image:
    registry: quay.io
    repository: giantswarm/alpine-testing
    tag: 0.1.0
`
