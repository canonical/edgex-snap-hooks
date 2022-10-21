// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package hooks

const (
	// AutostartConfig is a configuration key used indicate that a
	// service (application or device) should be autostarted on install
	AutostartConfig = "autostart"
	// EnvConfig is the prefix used for configure hook keys used for
	// EdgeX configuration overrides.
	EnvConfig = "env"
	// ProfileConfig is a configuration key that specifies a named
	// configuration profile
	ProfileConfig = "profile"
	// Deprecated: moved to snap
	// ServiceConsul is the service key for Consul.
	ServiceConsul = "consul"
	// Deprecated: moved to snap
	// ServiceRedis is the service key for Redis.
	ServiceRedis = "redis"
	// Deprecated: moved to snap
	// ServiceData is the service key for EdgeX Core Data.
	ServiceData = "core-data"
	// Deprecated: moved to snap
	// ServiceMetadata is the service key for EdgeX Core MetaData.
	ServiceMetadata = "core-metadata"
	// Deprecated: moved to snap
	// ServiceCommand is the service key for EdgeX Core Command.
	ServiceCommand = "core-command"
	// Deprecated: moved to snap
	// ServiceNotify is the service key for EdgeX Support Notifications.
	ServiceNotify = "support-notifications"
	// Deprecated: moved to snap
	// ServiceSched is the service key for EdgeX Support Scheduler.
	ServiceSched = "support-scheduler"
	// Deprecated
	// ServiceAppCfg is the service key for EdgeX App Service Configurable.
	ServiceAppCfg = "app-service-configurable"
	// Deprecated: app has been removed from edgexfoundry snap
	// ServiceDevVirt is the service key for EdgeX Device Virtual.
	ServiceDevVirt = "device-virtual"
	// Deprecated: moved to snap
	// ServiceSecStore is the service key for EdgeX Security Secret Store (aka Vault).
	ServiceSecStore = "security-secret-store"
	// Deprecated: moved to snap
	// ServiceProxy is the service key for EdgeX API Gateway (aka Kong).
	ServiceProxy = "security-proxy"
	// Deprecated: moved to snap
	// ServiceSysMgmt is the service key for EdgeX SMA (sys-mgmt-agent).
	// Deprecated: moved to snap
	ServiceSysMgmt = "sys-mgmt-agent"
	// ServiceKuiper is the service key for the Kuiper rules engine.
	// Deprecated: moved to snap
	ServiceKuiper = "kuiper"
	// Deprecated: moved to snap
	// ServiceSecBootstrapper is the service key for the security-bootstrapper,
	// a one-shot service that bootstraps per-service Consul ACLS required to
	// access Consul for configuration or registry services.
	ServiceSecBootstrapper = "security-bootstrapper"

	// Deprecated: use env package
	snapEnv         = "SNAP"
	snapCommonEnv   = "SNAP_COMMON"
	snapDataEnv     = "SNAP_DATA"
	snapInstNameEnv = "SNAP_INSTANCE_NAME"
	snapNameEnv     = "SNAP_NAME"
	snapRevEnv      = "SNAP_REVISION"
)

