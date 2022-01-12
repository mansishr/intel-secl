/*
* Copyright (C) 2020 Intel Corporation
* SPDX-License-Identifier: BSD-3-Clause
 */
package mocks

import (
	"github.com/intel-secl/intel-secl/v5/pkg/hvs/domain/models"
	"github.com/intel-secl/intel-secl/v5/pkg/lib/common/crypt"
	"github.com/intel-secl/intel-secl/v5/pkg/lib/common/log"
)

var defaultLog = log.GetDefaultLogger()
var secLog = log.GetSecurityLogger()

func NewFakeCertificatesPathStore() *crypt.CertificatesPathStore {
	// For ECA, to read list of certificates from directory
	ecCaPath := "../domain/mocks/resources/"
	// Mock path to create new certificate
	rootCaPath := "../domain/mocks/resources/"
	//Path for any certificate containing file, so instead of creating new use existing one
	caCertPath := "../domain/mocks/resources/EndorsementCA-external.pem"

	return &crypt.CertificatesPathStore{
		models.CaCertTypesRootCa.String(): crypt.CertLocation{
			CertPath: rootCaPath,
		},
		models.CaCertTypesEndorsementCa.String(): crypt.CertLocation{
			CertPath: ecCaPath,
		},
		models.CaCertTypesPrivacyCa.String(): crypt.CertLocation{
			CertPath: caCertPath,
		},
		models.CaCertTypesTagCa.String(): crypt.CertLocation{
			CertPath: caCertPath,
		},
		models.CertTypesSaml.String(): crypt.CertLocation{
			CertPath: caCertPath,
		},
		models.CertTypesTls.String(): crypt.CertLocation{
			CertPath: caCertPath,
		},
	}
}

func NewFakeCertificatesStore() *crypt.CertificatesStore {

	// Mock path to create new certificate
	rootCaPath := "../domain/mocks/resources/"
	//Path for any certificate containing file, so instead of creating new use existing one
	caCertPath := "../domain/mocks/resources/EndorsementCA-external.pem"

	return &crypt.CertificatesStore{
		models.CaCertTypesRootCa.String(): &crypt.CertificateStore{
			CertPath:     rootCaPath,
			Certificates: nil,
		},
		models.CaCertTypesEndorsementCa.String(): &crypt.CertificateStore{
			CertPath:     rootCaPath,
			Certificates: nil,
		},
		models.CaCertTypesPrivacyCa.String(): &crypt.CertificateStore{
			CertPath:     caCertPath,
			Certificates: nil,
		},
		models.CaCertTypesTagCa.String(): &crypt.CertificateStore{
			CertPath:     caCertPath,
			Certificates: nil,
		},
		models.CertTypesSaml.String(): &crypt.CertificateStore{
			CertPath:     caCertPath,
			Certificates: nil,
		},
		models.CertTypesTls.String(): &crypt.CertificateStore{
			CertPath:     caCertPath,
			Certificates: nil,
		},
		models.CertTypesFlavorSigning.String(): &crypt.CertificateStore{
			Key:          nil,
			CertPath:     caCertPath,
			Certificates: nil,
		},
	}
}
