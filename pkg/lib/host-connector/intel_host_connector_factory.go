/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package host_connector

import (
	"github.com/pkg/errors"
	client "github.com/intel-secl/intel-secl/v3/pkg/clients/ta"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/host-connector/types"
	"net/url"
	"strings"
)

type IntelConnectorFactory struct {

}

func (icf *IntelConnectorFactory) GetHostConnector(vendorConnector types.VendorConnector, aasApiUrl, trustedCaCerts string) (HostConnector, error) {

	log.Trace("intel_host_connector_factory:GetHostConnector() Entering")
	defer log.Trace("intel_host_connector_factory:GetHostConnector() Leaving")
	baseURL := vendorConnector.Url
	if !strings.Contains(baseURL, "/v2") {
		baseURL = baseURL + "/v2"
	}
	taApiURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, errors.New("intel_host_connector_factory:GetHostConnector() error retrieving TA API URL")
	}

	taClient, err := client.NewTAClient(aasApiUrl, 
										taApiURL,
										vendorConnector.Configuration.Username,
										vendorConnector.Configuration.Password,
										trustedCaCerts)

	if err != nil {
		return nil, errors.Wrap(err, "intel_host_connector_factory:GetHostConnector() Could not create Trust Agent client")
	}

	log.Debug("intel_host_connector_factory:GetHostConnector() TA client created")
	return &IntelConnector{taClient}, nil
}
