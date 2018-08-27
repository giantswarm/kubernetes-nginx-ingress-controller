// +build k8srequired

package basic

import (
	"context"
	"fmt"
	"testing"

	"github.com/giantswarm/e2etests/managedservices"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/env"
	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/templates"
)

const (
	chartName          = "kubernetes-nginx-ingress-controller"
	controllerName     = "nginx-ingress-controller"
	defaultBackendName = "default-http-backend"
	testName           = "basic"
)

func TestHelm(t *testing.T) {
	chartConfig := managedservices.ChartConfig{
		ChannelName: fmt.Sprintf("%s-%s", env.CircleSHA(), testName),
		ChartName:   chartName,
		ChartValues: templates.NginxIngressControllerBasicValues,
		Namespace:   metav1.NamespaceSystem,
	}
	chartResources := managedservices.ChartResources{
		Deployments: []managedservices.Deployment{
			{
				Name:      controllerName,
				Namespace: metav1.NamespaceSystem,
				Labels: map[string]string{
					"app": controllerName,
					"giantswarm.io/service-type": "managed",
					"k8s-app":                    controllerName,
				},
				MatchLabels: map[string]string{
					"k8s-app": controllerName,
				},
				Replicas: 3,
			},
			{
				Name:      defaultBackendName,
				Namespace: metav1.NamespaceSystem,
				Labels: map[string]string{
					"app": defaultBackendName,
					"giantswarm.io/service-type": "managed",
					"k8s-app":                    defaultBackendName,
				},
				MatchLabels: map[string]string{
					"k8s-app": defaultBackendName,
				},
				Replicas: 2,
			},
		},
	}

	err := ms.Test(context.Background(), chartConfig, chartResources)
	if err != nil {
		t.Fatalf("%#v", err)
	}
}