// Deprecated
// ConfToEnv defines mappings from snap config keys to EdgeX environment variable
// names that are used to override individual service configuration values via a
// .env file read by the snap service wrapper.
//
// The syntax to set a configuration key is:
//
// env.<service name>.<section>.<keyname>
//
var ConfToEnv = map[string]string{
	// [Writable] - not yet supported
	// conf_to_env["writable.log-level"]="BootTimeout"
	// See https://github.com/edgexfoundry/go-mod-bootstrap/blob/master/config/types.go

	// [Service]
	// HealthCheckInterval is the interval for Registry heal check callback
	"service.health-check-interval": "SERVICE_HEALTHCHECKINTERVAL",
	// Host is the hostname or IP address of the service.
	"service.host": "SERVICE_HOST",
	// Port is the HTTP port of the service.
	"service.port": "SERVICE_PORT",
	// ServerBindAddr specifies an IP address or hostname
	// for ListenAndServe to bind to, such as 0.0.0.0
	"service.server-bind-addr": "SERVICE_SERVERBINDADDR",
	// StartupMsg specifies a string to log once service
	// initialization and startup is completed.
	"service.startup-msg": "SERVICE_STARTUPMSG",
	// MaxResultCount specifies the maximum size list supported
	// in response to REST calls to other services.
	"service.max-result-count": "SERVICE_MAXRESULTCOUNT",
	// MaxRequestSize defines the maximum size of http request body in bytes
	"service.max-request-size": "SERVICE_MAXREQUESTSIZE",
	// RequestTimeout specifies a timeout (in milliseconds) for
	// processing REST request calls from other services.
	"service.request-timeout": "SERVICE_REQUESTTIMEOUT",

	// [SecretStore]
	"secret-store.secrets-file":               "SECRETSTORE_SECRETSFILE",
	"secret-store.disable-scrub-secrets-file": "SECRETSTORE_DISABLESCRUBSECRETSFILE",

	// [Registry] -- not yet supported, would also require consul changes

	// [Clients]
	// [Clients.core-command]
	"clients.core-command.port": "CLIENTS_CORECOMMAND_PORT",

	// [Clients.core-data]
	"clients.core-data.port": "CLIENTS_COREDATA_PORT",

	// [Clients.core-metadata]
	"clients.core-metadata.port": "CLIENTS_COREMETADATA_PORT",

	// [Clients.support-notifications]
	"clients.support-notifications.port": "CLIENTS_SUPPORTNOTIFICATIONS_PORT",

	// [Clients.support-scheduler]
	"clients.support-scheduler.port": "CLIENTS_SUPPORTSCHEDULER_PORT",

	// [Database] -- application services only; not supported
	// [Databases] -- not supported

	// [MessageQueue] -- core-data only
	// Indicates the message bus implementation to use, i.e. zero, mqtt, redisstreams...
	"messagequeue.type": "core-data,device-virtual/MESSAGEQUEUE_TYPE",
	// Protocol indicates the protocol to use when accessing the message bus.
	"messagequeue.protocol": "core-data,device-virtual/MESSAGEQUEUE_PROTOCOL",
	// Host is the hostname or IP address of the broker, if applicable.
	"messagequeue.host": "core-data,device-virtual/MESSAGEQUEUE_HOST",
	// Port defines the port on which to access the message bus.
	"messagequeue.port": "core-data,device-virtual/MESSAGEQUEUE_PORT",
	// PublishTopicPrefix indicates the topic prefix the data is published to.
	"messagequeue.publish-topic-prefix": "core-data,device-virtual/MESSAGEQUEUE_PUBLISHTOPICPREFIX",
	// SubscribeTopic indicates the topic in which to subscribe.
	"messagequeue.subscribe-topic": "core-data,device-virtual/MESSAGEQUEUE_SUBSCRIBETOPIC",
	// AuthMode specifies the type of secure connection to the message bus which are 'none', 'usernamepassword'
	// 'clientcert' or 'cacert'. Not all option supported by each implementation.
	// ZMQ doesn't support any Authmode beyond 'none', RedisStreams only supports 'none' & 'usernamepassword'
	// while MQTT supports all options.
	"messagequeue.auth-mode": "core-data,device-virtual/MESSAGEQUEUE_AUTHMODE",
	// SecretName is the name of the secret in the SecretStore that contains the Auth Credentials. The credential are
	// dynamically loaded using this name and store the Option property below where the implementation expected to
	// find them.
	"messagequeue.secret-name": "core-data,device-virtual/MESSAGEQUEUE_SECRETNAME",
	// SubscribeEnabled indicates whether enable the subscription to the Message Queue
	"messagequeue.subscribe-enabled": "core-data,device-virtual/MESSAGEQUEUE_SUBSCRIBEENABLED",

	// [MessageQueue.Optional] - not yet supported

	// [Trigger]
	// [Trigger.EdgexMessageBus]
	// [Trigger.EdgexMessageBus.SubscribeHost]
	"trigger.edgex-message-bus.subscribe-host.port":             "app-service-config/TRIGGER_EDGEXMESSAGEBUS_SUBSCRIBEHOST_PORT",
	"trigger.edgex-message-bus.subscribe-host.protocol":         "app-service-config/TRIGGER_EDGEXMESSAGEBUS_SUBSCRIBEHOST_PROTOCOL",
	"trigger.edgex-message-bus.subscribe-host.subscribe-topics": "app-service-config/TRIGGER_EDGEXMESSAGEBUS_SUBSCRIBEHOST_SUBSCRIBETOPICS",

	// [Trigger.EdgexMessageBus.PublishHost]
	"trigger.edgex-message-bus.publish-host.port":          "app-service-config/TRIGGER_EDGEXMESSAGEBUS_PUBLISHHOST_PORT",
	"trigger.edgex-message-bus.publish-host.protocol":      "app-service-config/TRIGGER_EDGEXMESSAGEBUS_PUBLISHHOST_PROTOCOL",
	"trigger.edgex-message-bus.publish-host.publish-topic": "app-service-config/TRIGGER_EDGEXMESSAGEBUS_PUBLISHHOST_PUBLISHTOPIC",

	// [Smtp]
	"smtp.host":                    "support-notifications/SMTP_HOST",
	"smtp.username":                "support-notifications/SMTP_USERNAME",
	"smtp.password":                "support-notifications/SMTP_PASSWORD",
	"smtp.port":                    "support-notifications/SMTP_PORT",
	"smtp.sender":                  "support-notifications/SMTP_SENDER",
	"smtp.enable-self-signed-cert": "support-notifications/SMTP_ENABLE_SELF_SIGNED_CERT",
	"smtp.subject":                 "support-notifications/SMTP_SUBJECT",

	// AuthMode is the SMTP authentication mechanism. Currently, 'usernamepassword' is the only
	// AuthMode supported by this service, and the secret keys are 'username' and 'password'.
	"smtp.auth-mode": "support-notifications/SMTP_AUTHMODE",

	// ADD_PROXY_ROUTE is a csv list of URLs to be added to the
	// API Gateway (aka Kong). For references:
	//
	// https://docs.edgexfoundry.org/1.3/microservices/security/Ch-APIGateway/
	//
	// NOTE - this setting is not a configuration override, it's a top-level
	// environment variable used by the security-proxy-setup.
	//
	// TODO: validation
	//
	"add-proxy-route": "security-proxy/ADD_PROXY_ROUTE",

	// [KongAuth]
	"kongauth.name": "security-proxy/KONGAUTH_NAME",

	// NOTE - these settings are not configuration overrides, they are environment
	// variables ready by security-secretstore-setup and security-bootstrapper on
	// startup
	//
	// TODO: validation
	//
	// ADD_SECRETSTORE_TOKENS is a csv list of service keys to be added to the
	// list of Vault tokens that security-file-token-provider (launched by
	// security-secretstore-setup) creates.
	//
	"add-secretstore-tokens": "security-secret-store/ADD_SECRETSTORE_TOKENS",

	// ADD_KNOWN_SECRETS is a csv list of service keys and list of known secrets
	// to be copied into the Vault namespace for the service. The primary use for
	// this variable is ensuring that the redis password is accessible.
	"add-known-secrets": "security-secret-store/ADD_KNOWN_SECRETS",

	// The default-token-ttl setting is a Go Duration string, a sequence of decimal
	// numbers, each with optional fraction and a unit suffix (e.g. "ns", "us" (or
	// "µs"), "ms", "s", "m", "h"). It's used to set the TTL of vault tokens generated
	// for EdgeX services during bootstrap. This setting can be used to increase
	// (or decrease) the default TTL (one hour). If the TTL of a token expires before
	// a service is started, the service will not be able to access the Secret Store.
	"default-token-ttl": "security-secret-store/TOKENFILEPROVIDER_DEFAULTTOKENTTL",

	// ADD_REGISTRY_ACL_ROLES is a csv list of service keys used to create
	// ACL roles in Vault to allow secure Consul access for the services.
	"add-registry-acl-roles": "security-bootstrapper/ADD_REGISTRY_ACL_ROLES",
}

// Deprecated: moved to snap
// Services is a string array of all of the edgexfoundry snap services.
var Services = []string{
	// base services
	ServiceConsul,
	ServiceRedis,
	// core services
	ServiceData,
	ServiceMetadata,
	ServiceCommand,
	// support services
	ServiceNotify,
	ServiceSched,
	// app-services
	ServiceAppCfg,
	// device services
	ServiceDevVirt,
	// security services
	ServiceSecStore,
	ServiceProxy,
	// sys mgmt services
	ServiceSysMgmt,
	// rules-engine
	ServiceKuiper,
	ServiceSecBootstrapper,
}
