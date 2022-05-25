/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package controllers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gorilla/mux"
	consts "github.com/intel-secl/intel-secl/v5/pkg/lib/common/constants"
	"github.com/intel-secl/intel-secl/v5/pkg/lib/common/context"
	ct "github.com/intel-secl/intel-secl/v5/pkg/model/aas"
	"github.com/intel-secl/intel-secl/v5/pkg/tagent/common"
	"github.com/intel-secl/intel-secl/v5/pkg/tagent/config"
	"github.com/intel-secl/intel-secl/v5/pkg/tagent/constants"
	"github.com/intel-secl/intel-secl/v5/pkg/tagent/controllers"
	tagentRouter "github.com/intel-secl/intel-secl/v5/pkg/tagent/router"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

const (
	testConfig      = "../test/resources/config.yaml"
	testConfig_test = "../test/resources/config-test.yaml"
)

var _ = Describe("GetAik Request", func() {
	var router *mux.Router
	var w *httptest.ResponseRecorder

	// Read Config
	testCfg, err := os.ReadFile(testConfig)
	if err != nil {
		log.Fatalf("Failed to load test tagent config file %v", err)
	}
	var tagentConfig *config.TrustAgentConfiguration
	yaml.Unmarshal(testCfg, &tagentConfig)

	testConfig_test, err := os.ReadFile(testConfig_test)
	if err != nil {
		log.Fatalf("Failed to load test tagent config file %v", err)
	}
	var testConfig *config.TrustAgentConfiguration
	yaml.Unmarshal(testConfig_test, &testConfig)

	var reqHandler common.RequestHandler
	var negReqHandler common.RequestHandler

	BeforeEach(func() {
		router = mux.NewRouter()
		reqHandler = common.NewMockRequestHandler(tagentConfig)
		negReqHandler = common.NewMockRequestHandler(testConfig)
	})

	Describe("GetAik", func() {
		Context("GetAik request", func() {
			It("Should get Aik", func() {
				router.HandleFunc("/v2/aik", tagentRouter.ErrorHandler(tagentRouter.RequiresPermission(
					controllers.GetAik(reqHandler), []string{"aik:retrieve"}))).Methods(http.MethodGet)

				req, err := http.NewRequest(http.MethodGet, "/v2/aik", nil)
				Expect(err).NotTo(HaveOccurred())

				permissions := ct.PermissionInfo{
					Service: constants.TAServiceName,
					Rules:   []string{"aik:retrieve"},
				}
				req = context.SetUserPermissions(req, []ct.PermissionInfo{permissions})
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})

		Context("Invalid RequestHandler in GetAik request", func() {
			It("Should not get Aik - Invalid RequestHandler", func() {
				router.HandleFunc("/v2/aik", tagentRouter.ErrorHandler(tagentRouter.RequiresPermission(
					controllers.GetAik(negReqHandler), []string{"aik:retrieve"}))).Methods(http.MethodGet)

				req, err := http.NewRequest(http.MethodGet, "/v2/aik", nil)
				Expect(err).NotTo(HaveOccurred())

				permissions := ct.PermissionInfo{
					Service: constants.TAServiceName,
					Rules:   []string{"aik:retrieve"},
				}
				req = context.SetUserPermissions(req, []ct.PermissionInfo{permissions})
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("Invalid Content-Type in get Aik request", func() {
			It("Should not get Aik - Invalid Content-Type", func() {
				router.HandleFunc("/v2/aik", tagentRouter.ErrorHandler(tagentRouter.RequiresPermission(
					controllers.GetAik(reqHandler), []string{"aik:retrieve"}))).Methods(http.MethodGet)

				req, err := http.NewRequest(http.MethodGet, "/v2/aik", nil)
				Expect(err).NotTo(HaveOccurred())

				permissions := ct.PermissionInfo{
					Service: constants.TAServiceName,
					Rules:   []string{"aik:retrieve"},
				}
				req = context.SetUserPermissions(req, []ct.PermissionInfo{permissions})
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
