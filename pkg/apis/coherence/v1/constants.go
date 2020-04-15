/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

const (
	LabelCoherenceDeployment = "coherenceDeployment"
	LabelCoherenceCluster    = "coherenceCluster"
	LabelCoherenceRole       = "coherenceRole"
	LabelComponent           = "component"
	LabelCoherenceWKAMember  = "coherenceWKAMember"

	LabelComponentCoherenceStatefulSet = "coherence"
	LabelComponentCoherencePod         = "coherencePod"
	LabelComponentCoherenceHeadless    = "coherence-headless"
	LabelComponentEfkConfig            = "coherence-efk-config"
	LabelComponentPVC                  = "coherence-volume"

	LabelComponentPortServiceTemplate = "coherence-service-%s"

	EfkConfigMapNameTemplate = "%s-efk-config"

	// The default k8s service account name.
	DefaultServiceAccount = "default"

	// The affinity topology key for fault domains.
	AffinityTopologyKey = "failure-domain.beta.kubernetes.io/zone"

	// Container Names
	ContainerNameCoherence   = "coherence"
	ContainerNameUtils       = "coherence-k8s-utils"
	ContainerNameApplication = "application"
	ContainerNameFluentd     = "fluentd"

	// Volume names
	VolumeNamePersistence   = "persistence-volume"
	VolumeNameSnapshots     = "snapshot-volume"
	VolumeNameLogs          = "log-dir"
	VolumeNameUtils         = "utils-dir"
	VolumeNameApplication   = "application-dir"
	VolumeNameJVM           = "jvm"
	VolumeNameScripts       = "coherence-scripts"
	VolumeNameFluentdConfig = "fluentd-coherence-conf"
	VolumeNameManagementSSL = "management-ssl-config"
	VolumeNameMetricsSSL    = "metrics-ssl-config"
	VolumeNameLoggingConfig = "logging-config"

	// Volume mount paths
	VolumeMountPathPersistence     = "/persistence"
	VolumeMountPathSnapshots       = "/snapshot"
	VolumeMountPathUtils           = "/utils"
	VolumeMountPathJVM             = "/jvm"
	VolumeMountPathLogs            = "/logs"
	VolumeMountPathScripts         = "/scripts"
	VolumeMountPathManagementCerts = "/coherence/certs/management"
	VolumeMountPathMetricsCerts    = "/coherence/certs/metrics"
	VolumeMountPathLoggingConfig   = "/loggingconfig"
	VolumeMountPathFluentdConfig   = "/fluentd/etc/fluentd-coherence.conf"

	VolumeMountSubPathFluentdConfig = "fluentd-coherence.conf"

	AppDir          = "/app"
	LibDir          = AppDir + "/lib"
	ConfDir         = AppDir + "/conf"
	ExternalAppDir  = "/u01/oracle/oracle_home/coherence/app"
	ExternalConfDir = ExternalAppDir + "/conf"
	ExternalLibDir  = ExternalAppDir + "/lib"

	// Port names
	PortNameCoherence = "coherence"
	PortNameDebug     = "debug-port"
	PortNameHealth    = "health"

	// Logging config secret name
	CoherenceMonitoringConfigName = "coherence-monitoring-config"

	// Logging config keys
	LoggingConfigKeyElasticSearchHost  = "elasticsearchhost"
	LoggingConfigElasticSearchPort     = "elasticsearchport"
	LoggingConfigElasticSearchUser     = "elasticsearchuser"
	LoggingConfigElasticSearchPassword = "elasticsearchpassword"

	DefaultLoggingConfig = "/scripts/logging.properties"

	DefaultDebugPort      int32 = 5005
	DefaultManagementPort int32 = 30000
	DefaultMetricsPort    int32 = 9612
	DefaultJmxmpPort      int32 = 9099

	ConfigMapNameOperatorConfig = "coherence-operator-config"
	ConfigMapNameScripts        = "coherence-operator-scripts"

	OperatorConfigKeyHost = "operatorhost"

	DefaultReadinessPath = "/ready"
	DefaultLivenessPath  = "/healthz"

	DefaultCommandApplication = "/utils/copy"

	DefaultFluentdImage = "fluent/fluentd-kubernetes-daemonset:v1.3.3-debian-elasticsearch-1.3"
)

const EfkConfig = `# Coherence fluentd configuration
{{- if .Logging.Fluentd }}
{{-   if .Logging.Fluentd.ConfigFile }}
@include {{ .Logging.Fluentd.ConfigFile }}
{{-   end }}
{{- end }}

# Ignore fluentd messages
<match fluent.**>
  @type null
</match>

# Coherence Logs
<source>
  @type tail
  path /logs/coherence-*.log
  pos_file /tmp/cohrence.log.pos
  read_from_head true
  tag coherence-cluster
  multiline_flush_interval 20s
  <parse>
    @type multiline
    format_firstline /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}/
    format1 /^(?<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3})\/(?<uptime>[0-9\.]+) (?<product>.+) <(?<level>[^\s]+)> \(thread=(?<thread>.+), member=(?<member>.+)\):[\S\s](?<log>.*)/
  </parse>
</source>

<filter coherence-cluster>
  @type record_transformer
  <record>
    cluster "{{ $.ClusterName }}"
    role "{{ .RoleName }}"
    host "#{ENV['HOSTNAME']}"
    pod-uid "#{ENV['COHERENCE_POD_ID']}"
  </record>
</filter>

<match coherence-cluster>
  @type elasticsearch
  host "#{ENV['ELASTICSEARCH_HOST']}"
  port "#{ENV['ELASTICSEARCH_PORT']}"
  user "#{ENV['ELASTICSEARCH_USER']}"
  password "#{ENV['ELASTICSEARCH_PASSWORD']}"
  logstash_format true
  logstash_prefix coherence-cluster
</match>

{{- if .Logging.Fluentd }}
{{-   if .Logging.Fluentd.tag }}
<match {{ .Logging.Fluentd.Tag }} >
  @type elasticsearch
  host "#{ENV['ELASTICSEARCH_HOST']}"
  port "#{ENV['ELASTICSEARCH_PORT']}"
  user "#{ENV['ELASTICSEARCH_USER']}"
  password "#{ENV['ELASTICSEARCH_PASSWORD']}"
  logstash_format true
  logstash_prefix {{ .Logging.Fluentd.Tag }}
</match>
{{-   end }}
{{- end }}`
