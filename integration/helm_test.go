// +build k8srequired

package integration

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	resourceNamespace = "kube-system"
)

var (
	f *framework.Host
)

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	var v int
	var err error

	f, err = framework.NewHost(framework.HostConfig{})
	if err != nil {
		panic(err.Error())
	}

	err = f.CreateNamespace("giantswarm")
	if err != nil {
		log.Printf("%#v\n", err)
		v = 1
	}

	if v == 0 {
		v = m.Run()
	}

	if os.Getenv("KEEP_RESOURCES") != "true" {
		f.Teardown()
	}

	os.Exit(v)
}

func TestHelm(t *testing.T) {
	channel := os.Getenv("CIRCLE_SHA1")

	err := framework.HelmCmd(fmt.Sprintf("registry install --wait quay.io/giantswarm/kubernetes-nginx-ingress-controller-chart:%s -n test-deploy", channel))
	if err != nil {
		t.Errorf("unexpected error during installation of the chart: %v", err)
	}
	defer framework.HelmCmd("delete test-deploy --purge")

	controllerLabels := map[string]string{
		"app":                        "nginx-ingress-controller",
		"k8s-app":                    "nginx-ingress-controller",
		"giantswarm.io/service-type": "managed",
	}
	err = checkDeployment("nginx-ingress-controller", 3, controllerLabels)
	if err != nil {
		t.Fatalf("controller manifest is incorrect: %v", err)
	}

	defaultBackendLabels := map[string]string{
		"k8s-app":                    "default-http-backend",
		"giantswarm.io/service-type": "managed",
	}
	err = checkDeployment("default-http-backend", 2, defaultBackendLabels)
	if err != nil {
		t.Fatalf("default backend manifest is incorrect: %v", err)
	}

	err = framework.HelmCmd("test --debug test-deploy")
	if err != nil {
		t.Errorf("unexpected error during test of the chart: %v", err)
	}
}

// TestMigration ensures that previously deployed resources are properly
// removed.
// It installs a chart with the same resources as nginx-ingress-controller with
// appropriate labels so that we can query for them. Then installs the
// nginx-ingress-controller chart and checks that the previous resources are
// removed and the ones for nginx-ingress-controller are in place.
func TestMigration(t *testing.T) {
	// Install legacy resources.
	err := framework.HelmCmd("install /e2e/fixtures/resources-chart -n resources")
	if err != nil {
		t.Fatalf("could not install resources chart: %v", err)
	}
	defer framework.HelmCmd("delete resources --purge")

	// Check legacy resources are present.
	err = checkResourcesPresent("kind=legacy")
	if err != nil {
		t.Fatalf("legacy resources present: %v", err)
	}
	// Check managed resources are not present.
	err = checkResourcesNotPresent("giantswarm.io/service-type=managed")
	if err != nil {
		t.Fatalf("managed resources not present: %v", err)
	}

	// Install kubernetes-nginx-ingress-controller-chart.
	channel := os.Getenv("CIRCLE_SHA1")
	err = framework.HelmCmd(fmt.Sprintf("registry install --wait quay.io/giantswarm/kubernetes-nginx-ingress-controller-chart:%s -n test-deploy", channel))
	if err != nil {
		t.Fatalf("could not install kubernetes-nginx-ingress-controller-chart: %v", err)
	}
	defer framework.HelmCmd("delete test-deploy --purge")

	// Check legacy resources are not present.
	err = checkResourcesNotPresent("kind=legacy")
	if err != nil {
		t.Fatalf("legacy resources present: %v", err)
	}
	// Check managed resources are present.
	err = checkResourcesPresent("giantswarm.io/service-type=managed")
	if err != nil {
		t.Fatalf("managed resources not present: %v", err)
	}
}

