/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"fmt"
	"github.com/go-test/deep"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sort"
	"testing"
)

func TestCreateServicesFromRoleWithAdditionalPortsEmpty(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Ports: []coh.NamedPortSpec{},
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// assert that the Services are as expected
	assertService(t, role, cluster)
}

func TestCreateServicesFromRoleWithPortsWithOneAdditionalPortWithServiceEnabledFalse(t *testing.T) {

	protocol := corev1.ProtocolUDP

	role := coh.CoherenceRoleSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				PortSpec: coh.PortSpec{
					Port:     9876,
					Protocol: &protocol,
					NodePort: int32Ptr(2020),
					HostPort: int32Ptr(1234),
					HostIP:   stringPtr("10.10.1.0"),
					Service: &coh.ServiceSpec{
						Enabled: boolPtr(false),
					},
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// assert that the Services are as expected
	assertService(t, role, cluster)
}

func TestCreateServicesFromRoleWithPortsWithOneAdditionalPort(t *testing.T) {

	protocol := corev1.ProtocolUDP

	role := coh.CoherenceRoleSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				PortSpec: coh.PortSpec{
					Port:     9876,
					Protocol: &protocol,
					NodePort: int32Ptr(2020),
					HostPort: int32Ptr(1234),
					HostIP:   stringPtr("10.10.1.0"),
					Service: &coh.ServiceSpec{
						Enabled: boolPtr(true),
					},
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// Create the expected labels
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPortService

	// Create the expected service selector labels
	selectorLabels := role.CreateCommonLabels(cluster)
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-%s-test-port-one", cluster.Name, role.GetRoleName()),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Protocol:   protocol,
					Port:       9876,
					TargetPort: intstr.FromString("test-port-one"),
					NodePort:   2020,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, role, cluster, &svcExpected)
}

func TestCreateServicesFromRoleWithPortsWithOneAdditionalPortWithServiceName(t *testing.T) {

	protocol := corev1.ProtocolUDP

	role := coh.CoherenceRoleSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				PortSpec: coh.PortSpec{
					Port:     9876,
					Protocol: &protocol,
					NodePort: int32Ptr(2020),
					HostPort: int32Ptr(1234),
					HostIP:   stringPtr("10.10.1.0"),
					Service: &coh.ServiceSpec{
						Enabled: boolPtr(true),
						Name:    stringPtr("test-service"),
					},
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// Create the expected labels
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPortService

	// Create the expected service selector labels
	selectorLabels := role.CreateCommonLabels(cluster)
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test-service",
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Protocol:   protocol,
					Port:       9876,
					TargetPort: intstr.FromString("test-port-one"),
					NodePort:   2020,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, role, cluster, &svcExpected)
}

func TestCreateServicesFromRoleWithPortsWithOneAdditionalPortWithServicePort(t *testing.T) {

	protocol := corev1.ProtocolUDP

	role := coh.CoherenceRoleSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				PortSpec: coh.PortSpec{
					Port:     9876,
					Protocol: &protocol,
					NodePort: int32Ptr(2020),
					HostPort: int32Ptr(1234),
					HostIP:   stringPtr("10.10.1.0"),
					Service: &coh.ServiceSpec{
						Enabled: boolPtr(true),
						Port:    int32Ptr(80),
					},
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// Create the expected labels
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPortService

	// Create the expected service selector labels
	selectorLabels := role.CreateCommonLabels(cluster)
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-%s-test-port-one", cluster.Name, role.GetRoleName()),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Protocol:   protocol,
					Port:       80,
					TargetPort: intstr.FromString("test-port-one"),
					NodePort:   2020,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, role, cluster, &svcExpected)
}

func TestCreateServicesFromRoleWithPortsWithOneAdditionalPortWithServiceFields(t *testing.T) {

	protocol := corev1.ProtocolUDP
	svcType := corev1.ServiceTypeNodePort
	trafficPolicy := corev1.ServiceExternalTrafficPolicyTypeLocal
	ipFamily := corev1.IPv4Protocol
	affinity := corev1.ServiceAffinityNone
	cfg := corev1.SessionAffinityConfig{
		ClientIP: &corev1.ClientIPConfig{
			TimeoutSeconds: int32Ptr(1000),
		},
	}

	role := coh.CoherenceRoleSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				PortSpec: coh.PortSpec{
					Port:     9876,
					Protocol: &protocol,
					NodePort: int32Ptr(2020),
					HostPort: int32Ptr(1234),
					HostIP:   stringPtr("10.10.1.0"),
					Service: &coh.ServiceSpec{
						Enabled:                  boolPtr(true),
						Type:                     &svcType,
						ClusterIP:                stringPtr("192.168.1.30"),
						ExternalIPs:              []string{"10.10.10.99", "10.10.10.100"},
						LoadBalancerIP:           stringPtr("10.99.0.0"),
						LoadBalancerSourceRanges: []string{"10.10.10.0", "10.10.10.255"},
						ExternalName:             stringPtr("test-external-name"),
						HealthCheckNodePort:      int32Ptr(1000),
						PublishNotReadyAddresses: boolPtr(true),
						ExternalTrafficPolicy:    &trafficPolicy,
						SessionAffinity:          &affinity,
						SessionAffinityConfig:    &cfg,
						IPFamily:                 &ipFamily,
					},
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// Create the expected labels
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPortService

	// Create the expected service selector labels
	selectorLabels := role.CreateCommonLabels(cluster)
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-%s-test-port-one", cluster.Name, role.GetRoleName()),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Protocol:   protocol,
					Port:       9876,
					TargetPort: intstr.FromString("test-port-one"),
					NodePort:   2020,
				},
			},
			Selector:                 selectorLabels,
			Type:                     svcType,
			ClusterIP:                "192.168.1.30",
			ExternalIPs:              []string{"10.10.10.99", "10.10.10.100"},
			LoadBalancerIP:           "10.99.0.0",
			LoadBalancerSourceRanges: []string{"10.10.10.0", "10.10.10.255"},
			ExternalName:             "test-external-name",
			HealthCheckNodePort:      1000,
			PublishNotReadyAddresses: true,
			ExternalTrafficPolicy:    trafficPolicy,
			SessionAffinity:          affinity,
			SessionAffinityConfig:    &cfg,
			IPFamily:                 &ipFamily,
		},
	}

	// assert that the Services are as expected
	assertService(t, role, cluster, &svcExpected)
}

