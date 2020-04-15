/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestCreateStatefulSetFromRoleWithPersistenceEmpty(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{},
			Snapshot:    &coh.PersistentStorageSpec{},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithPersistenceDisabledAndSnapshotDisabled(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{
				Enabled: boolPtr(false),
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
					StorageClassName: stringPtr("Foo"),
				},
				Volume: &corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/data/persistence"},
				},
			},
			Snapshot: &coh.PersistentStorageSpec{
				Enabled: boolPtr(false),
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
					StorageClassName: stringPtr("Bar"),
				},
				Volume: &corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/data/snapshot"},
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

func TestCreateStatefulSetFromRoleWithPersistenceDisabledAndSnapshotEnabled(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{
				Enabled: boolPtr(false),
			},
			Snapshot: &coh.PersistentStorageSpec{
				Enabled: boolPtr(true),
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"})
	// add the expected volume mount too the utils container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})
	// add the expected volume mount too the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})
	// add the expected snapshot PVC
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPVC
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   coh.VolumeNameSnapshots,
			Labels: labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithPersistenceDisabledAndSnapshotEnabledWithPVC(t *testing.T) {

	mode := corev1.PersistentVolumeFilesystem
	pvc := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"Foo": "Bar"},
		},
		VolumeName:       "test-volume",
		StorageClassName: stringPtr("Foo"),
		VolumeMode:       &mode,
	}

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{
				Enabled: boolPtr(false),
			},
			Snapshot: &coh.PersistentStorageSpec{
				Enabled:               boolPtr(true),
				PersistentVolumeClaim: &pvc,
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"})
	// add the expected volume mount too the utils container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})
	// add the expected volume mount too the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})
	// add the expected snapshot PVC
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPVC
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   coh.VolumeNameSnapshots,
			Labels: labels,
		},
		Spec: pvc,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithPersistenceDisabledAndSnapshotEnabledWithVolume(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{
				Enabled: boolPtr(false),
			},
			Snapshot: &coh.PersistentStorageSpec{
				Enabled: boolPtr(true),
				Volume: &corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/data/snapshots"},
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"})
	// add the expected volume mount too the utils container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})
	// add the expected volume mount too the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})
	// add the expected volume to the Pod
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: coh.VolumeNameSnapshots,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{Path: "/data/snapshots"},
		},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithPersistenceDisabledAndSnapshotEnabledWithVolumeAndPvcVolumeOnlyIsAdded(t *testing.T) {

	mode := corev1.PersistentVolumeFilesystem
	pvc := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"Foo": "Bar"},
		},
		VolumeName:       "test-volume",
		StorageClassName: stringPtr("Foo"),
		VolumeMode:       &mode,
	}

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{
				Enabled: boolPtr(false),
			},
			Snapshot: &coh.PersistentStorageSpec{
				Enabled: boolPtr(true),
				Volume: &corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/data/snapshots"},
				},
				PersistentVolumeClaim: &pvc,
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"})
	// add the expected volume mount too the utils container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})
	// add the expected volume mount too the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})
	// add the expected volume to the Pod
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: coh.VolumeNameSnapshots,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{Path: "/data/snapshots"},
		},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithPersistenceEnabledSnapshotDisabled(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{
				Enabled: boolPtr(true),
			},
			Snapshot: &coh.PersistentStorageSpec{
				Enabled: boolPtr(false),
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"})
	// add the expected volume mount too the utils container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})
	// add the expected volume mount too the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})
	// add the expected snapshot PVC
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPVC
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   coh.VolumeNamePersistence,
			Labels: labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithPersistenceEnabledWithPVCAndSnapshotDisabled(t *testing.T) {

	mode := corev1.PersistentVolumeFilesystem
	pvc := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"Foo": "Bar"},
		},
		VolumeName:       "test-volume",
		StorageClassName: stringPtr("Foo"),
		VolumeMode:       &mode,
	}

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{
				Enabled:               boolPtr(true),
				PersistentVolumeClaim: &pvc,
			},
			Snapshot: &coh.PersistentStorageSpec{
				Enabled: boolPtr(false),
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"})
	// add the expected volume mount too the utils container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})
	// add the expected volume mount too the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})
	// add the expected snapshot PVC
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPVC
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   coh.VolumeNamePersistence,
			Labels: labels,
		},
		Spec: pvc,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithPersistenceEnabledWithVolumeAndSnapshotDisabled(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{
				Enabled: boolPtr(true),
				Volume: &corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/data/persistence"},
				},
			},
			Snapshot: &coh.PersistentStorageSpec{
				Enabled: boolPtr(false),
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"})
	// add the expected volume mount too the utils container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})
	// add the expected volume mount too the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})
	// add the expected volume to the Pod
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: coh.VolumeNamePersistence,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{Path: "/data/persistence"},
		},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithPersistenceEnabledWithVolumeAndPvcAndSnapshotDisabledVolumeOnlyIsAdded(t *testing.T) {

	mode := corev1.PersistentVolumeFilesystem
	pvc := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"Foo": "Bar"},
		},
		VolumeName:       "test-volume",
		StorageClassName: stringPtr("Foo"),
		VolumeMode:       &mode,
	}

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{
				Enabled: boolPtr(true),
				Volume: &corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/data/snapshots"},
				},
				PersistentVolumeClaim: &pvc,
			},
			Snapshot: &coh.PersistentStorageSpec{
				Enabled: boolPtr(false),
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"})
	// add the expected volume mount too the utils container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})
	// add the expected volume mount too the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})
	// add the expected volume to the Pod
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: coh.VolumeNamePersistence,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{Path: "/data/snapshots"},
		},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithPersistenceEnabledSnapshotEnabled(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistentStorageSpec{
				Enabled: boolPtr(true),
			},
			Snapshot: &coh.PersistentStorageSpec{
				Enabled: boolPtr(true),
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"})
	// add the expected volume mount too the utils container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})
	// add the expected volume mount too the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})
	// add the expected snapshot PVC
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPVC
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   coh.VolumeNamePersistence,
			Labels: labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{},
	})
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   coh.VolumeNameSnapshots,
			Labels: labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}
