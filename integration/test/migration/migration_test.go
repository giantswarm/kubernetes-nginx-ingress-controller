// +build k8srequired

package migration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/env"
	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/templates"
)

const (
	resourceNamespace = "kube-system"
	testName          = "migration"
)

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

	// Check controller service is present.
	err = checkControllerServicePresent()
	if err != nil {
		t.Fatalf("controller service present: %v", err)
	}
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

	channel := fmt.Sprintf("%s-%s", env.CircleSHA(), testName)
	releaseName := "kubernetes-nginx-ingress-controller"
	err = r.InstallResource(releaseName, templates.NginxIngressControllerMigrationValues, channel)
	if err != nil {
		t.Fatalf("could not install %q %v", releaseName, err)
	}

	err = r.WaitForStatus(releaseName, "DEPLOYED")
	if err != nil {
		t.Fatalf("could not get release status of %q %v", releaseName, err)
	}
	l.Log("level", "debug", "message", fmt.Sprintf("%s succesfully deployed", releaseName))

	defer framework.HelmCmd(fmt.Sprintf("delete %s --purge", releaseName))

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
	// Check controller service is still present.
	err = checkControllerServicePresent()
	if err != nil {
		t.Fatalf("controller service present: %v", err)
	}
}

func checkControllerServicePresent() error {
	c := h.K8sClient()

	controllerListOptions := metav1.ListOptions{
		LabelSelector: "k8s-app=nginx-ingress-controller,kind=legacy",
	}
	s, err := c.Core().Services(resourceNamespace).List(controllerListOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(s.Items) != 1 {
		return microerror.Newf("unexpected number of services, want 1, got %d", len(s.Items))
	}

	return nil
}

func checkResourcesPresent(labelSelector string) error {
	c := h.K8sClient()
	controllerListOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("k8s-app=nginx-ingress-controller,%s", labelSelector),
	}
	backendListOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("k8s-app=default-http-backend,%s", labelSelector),
	}
	configMapListOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("k8s-addon=ingress-nginx.addons.k8s.io,%s", labelSelector),
	}

	cm, err := c.Core().ConfigMaps(resourceNamespace).List(configMapListOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(cm.Items) != 1 {
		return microerror.Newf("unexpected number of configmaps, want 1, got %d", len(cm.Items))
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
	c := h.K8sClient()
	controllerListOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("k8s-app=nginx-ingress-controller,%s", labelSelector),
	}
	backendListOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("k8s-app=default-http-backend,%s", labelSelector),
	}
	configMapListOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("k8s-addon=ngress-nginx.addons.k8s.io,%s", labelSelector),
	}

	cm, err := c.Core().ConfigMaps(resourceNamespace).List(configMapListOptions)
	if err == nil && len(cm.Items) > 0 {
		return microerror.New("expected error querying for configmaps didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
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
func checkDeployment(name string, replicas int) error {
	expectedMatchLabels := map[string]string{
		"k8s-app": name,
	}
	expectedLabels := map[string]string{
		"k8s-app":                    name,
		"giantswarm.io/service-type": "managed",
	}

	c := h.K8sClient()
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
