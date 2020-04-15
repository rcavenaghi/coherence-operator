/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateStatefulSetFromRoleWithApplicationType(t *testing.T) {
	// Create a role with an ApplicationSpec with an application type
	role := coh.CoherenceRoleSpec{
		Application: &coh.ApplicationSpec{
			Type: stringPtr("foo"),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "APP_TYPE", Value: "foo"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithApplicationMain(t *testing.T) {
	// Create a role with an ApplicationSpec with a main
	mainClass := "com.tangosol.net.CacheFactory"
	role := coh.CoherenceRoleSpec{
		Application: &coh.ApplicationSpec{
			Main: stringPtr(mainClass),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_MAIN_CLASS", Value: mainClass})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithApplicationMainArgs(t *testing.T) {
	// Create a role with an ApplicationSpec with a main
	role := coh.CoherenceRoleSpec{
		Application: &coh.ApplicationSpec{
			Args: []string{"arg1", "arg2"},
			ImageSpec: coh.ImageSpec{
				Image:           nil,
				ImagePullPolicy: nil,
			},
			AppDir:    nil,
			LibDir:    nil,
			ConfigDir: nil,
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_MAIN_ARGS", Value: "arg1 arg2"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithApplicationEmptyMainArgs(t *testing.T) {
	// Create a role with an ApplicationSpec with a main
	role := coh.CoherenceRoleSpec{
		Application: &coh.ApplicationSpec{
			Args: []string{},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithApplicationImageName(t *testing.T) {
	// Create a role with an ApplicationSpec with an application image
	role := coh.CoherenceRoleSpec{
		Application: &coh.ApplicationSpec{
			ImageSpec: coh.ImageSpec{
				Image: stringPtr("app-image:1.0"),
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// Append the expected application container
	container := createDefaultApplicationContainer("app-image:1.0")
	stsExpected.Spec.Template.Spec.InitContainers = append(stsExpected.Spec.Template.Spec.InitContainers, container)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithApplicationImagePullPolicy(t *testing.T) {
	// Create a role with an ApplicationSpec with an image pull policy
	policy := corev1.PullAlways
	role := coh.CoherenceRoleSpec{
		Application: &coh.ApplicationSpec{
			ImageSpec: coh.ImageSpec{
				Image:           stringPtr("app-image:1.0"),
				ImagePullPolicy: &policy,
			},
			AppDir:    nil,
			LibDir:    nil,
			ConfigDir: nil,
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// Append the expected application container
	container := createDefaultApplicationContainer("app-image:1.0")
	container.ImagePullPolicy = policy
	stsExpected.Spec.Template.Spec.InitContainers = append(stsExpected.Spec.Template.Spec.InitContainers, container)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithApplicationDirectory(t *testing.T) {
	// Create a role with an ApplicationSpec with an application directory
	dir := "/home/foo/app"
	role := coh.CoherenceRoleSpec{
		Application: &coh.ApplicationSpec{
			ImageSpec: coh.ImageSpec{
				Image: stringPtr("app-image:1.0"),
			},
			AppDir: &dir,
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// Append the expected application container
	container := createDefaultApplicationContainer("app-image:1.0")

	stsExpected.Spec.Template.Spec.InitContainers = append(stsExpected.Spec.Template.Spec.InitContainers, container)
	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameApplication, corev1.EnvVar{Name: "APP_DIR", Value: dir})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithApplicationLibDirectory(t *testing.T) {
	// Create a role with an ApplicationSpec with an application lib directory
	dir := "/home/foo/lib"
	role := coh.CoherenceRoleSpec{
		Application: &coh.ApplicationSpec{
			ImageSpec: coh.ImageSpec{
				Image: stringPtr("app-image:1.0"),
			},
			LibDir: &dir,
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// Append the expected application container
	container := createDefaultApplicationContainer("app-image:1.0")

	stsExpected.Spec.Template.Spec.InitContainers = append(stsExpected.Spec.Template.Spec.InitContainers, container)
	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameApplication, corev1.EnvVar{Name: "LIB_DIR", Value: dir})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithApplicationConfDirectory(t *testing.T) {
	// Create a role with an ApplicationSpec with an application lib directory
	dir := "/home/foo/lib"
	role := coh.CoherenceRoleSpec{
		Application: &coh.ApplicationSpec{
			ImageSpec: coh.ImageSpec{
				Image: stringPtr("app-image:1.0"),
			},
			ConfigDir: &dir,
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// Append the expected application container
	container := createDefaultApplicationContainer("app-image:1.0")

	stsExpected.Spec.Template.Spec.InitContainers = append(stsExpected.Spec.Template.Spec.InitContainers, container)
	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameApplication, corev1.EnvVar{Name: "CONF_DIR", Value: dir})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}
