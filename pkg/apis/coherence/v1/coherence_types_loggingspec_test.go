/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"os"
	"testing"
	"text/template"
)

var _ = Describe("Testing LoggingSpec struct", func() {

	Context("Copying a LoggingSpec using DeepCopyWithDefaults", func() {
		var original *coh.LoggingSpec
		var defaults *coh.LoggingSpec
		var clone *coh.LoggingSpec

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

		When("defaults is nil", func() {
			BeforeEach(func() {
				original = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: stringPtr("loggingMap"),
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(true)},
				}

				defaults = nil
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*original.ConfigMapName))
			})

			It("should copy the original Fluentd", func() {
				Expect(*clone.Fluentd).To(Equal(*original.Fluentd))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: stringPtr("loggingMap"),
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(true)},
				}

				original = nil
			})

			It("should copy the defaults ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*defaults.ConfigFile))
			})

			It("should copy the defaults ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*defaults.ConfigMapName))
			})

			It("should copy the defaults Fluentd", func() {
				Expect(*clone.Fluentd).To(Equal(*defaults.Fluentd))
			})
		})

		When("all original fields are set", func() {
			BeforeEach(func() {
				original = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: stringPtr("loggingMap"),
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(true)},
				}

				defaults = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging2.properties"),
					ConfigMapName: stringPtr("loggingMap2"),
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(false)},
				}
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*original.ConfigMapName))
			})

			It("should copy the original Fluentd", func() {
				Expect(*clone.Fluentd).To(Equal(*original.Fluentd))
			})
		})

		When("original Level is nil", func() {
			BeforeEach(func() {
				original = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: stringPtr("loggingMap"),
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(true)},
				}

				defaults = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging2.properties"),
					ConfigMapName: stringPtr("loggingMap2"),
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(false)},
				}
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*original.ConfigMapName))
			})

			It("should copy the original Fluentd", func() {
				Expect(*clone.Fluentd).To(Equal(*original.Fluentd))
			})
		})

		When("original ConfigFile is nil", func() {
			BeforeEach(func() {
				original = &coh.LoggingSpec{
					ConfigFile:    nil,
					ConfigMapName: stringPtr("loggingMap"),
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(true)},
				}

				defaults = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging2.properties"),
					ConfigMapName: stringPtr("loggingMap2"),
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(false)},
				}
			})

			It("should copy the defaults ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*defaults.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*original.ConfigMapName))
			})

			It("should copy the original Fluentd", func() {
				Expect(*clone.Fluentd).To(Equal(*original.Fluentd))
			})
		})

		When("original ConfigMapName is nil", func() {
			BeforeEach(func() {
				original = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: nil,
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(true)},
				}

				defaults = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging2.properties"),
					ConfigMapName: stringPtr("loggingMap2"),
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(false)},
				}
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the defaults ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*defaults.ConfigMapName))
			})

			It("should copy the original Fluentd", func() {
				Expect(*clone.Fluentd).To(Equal(*original.Fluentd))
			})
		})

		When("original Fluentd is nil", func() {
			BeforeEach(func() {
				original = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: stringPtr("loggingMap"),
					Fluentd:       nil,
				}

				defaults = &coh.LoggingSpec{
					ConfigFile:    stringPtr("logging2.properties"),
					ConfigMapName: stringPtr("loggingMap2"),
					Fluentd:       &coh.FluentdSpec{Enabled: boolPtr(false)},
				}
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*original.ConfigMapName))
			})

			It("should copy the defaults Fluentd", func() {
				Expect(*clone.Fluentd).To(Equal(*defaults.Fluentd))
			})
		})
	})
})

func TestLoggingSpec_CreateConfigMap(t *testing.T) {
	l := coh.LoggingConfigTemplate{
		ClusterName: "test-cluster",
		RoleName:    "storage",
		Logging:     &coh.LoggingSpec{},
	}

	temp, err := template.New("efk").Parse(coh.EfkConfig)
	if err != nil {
		panic(err)
	}
	err = temp.Execute(os.Stdout, l)
	if err != nil {
		panic(err)
	}

}
