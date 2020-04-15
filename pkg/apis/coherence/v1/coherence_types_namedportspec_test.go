/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"fmt"
	"github.com/go-test/deep"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

var _ = Describe("Testing NamedPortSpec struct", func() {
	tcp := corev1.ProtocolTCP
	udp := corev1.ProtocolUDP

	Context("Copying a NamedPortSpec using DeepCopyWithDefaults", func() {
		var original *coh.NamedPortSpec
		var defaults *coh.NamedPortSpec
		var clone *coh.NamedPortSpec
		var expected *coh.NamedPortSpec

		NewPortSpecOne := func() *coh.NamedPortSpec {
			return &coh.NamedPortSpec{
				Name: "foo",
				PortSpec: coh.PortSpec{
					Port: 8000,
				},
			}
		}

		NewPortSpecTwo := func() *coh.NamedPortSpec {
			return &coh.NamedPortSpec{
				Name: "bar",
				PortSpec: coh.PortSpec{
					Port: 9000,
				},
			}
		}

		ValidateResult := func() {
			It("should have correct Name", func() {
				Expect(clone.Name).To(Equal(expected.Name))
			})

			It("should have correct Protocol", func() {
				Expect(clone.Protocol).To(Equal(expected.Protocol))
			})

			It("should have correct PortSpec", func() {
				Expect(clone.PortSpec).To(Equal(expected.PortSpec))
			})
		}

		JustBeforeEach(func() {
			clone = original.DeepCopyWithDefaults(defaults)
		})

		When("original and defaults are nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = nil
			})

			It("the copy should be nil", func() {
				Expect(clone).Should(BeNil())
			})
		})

		When("defaults is nil copy should match original", func() {
			BeforeEach(func() {
				original = NewPortSpecOne()
				defaults = nil
				expected = original
			})

			ValidateResult()
		})

		When("original is nil copy should match defaults", func() {
			BeforeEach(func() {
				defaults = NewPortSpecOne()
				original = nil
				expected = defaults
			})

			ValidateResult()
		})

		When("all original fields are set copy should match original", func() {
			BeforeEach(func() {
				original = NewPortSpecOne()
				defaults = NewPortSpecTwo()
				expected = original
			})

			ValidateResult()
		})

		When("original Name is blank copy should have defaults name", func() {
			BeforeEach(func() {
				original = NewPortSpecOne()
				original.Name = ""
				defaults = NewPortSpecTwo()

				expected = NewPortSpecOne()
				expected.Name = defaults.Name
			})

			ValidateResult()
		})

		Context("Merging []NamedPortSpec", func() {
			var primary []coh.NamedPortSpec
			var secondary []coh.NamedPortSpec
			var merged []coh.NamedPortSpec

			var portOne = coh.NamedPortSpec{
				Name: "One",
				PortSpec: coh.PortSpec{
					Port:     7000,
					Protocol: &tcp,
				},
			}

			var portTwo = coh.NamedPortSpec{
				Name: "Two",
				PortSpec: coh.PortSpec{
					Port:     8000,
					Protocol: &udp,
				},
			}

			var portThree = coh.NamedPortSpec{
				Name: "Three",
				PortSpec: coh.PortSpec{
					Port:     9000,
					Protocol: &tcp,
				},
			}

			JustBeforeEach(func() {
				merged = coh.MergeNamedPortSpecs(primary, secondary)
			})

			When("primary and secondary slices are nil", func() {
				BeforeEach(func() {
					primary = nil
					secondary = nil
				})

				It("the result should be nil", func() {
					Expect(merged).To(BeNil())
				})
			})

			When("primary slice is not nil and the secondary slice is nil", func() {
				BeforeEach(func() {
					primary = []coh.NamedPortSpec{portOne, portTwo, portThree}
					secondary = nil
				})

				It("the result should be the primary slice", func() {
					Expect(merged).To(Equal(primary))
				})
			})

			When("primary slice is not nil and the secondary slice is empty", func() {
				BeforeEach(func() {
					primary = []coh.NamedPortSpec{portOne, portTwo, portThree}
					secondary = []coh.NamedPortSpec{}
				})

				It("the result should be the primary slice", func() {
					Expect(merged).To(Equal(primary))
				})
			})

			When("primary slice is nil and the secondary slice is not nil", func() {
				BeforeEach(func() {
					primary = nil
					secondary = []coh.NamedPortSpec{portOne, portTwo, portThree}
				})

				It("the result should be the secondary slice", func() {
					Expect(merged).To(Equal(secondary))
				})
			})

			When("primary slice is empty and the secondary slice is not nil", func() {
				BeforeEach(func() {
					primary = []coh.NamedPortSpec{}
					secondary = []coh.NamedPortSpec{portOne, portTwo, portThree}
				})

				It("the result should be the secondary slice", func() {
					Expect(merged).To(Equal(secondary))
				})
			})

			When("primary slice is populated and the secondary slice is populated", func() {
				BeforeEach(func() {
					primary = []coh.NamedPortSpec{portOne, portTwo}
					secondary = []coh.NamedPortSpec{portThree}
				})

				It("the result should contain the correct number of ports", func() {
					Expect(len(merged)).To(Equal(3))
				})

				It("the result should contain portOne at position 0", func() {
					Expect(merged[0]).To(Equal(portOne))
				})

				It("the result should contain portTwo at position 1", func() {
					Expect(merged[1]).To(Equal(portTwo))
				})

				It("the result should contain portThree at position 2", func() {
					Expect(merged[2]).To(Equal(portThree))
				})
			})

			When("primary slice is populated and the secondary slice is populated with matching ports", func() {
				var p1 = coh.NamedPortSpec{
					Name: "Foo",
					PortSpec: coh.PortSpec{
						Port: 7000,
					},
				}

				var p2 = coh.NamedPortSpec{
					Name: "Foo",
					PortSpec: coh.PortSpec{
						Protocol: &tcp,
					},
				}

				var pm = coh.NamedPortSpec{
					Name: "Foo",
					PortSpec: coh.PortSpec{
						Port:     7000,
						Protocol: &tcp,
					},
				}

				BeforeEach(func() {
					primary = []coh.NamedPortSpec{portOne, p1}
					secondary = []coh.NamedPortSpec{portTwo, p2}
				})

				It("the result should contain the correct number of ports", func() {
					Expect(len(merged)).To(Equal(3))
				})

				It("the result should contain portOne at position 0", func() {
					Expect(merged[0]).To(Equal(portOne))
				})

				It("the result should contain the merged port at position 1", func() {
					Expect(merged[1]).To(Equal(pm))
				})

				It("the result should contain portThree at position 2", func() {
					Expect(merged[2]).To(Equal(portTwo))
				})
			})
		})
	})
})

