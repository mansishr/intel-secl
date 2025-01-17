/*
Copyright © 2021 Intel Corporation
SPDX-License-Identifier: BSD-3-Clause
*/

package constants

const (
	ExplicitServiceName = "ISecL K8s Admission Controller"
)

const (
	LogLevelEnv        = "LOG_LEVEL"
	LogMaxLengthEnv    = "LOG_MAX_LENGTH"
	PortEnv            = "PORT"
	DefaultLogFilePath = "/var/log/admission-controller/admission-controller.log"
	LogBasePath        = "/var/log/isecl-k8s-extensions/"
)

var (
	HttpLogFile = "/var/log/admission-controller/admission-controller-http.log"
)

const (
	LogLevelDefault     = "INFO"
	LogMaxLengthDefault = 1500
	PortDefault         = 8889
	TlsCertPath         = "/etc/webhook/certs/tls.crt"
	TlsKeyPath          = "/etc/webhook/certs/tls.key"
)

const (
	TaintNameNoschedule   = "untrusted"
	TaintNameNoexecute    = "untrusted"
	TaintEffectNoSchedule = "NoSchedule"
	TaintEffectNoExecute  = "NoExecute"
	TaintValueTrue        = "true"
)

const (
	MutateRoute = "/mutate"
)
