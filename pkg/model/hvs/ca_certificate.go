/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package hvs

// CaCertificateCollection is collection of certificates/ca certificates
type CaCertificateCollection struct {
	CaCerts []*CaCertificate `json:"ca_certificate,omitempty"`
}

// CaCertificate is to represent certificate/ca certificate
type CaCertificate struct {
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Certificate []byte `json:"certificate,omitempty"`
}