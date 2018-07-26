// +build k8srequired

package basic

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/release"
	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/templates"
)

const (
	resourceNamespace = "kube-system"
)

func TestHelm(t *testing.T) {
	channel := os.Getenv("CIRCLE_SHA1")
	releaseName := "kubernetes-nginx-ingress-controller"

	/*
		gsHelmClient, err := createGsHelmClient()
		if err != nil {
			t.Fatalf("could not create giantswarm helmClient %v", err)
		}
	*/

	err := r.InstallResource(releaseName, templates.NginxIngressControllerValues, channel)
	if err != nil {
		t.Fatalf("could not install %q %v", releaseName, err)
	}

	err = release.WaitForStatus(helmClient, releaseName, "DEPLOYED")
	if err != nil {
		t.Fatalf("could not get release status of %q %v", releaseName, err)
	}
	l.Log("level", "debug", "message", fmt.Sprintf("%s succesfully deployed", releaseName))

	err = checkDeployment("nginx-ingress-controller", 3)
	if err != nil {
		t.Fatalf("controller manifest is incorrect: %v", err)
	}

	err = checkDeployment("default-http-backend", 2)
	if err != nil {
		t.Fatalf("default backend manifest is incorrect: %v", err)
	}

	err = framework.HelmCmd(fmt.Sprintf("test --debug %s", releaseName))
	if err != nil {
		t.Errorf("unexpected error during test of the chart: %v", err)
	}
}

// checkDeployment ensures that key properties of the deployment are correct.
func checkDeployment(name string, replicas int) error {
	expectedMatchLabels := map[string]string{
		"k8s-app": name,
	}
	expectedLabels := map[string]string{
		"k8s-app":                    name,
		"giantswarm.io/service-type": "managed",
	}

	c := f.K8sClient()
	ds, err := c.Apps().Deployments(resourceNamespace).Get(name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return microerror.Newf("could not find deployment: '%s' %v", name, err)
	} else if err != nil {
		return microerror.Newf("unexpected error getting deployment: %v", err)
	}

	// Check deployment labels.
	if !reflect.DeepEqual(expectedLabels, ds.ObjectMeta.Labels) {
		return microerror.Newf("expected labels: %v got: %v", expectedLabels, ds.ObjectMeta.Labels)
	}

	// Check selector match labels.
	if !reflect.DeepEqual(expectedMatchLabels, ds.Spec.Selector.MatchLabels) {
		return microerror.Newf("expected match labels: %v got: %v", expectedMatchLabels, ds.Spec.Selector.MatchLabels)
	}

	// Check pod labels.
	if !reflect.DeepEqual(expectedLabels, ds.Spec.Template.ObjectMeta.Labels) {
		return microerror.Newf("expected pod labels: %v got: %v", expectedLabels, ds.Spec.Template.ObjectMeta.Labels)
	}

	// Check replica count.
	if *ds.Spec.Replicas != int32(replicas) {
		return microerror.Newf("expected replicas: %d got: %d", replicas, ds.Spec.Replicas)
	}

	return nil
}

func createGsHelmClient() (*helmclient.Client, error) {
	l, err := micrologger.New(micrologger.Config{})
	if err != nil {
		return nil, microerror.Maskf(err, "could not create logger")
	}

	c := helmclient.Config{
		Logger:          l,
		K8sClient:       f.K8sClient(),
		RestConfig:      f.RestConfig(),
		TillerNamespace: "giantswarm",
	}

	gsHelmClient, err := helmclient.New(c)
	if err != nil {
		return nil, microerror.Maskf(err, "could not create helmClient")
	}

	return gsHelmClient, nil
}