func checkResourcesPresent(labelSelector string) error {
	c := f.K8sClient()
	controllerListOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("k8s-app=nginx-ingress-controller,%s", labelSelector),
	}
	backendListOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("k8s-app=default-http-backend,%s", labelSelector),
	}

	d, err := c.Extensions().Deployments(resourceNamespace).List(controllerListOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(d.Items) != 1 {
		return microerror.Newf("unexpected number of deployments, want 1, got %d", len(d.Items))
	}

	db, err := c.Extensions().Deployments(resourceNamespace).List(backendListOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(db.Items) != 1 {
		return microerror.Newf("unexpected number of deployments, want 1, got %d", len(db.Items))
	}

	cr, err := c.Rbac().ClusterRoles().List(controllerListOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(cr.Items) != 1 {
		return microerror.Newf("unexpected number of roles, want 1, got %d", len(cr.Items))
	}

	crb, err := c.Rbac().ClusterRoleBindings().List(controllerListOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(crb.Items) != 1 {
		return microerror.Newf("unexpected number of rolebindings, want 1, got %d", len(crb.Items))
	}

	r, err := c.Rbac().Roles(resourceNamespace).List(controllerListOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(r.Items) != 1 {
		return microerror.Newf("unexpected number of roles, want 1, got %d", len(r.Items))
	}

	rb, err := c.Rbac().RoleBindings(resourceNamespace).List(controllerListOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(rb.Items) != 1 {
		return microerror.Newf("unexpected number of rolebindings, want 1, got %d", len(rb.Items))
	}

	sb, err := c.Core().Services(resourceNamespace).List(backendListOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(sb.Items) != 1 {
		return microerror.Newf("unexpected number of services, want 1, got %d", len(sb.Items))
	}

	sa, err := c.Core().ServiceAccounts(resourceNamespace).List(controllerListOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(sa.Items) != 1 {
		return microerror.Newf("unexpected number of serviceaccountss, want 1, got %d", len(sa.Items))
	}

	return nil
}

func checkResourcesNotPresent(labelSelector string) error {
	c := f.K8sClient()
	controllerListOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("k8s-app=nginx-ingress-controller,%s", labelSelector),
	}
	backendListOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("k8s-app=default-http-backend,%s", labelSelector),
	}

	d, err := c.Extensions().Deployments(resourceNamespace).List(controllerListOptions)
	if err == nil && len(d.Items) > 0 {
		return microerror.New("expected error querying for deployments didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	db, err := c.Extensions().Deployments(resourceNamespace).List(backendListOptions)
	if err == nil && len(db.Items) > 0 {
		return microerror.New("expected error querying for deployments didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	cr, err := c.Rbac().ClusterRoles().List(controllerListOptions)
	if err == nil && len(cr.Items) > 0 {
		return microerror.New("expected error querying for roles didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	crb, err := c.Rbac().ClusterRoleBindings().List(controllerListOptions)
	if err == nil && len(crb.Items) > 0 {
		return microerror.New("expected error querying for rolebindings didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	r, err := c.Rbac().Roles(resourceNamespace).List(controllerListOptions)
	if err == nil && len(r.Items) > 0 {
		return microerror.New("expected error querying for roles didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	rb, err := c.Rbac().RoleBindings(resourceNamespace).List(controllerListOptions)
	if err == nil && len(rb.Items) > 0 {
		return microerror.New("expected error querying for rolebindings didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	sb, err := c.Core().Services(resourceNamespace).List(backendListOptions)
	if err == nil && len(sb.Items) > 0 {
		return microerror.New("expected error querying for services didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	sa, err := c.Core().ServiceAccounts(resourceNamespace).List(controllerListOptions)
	if err == nil && len(sa.Items) > 0 {
		return microerror.New("expected error querying for serviceaccounts didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	return nil
}

// checkDeployment ensures that key properties of the deployment are correct.
func checkDeployment(name string, replicas int, expctedLabels map[string]string) error {
	expectedMatchLabels := map[string]string{
		"k8s-app": name,
	}
	expectedLabels := map[string]string{
		"app":                        name,
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
