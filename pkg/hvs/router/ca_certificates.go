/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package router

import (
	"github.com/gorilla/mux"
	"github.com/intel-secl/intel-secl/v5/pkg/hvs/controllers"
	"github.com/intel-secl/intel-secl/v5/pkg/lib/common/constants"
	"github.com/intel-secl/intel-secl/v5/pkg/lib/common/crypt"
)

func SetCaCertificatesRoutes(router *mux.Router, certStore *crypt.CertificatesStore) *mux.Router {
	defaultLog.Trace("router/ca_certificates:SetCaCertificatesRoutes() Entering")
	defer defaultLog.Trace("router/ca_certificates:SetCaCertificatesRoutes() Leaving")

	caCertController := controllers.CaCertificatesController{CertStore: certStore}

	router.Handle("/ca-certificates/{certType}", ErrorHandler(JsonResponseHandler(caCertController.Retrieve))).Methods("GET")
	router.Handle("/ca-certificates", ErrorHandler(ResponseHandler(caCertController.SearchPem))).Methods("GET").Headers("Accept", constants.HTTPMediaTypePemFile)
	router.Handle("/ca-certificates", ErrorHandler(JsonResponseHandler(caCertController.Search))).Methods("GET")
	return router
}
