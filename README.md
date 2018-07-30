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

Provide a custom `values.yaml`:

```bash
$ helm install kubernetes-nginx-ingress-controller-chart -f values.yaml
```

Deployment to Guest Clusters is handled by [chart-operator](https://github.com/giantswarm/chart-operator).