func TestCreateServicesFromRoleWithPortsWithOneAdditionalPortWithServiceLabels(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				PortSpec: coh.PortSpec{
					Port: 9876,
					Service: &coh.ServiceSpec{
						Enabled: boolPtr(true),
						Labels:  map[string]string{"LabelOne": "One", "LabelTwo": "Two"},
					},
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// Create the expected labels
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPortService
	labels["LabelOne"] = "One"
	labels["LabelTwo"] = "Two"

	// Create the expected service selector labels
	selectorLabels := role.CreateCommonLabels(cluster)
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-%s-test-port-one", cluster.Name, role.GetRoleName()),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Port:       9876,
					TargetPort: intstr.FromString("test-port-one"),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, role, cluster, &svcExpected)
}

func TestCreateServicesFromRoleWithPortsWithOneAdditionalPortWithServiceAnnotations(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				PortSpec: coh.PortSpec{
					Port: 9876,
					Service: &coh.ServiceSpec{
						Enabled:     boolPtr(true),
						Annotations: map[string]string{"AnnOne": "One", "AnnTwo": "Two"},
					},
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// Create the expected labels
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPortService

	// Create the expected service selector labels
	selectorLabels := role.CreateCommonLabels(cluster)
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-%s-test-port-one", cluster.Name, role.GetRoleName()),
			Labels:      labels,
			Annotations: map[string]string{"AnnOne": "One", "AnnTwo": "Two"},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Port:       9876,
					TargetPort: intstr.FromString("test-port-one"),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, role, cluster, &svcExpected)
}

func TestCreateServicesFromRoleWithPortsWithTwoAdditionalPorts(t *testing.T) {

	protocolOne := corev1.ProtocolUDP
	protocolTwo := corev1.ProtocolSCTP

	role := coh.CoherenceRoleSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				PortSpec: coh.PortSpec{
					Port:     9876,
					Protocol: &protocolOne,
				},
			},
			{
				Name: "test-port-two",
				PortSpec: coh.PortSpec{
					Port:     5678,
					Protocol: &protocolTwo,
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// Create the expected labels
	labels := role.CreateCommonLabels(cluster)
	labels[coh.LabelComponent] = coh.LabelComponentPortService

	// Create the expected service selector labels
	selectorLabels := role.CreateCommonLabels(cluster)
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected first Service
	svcExpectedOne := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-%s-test-port-one", cluster.Name, role.GetRoleName()),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Port:       9876,
					TargetPort: intstr.FromString("test-port-one"),
					Protocol:   protocolOne,
				},
			},
			Selector: selectorLabels,
		},
	}

	// Create expected second Service
	svcExpectedTwo := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-%s-test-port-two", cluster.Name, role.GetRoleName()),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-two",
					Port:       5678,
					TargetPort: intstr.FromString("test-port-two"),
					Protocol:   protocolTwo,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, role, cluster, &svcExpectedOne, &svcExpectedTwo)
}

func assertService(t *testing.T, role coh.CoherenceRoleSpec, cluster *coh.CoherenceCluster, servicesExpected ...metav1.Object) {
	g := NewGomegaWithT(t)

	servicesActual := role.CreateServicesForPort(cluster)

	// Sort the expected services
	sort.SliceStable(servicesExpected, func(i, j int) bool {
		return servicesExpected[i].GetName() < servicesExpected[j].GetName()
	})

	// Sort the actual services
	sort.SliceStable(servicesActual, func(i, j int) bool {
		return servicesActual[i].GetName() < servicesActual[j].GetName()
	})

	diffs := deep.Equal(servicesActual, servicesExpected)
	g.Expect(diffs).To(BeNil())
}
