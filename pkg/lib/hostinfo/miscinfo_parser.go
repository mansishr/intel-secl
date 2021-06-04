/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package hostinfo

import (
	"runtime"

	model "github.com/intel-secl/intel-secl/v4/pkg/model/ta"
)

// miscInfoParser currenty collects the HostInfo's NumberOfSockets field.
type miscInfoParser struct{}

func (miscInfoParser *miscInfoParser) Init() error {
	return nil
}

func (miscInfoParser *miscInfoParser) Parse(hostInfo *model.HostInfo) error {
	hostInfo.NumberOfSockets = runtime.NumCPU()

	// Currently, hostinfo is not used on vmware and windows is not supported
	// (i.e., Linux only).
	hostInfo.OSType = model.OsTypeLinux
	return nil
}