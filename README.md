# THIS REPOSITORY IS DEPRECATED.

Use [nginx-ingress-controller-app](https://github.com/giantswarm/nginx-ingress-controller-app) for further development.

# kubernetes-nginx-ingress-controller

DEPRECATED Changes should be made to https://github.com/giantswarm/nginx-ingress-controller-app.

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
