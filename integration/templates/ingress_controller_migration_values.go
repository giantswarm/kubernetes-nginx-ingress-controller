// +build k8srequired

package templates

// NginxIngressControllerMigrationValues values required by kubernetes-nginx-ingress-controller-chart.
const NginxIngressControllerMigrationValues = `namespace: kube-system

controller:
  name: nginx-ingress-controller-2
  k8sAppLabel: nginx-ingress-controller
  metricsPort: 10254

  replicas: 3

  configmap:
    name: ingress-nginx-2

  image:
    registry: quay.io
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
  name: default-http-backend-2
  k8sAppLabel: default-http-backend
  port: 8080

  replicas: 2

  image:
    registry: quay.io
    repository: giantswarm/defaultbackend
    tag: 1.2

  resources:
    limits:
      cpu: 10m
      memory: 20Mi
    requests:
      cpu: 10m
      memory: 20Mi

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
