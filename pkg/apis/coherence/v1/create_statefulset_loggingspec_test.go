/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"fmt"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateStatefulSetFromRoleWithEmptyLoggingSpec(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Logging: &coh.LoggingSpec{},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLoggingSpecWithConfigFile(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Logging: &coh.LoggingSpec{
			ConfigFile: stringPtr("/conf/test-logging.config"),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: "/conf/test-logging.config"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLoggingSpecWithConfigMapName(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Logging: &coh.LoggingSpec{
			ConfigMapName: stringPtr("test-logging-configmap"),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	mount := corev1.VolumeMount{
		Name:      coh.VolumeNameLoggingConfig,
		MountPath: coh.VolumeMountPathLoggingConfig,
	}
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, mount)

	vol := corev1.Volume{
		Name: coh.VolumeNameLoggingConfig,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test-logging-configmap",
				},
				DefaultMode: int32Ptr(0777),
				Optional:    nil,
			},
		},
	}
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, vol)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLoggingSpecWithConfigMapNameAndConfigFile(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Logging: &coh.LoggingSpec{
			ConfigMapName: stringPtr("test-logging-configmap"),
			ConfigFile:    stringPtr("test-logging.config"),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	mount := corev1.VolumeMount{
		Name:      coh.VolumeNameLoggingConfig,
		MountPath: coh.VolumeMountPathLoggingConfig,
	}
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, mount)

	vol := corev1.Volume{
		Name: coh.VolumeNameLoggingConfig,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test-logging-configmap",
				},
				DefaultMode: int32Ptr(0777),
				Optional:    nil,
			},
		},
	}
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, vol)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: coh.VolumeMountPathLoggingConfig + "/test-logging.config"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLoggingSpecWithEmptyFluentdSpec(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Logging: &coh.LoggingSpec{
			Fluentd: &coh.FluentdSpec{},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLoggingSpecWithFluentdSpecEnabledFalse(t *testing.T) {
	policy := corev1.PullAlways
	role := coh.CoherenceRoleSpec{
		Logging: &coh.LoggingSpec{
			Fluentd: &coh.FluentdSpec{
				Enabled:    boolPtr(false),
				ConfigFile: stringPtr("test-fluentd-config.cfg"),
				Tag:        stringPtr("fluentd-tag"),
				ImageSpec: coh.ImageSpec{
					Image:           stringPtr("test-fluentd:1.0"),
					ImagePullPolicy: &policy,
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLoggingSpecWithFluentdSpecEnabledTrue(t *testing.T) {
	role := coh.CoherenceRoleSpec{
		Logging: &coh.LoggingSpec{
			Fluentd: &coh.FluentdSpec{
				Enabled: boolPtr(true),
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	// Add the Fluentd container to the expected StatefulSet
	c := createExpectedFluentdContainer()
	stsExpected.Spec.Template.Spec.Containers = append(stsExpected.Spec.Template.Spec.Containers, c)
	// Add the expected Fluentd ConfigMap volume
	cmName := fmt.Sprintf(coh.EfkConfigMapNameTemplate, role.GetFullRoleName(cluster))
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: coh.VolumeNameFluentdConfig,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: cmName},
				DefaultMode:          int32Ptr(420),
			},
		},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLoggingSpecWithFluentdSpecWithImageName(t *testing.T) {
	role := coh.CoherenceRoleSpec{
		Logging: &coh.LoggingSpec{
			Fluentd: &coh.FluentdSpec{
				Enabled: boolPtr(true),
				ImageSpec: coh.ImageSpec{
					Image: stringPtr("test-fluentd:1.0"),
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	// Add the Fluentd container to the expected StatefulSet
	c := createExpectedFluentdContainer()
	c.Image = "test-fluentd:1.0"
	stsExpected.Spec.Template.Spec.Containers = append(stsExpected.Spec.Template.Spec.Containers, c)
	// Add the expected Fluentd ConfigMap volume
	cmName := fmt.Sprintf(coh.EfkConfigMapNameTemplate, role.GetFullRoleName(cluster))
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: coh.VolumeNameFluentdConfig,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: cmName},
				DefaultMode:          int32Ptr(420),
			},
		},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLoggingSpecWithFluentdSpecWithImagePullPolicy(t *testing.T) {
	policy := corev1.PullAlways
	role := coh.CoherenceRoleSpec{
		Logging: &coh.LoggingSpec{
			Fluentd: &coh.FluentdSpec{
				Enabled: boolPtr(true),
				ImageSpec: coh.ImageSpec{
					ImagePullPolicy: &policy,
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	// Add the Fluentd container to the expected StatefulSet
	c := createExpectedFluentdContainer()
	c.ImagePullPolicy = policy
	stsExpected.Spec.Template.Spec.Containers = append(stsExpected.Spec.Template.Spec.Containers, c)
	// Add the expected Fluentd ConfigMap volume
	cmName := fmt.Sprintf(coh.EfkConfigMapNameTemplate, role.GetFullRoleName(cluster))
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: coh.VolumeNameFluentdConfig,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: cmName},
				DefaultMode:          int32Ptr(420),
			},
		},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func createExpectedFluentdContainer() corev1.Container {
	pullPolicy := corev1.PullIfNotPresent

	return corev1.Container{
		Name:            coh.ContainerNameFluentd,
		Image:           coh.DefaultFluentdImage,
		ImagePullPolicy: pullPolicy,
		Args:            []string{"-c", "/etc/fluent.conf"},
		Env: []corev1.EnvVar{
			{
				Name: "COHERENCE_POD_ID",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.uid",
					},
				},
			},
			{
				Name:  "FLUENTD_CONF",
				Value: "fluentd-coherence.conf",
			},
			{
				Name:  "FLUENT_ELASTICSEARCH_SED_DISABLE",
				Value: "true",
			},
			{
				Name: "ELASTICSEARCH_HOST",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: coh.CoherenceMonitoringConfigName,
						},
						Key: coh.LoggingConfigKeyElasticSearchHost,
					},
				},
			},
			{
				Name: "ELASTICSEARCH_PORT",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: coh.CoherenceMonitoringConfigName,
						},
						Key: coh.LoggingConfigElasticSearchPort,
					},
				},
			},
			{
				Name: "ELASTICSEARCH_USER",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: coh.CoherenceMonitoringConfigName,
						},
						Key: coh.LoggingConfigElasticSearchUser,
					},
				},
			},
			{
				Name: "ELASTICSEARCH_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: coh.CoherenceMonitoringConfigName,
						},
						Key: coh.LoggingConfigElasticSearchPassword,
					},
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      coh.VolumeNameFluentdConfig,
				MountPath: coh.VolumeMountPathFluentdConfig,
				SubPath:   coh.VolumeMountSubPathFluentdConfig,
			},
			{
				Name:      coh.VolumeNameLogs,
				MountPath: coh.VolumeMountPathLogs,
			},
		},
	}
}