func TestNamedPortSpec_CreateServiceWithMinimalFields(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceCluster{}
	c.Name = "test-cluster"
	r := coh.CoherenceRoleSpec{}
	r.Role = "storage"

	np := coh.NamedPortSpec{
		Name: "foo",
		PortSpec: coh.PortSpec{
			Port: 19,
		},
	}

	labels := r.CreateCommonLabels(&c)
	labels[coh.LabelComponent] = fmt.Sprintf(coh.LabelComponentPortServiceTemplate, np.Name)

	selector := r.CreateCommonLabels(&c)
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod

	expected := corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:   "test-cluster-storage-foo",
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "foo",
					Protocol:   corev1.ProtocolTCP,
					Port:       19,
					TargetPort: intstr.FromString("foo"),
					NodePort:   0,
				},
			},
			Selector: selector,
		},
	}

	svc := np.CreateService(&c, &r)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(deep.Equal(*svc, expected)).To(BeNil())
}

func TestNamedPortSpec_CreateServiceWithProtocol(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceCluster{}
	c.Name = "test-cluster"
	r := coh.CoherenceRoleSpec{}
	r.Role = "storage"

	udp := corev1.ProtocolUDP

	np := coh.NamedPortSpec{
		Name: "foo",
		PortSpec: coh.PortSpec{
			Port:     19,
			Protocol: &udp,
		},
	}

	svc := np.CreateService(&c, &r)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(svc.Spec.Ports[0].Protocol).To(Equal(udp))
}

func TestNamedPortSpec_CreateServiceWithNodePort(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceCluster{}
	c.Name = "test-cluster"
	r := coh.CoherenceRoleSpec{}
	r.Role = "storage"

	np := coh.NamedPortSpec{
		Name: "foo",
		PortSpec: coh.PortSpec{
			Port:     19,
			NodePort: int32Ptr(6676),
		},
	}

	svc := np.CreateService(&c, &r)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(svc.Spec.Ports[0].NodePort).To(Equal(int32(6676)))
}

func TestNamedPortSpec_CreateServiceWithService(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceCluster{}
	c.Name = "test-cluster"
	r := coh.CoherenceRoleSpec{}
	r.Role = "storage"

	tp := corev1.ServiceTypeClusterIP
	ipf := corev1.IPv4Protocol
	etpt := corev1.ServiceExternalTrafficPolicyTypeLocal
	sa := corev1.ServiceAffinityClientIP
	sac := corev1.SessionAffinityConfig{
		ClientIP: &corev1.ClientIPConfig{TimeoutSeconds: int32Ptr(9876)},
	}

	np := coh.NamedPortSpec{
		Name: "foo",
		PortSpec: coh.PortSpec{
			Port: 19,
			Service: &coh.ServiceSpec{
				Name:                     stringPtr("bar"),
				Port:                     int32Ptr(99),
				Type:                     &tp,
				ClusterIP:                stringPtr("10.10.10.99"),
				ExternalIPs:              []string{"192.164.1.99", "192.164.1.100"},
				LoadBalancerIP:           stringPtr("10.10.10.10"),
				Labels:                   nil,
				Annotations:              nil,
				SessionAffinity:          &sa,
				LoadBalancerSourceRanges: []string{"A", "B"},
				ExternalName:             stringPtr("ext-bar"),
				ExternalTrafficPolicy:    &etpt,
				HealthCheckNodePort:      int32Ptr(1234),
				PublishNotReadyAddresses: boolPtr(true),
				SessionAffinityConfig:    &sac,
				IPFamily:                 &ipf,
			},
		},
	}

	labels := r.CreateCommonLabels(&c)
	labels[coh.LabelComponent] = fmt.Sprintf(coh.LabelComponentPortServiceTemplate, np.Name)

	selector := r.CreateCommonLabels(&c)
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod

	expected := corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:   "bar",
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "foo",
					Protocol:   corev1.ProtocolTCP,
					Port:       99,
					TargetPort: intstr.FromString("foo"),
				},
			},
			Selector:                 selector,
			ClusterIP:                "10.10.10.99",
			Type:                     tp,
			ExternalIPs:              []string{"192.164.1.99", "192.164.1.100"},
			SessionAffinity:          sa,
			LoadBalancerIP:           "10.10.10.10",
			LoadBalancerSourceRanges: []string{"A", "B"},
			ExternalName:             "ext-bar",
			ExternalTrafficPolicy:    etpt,
			HealthCheckNodePort:      1234,
			PublishNotReadyAddresses: true,
			SessionAffinityConfig:    &sac,
			IPFamily:                 &ipf,
		},
	}

	svc := np.CreateService(&c, &r)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(deep.Equal(*svc, expected)).To(BeNil())
}
