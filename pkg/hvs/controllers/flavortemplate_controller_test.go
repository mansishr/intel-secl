/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/controllers"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/domain/mocks"
	hvsRoutes "github.com/intel-secl/intel-secl/v3/pkg/hvs/router"
	consts "github.com/intel-secl/intel-secl/v3/pkg/lib/common/constants"
	"github.com/intel-secl/intel-secl/v3/pkg/model/hvs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlavorTemplateController", func() {
	var router *mux.Router
	var w *httptest.ResponseRecorder
	var flavorTemplateStore *mocks.MockFlavorTemplateStore
	var flavorTemplateController *controllers.FlavorTemplateController
	BeforeEach(func() {
		router = mux.NewRouter()
		flavorTemplateStore = mocks.NewFakeFlavorTemplateStore()

		flavorTemplateController = controllers.NewFlavorTemplateController(flavorTemplateStore,
			"../../../build/linux/hvs/schema/common.schema.json", "../../../build/linux/hvs/schema/flavor-template.json")
	})

	// Specs for HTTP Post to "/flavor-template"
	Describe("Post a new FlavorTemplate", func() {
		Context("Provide a valid FlavorTemplate data", func() {
			It("Should create a new Flavortemplate and get HTTP Status: 201", func() {
				router.Handle("/flavor-templates", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Create))).Methods("POST")
				flavorTemplateJson := `{
					"label": "test-uefi",
					"condition": [
						"//host_info/vendor='Linux'",
						"//host_info/tpm_version='2.0'",
						"//host_info/uefi_enabled='true'",
						"//host_info/suefi_enabled='true'"
					],
					"flavor_parts": {
						"PLATFORM": {
							"meta": {
								"tpm_version": "2.0",
								"uefi_enabled": true,
								"vendor": "Linux"
							},
							"pcr_rules": [
								{
									"pcr": {
										"index": 0,
										"bank": "SHA256"
									},
									"pcr_matches": true,
									"eventlog_equals": {}
								}
							]
						},
						"OS": {
							"meta": {
								"tpm_version": "2.0",
								"uefi_enabled": true,
								"vendor": "Linux"
							},
							"pcr_rules": [
								{
									"pcr": {
										"index": 7,
										"bank": "SHA256"
									},
									"pcr_matches": true,
									"eventlog_includes": [
										"shim",
										"db",
										"kek",
										"vmlinuz"
									]
								}
							]
						}
					}
				}`

				req, err := http.NewRequest(
					"POST",
					"/flavor-templates",
					strings.NewReader(flavorTemplateJson),
				)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusCreated))
			})
		})

		Context("Provide a FlavorTemplate data that contains invalid fileds, to validate against schema", func() {
			It("Should get HTTP Status: 400", func() {
				router.Handle("/flavor-templates", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Create))).Methods("POST")
				flavorgroupJson := `{
					"label": "",
					"condition": [
						"//host_info/vendor='Linux'",
						"//host_info/tpm_version='2.0'",
						"//host_info/uefi_enabled='true'",
						"//host_info/suefi_enabled='true'"
					],
					"flavor_parts": {
						"PLATFORM": {
							"meta": {
								"tpm_version": "2.0",
								"tboot_installed": true
							},
							"pcr_rules": [
								{
									"pcr": {
										"index": 0,
										"bank": "SHA256"
									},
									"pcr_matches": true
								}
							]
						}
					}
				}`

				req, err := http.NewRequest(
					"POST",
					"/flavor-templates",
					strings.NewReader(flavorgroupJson),
				)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(400))
			})
		})

		Context("Provide a empty data that should give bad request error", func() {
			It("Should get HTTP Status: 400", func() {
				router.Handle("/flavor-templates", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Create))).Methods("POST")
				req, err := http.NewRequest(
					"POST",
					"/flavor-templates",
					strings.NewReader(""),
				)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(400))
			})
		})

		Context("Provide a valid FlavorTemplate data without ACCEPT header", func() {
			It("Should give HTTP Status: 415", func() {
				router.Handle("/flavor-templates", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Create))).Methods("POST")
				flavorTemplateJson := `{
					"label": "default-uefi",
					"condition": [
						"//host_info/vendor='Linux'",
						"//host_info/tpm_version='2.0'",
						"//host_info/uefi_enabled='true'",
						"//host_info/suefi_enabled='true'"
					],
					"flavor_parts": {
						"PLATFORM": {
							"meta": {
								"tpm_version": "2.0",
								"uefi_enabled": true,
								"vendor": "Linux"
							},
							"pcr_rules": [
								{
									"pcr": {
										"index": 0,
										"bank": "SHA256"
									},
									"pcr_matches": true,
									"eventlog_equals": {}
								}
							]
						},
						"OS": {
							"meta": {
								"tpm_version": "2.0",
								"uefi_enabled": true,
								"vendor": "Linux"
							},
							"pcr_rules": [
								{
									"pcr": {
										"index": 7,
										"bank": "SHA256"
									},
									"pcr_matches": true,
									"eventlog_includes": [
										"shim",
										"db",
										"kek",
										"vmlinuz"
									]
								}
							]
						}
					}
				}`

				req, err := http.NewRequest(
					"POST",
					"/flavor-templates",
					strings.NewReader(flavorTemplateJson),
				)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusUnsupportedMediaType))
			})
		})

		Context("Provide a valid FlavorTemplate data without Content-Type header", func() {
			It("Should give HTTP Status: 415", func() {
				router.Handle("/flavor-templates", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Create))).Methods("POST")
				flavorTemplateJson := `{
					"label": "default-uefi",
					"condition": [
						"//host_info/vendor='Linux'",
						"//host_info/tpm_version='2.0'",
						"//host_info/uefi_enabled='true'",
						"//host_info/suefi_enabled='true'"
					],
					"flavor_parts": {
						"PLATFORM": {
							"meta": {
								"tpm_version": "2.0",
								"uefi_enabled": true,
								"vendor": "Linux"
							},
							"pcr_rules": [
								{
									"pcr": {
										"index": 0,
										"bank": "SHA256"
									},
									"pcr_matches": true,
									"eventlog_equals": {}
								}
							]
						},
						"OS": {
							"meta": {
								"tpm_version": "2.0",
								"uefi_enabled": true,
								"vendor": "Linux"
							},
							"pcr_rules": [
								{
									"pcr": {
										"index": 7,
										"bank": "SHA256"
									},
									"pcr_matches": true,
									"eventlog_includes": [
										"shim",
										"db",
										"kek",
										"vmlinuz"
									]
								}
							]
						}
					}
				}`

				req, err := http.NewRequest(
					"POST",
					"/flavor-templates",
					strings.NewReader(flavorTemplateJson),
				)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusUnsupportedMediaType))
			})
		})

	})

	// Specs for HTTP Post to "/flavor-template/{flavor-template-id}"
	Describe("Retrieve a FlavorTemplate", func() {
		Context("Retrieve data with valid FlavorTemplate ID", func() {
			It("Should retrieve Flavortemplate data and get HTTP Status: 200", func() {
				router.Handle("/flavor-templates/{id}", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Retrieve))).Methods("GET")
				req, err := http.NewRequest("GET", "/flavor-templates/426912bd-39b0-4daa-ad21-0c6933230b50", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})

		Context("Retrieve data with unavailable FlavorTemplate ID", func() {
			It("Should not retrieve Flavortemplate data and get HTTP Status: 404", func() {
				router.Handle("/flavor-templates/{id}", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Retrieve))).Methods("GET")
				req, err := http.NewRequest("GET", "/flavor-templates/73755fda-c910-46be-821f-e8ddeab189e9", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("Retrieve data with invalid FlavorTemplate ID", func() {
			It("Should not retrieve Flavortemplate data and get HTTP Status: 404", func() {
				router.Handle("/flavor-templates/{id}", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Retrieve))).Methods("GET")
				req, err := http.NewRequest("GET", "/flavor-templates/2f5cf0ec-0000-0000-0000-000000000000", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

	})

	Describe("Search And Delete Flavor Templates", func() {
		Context("When no filter arguments are passed", func() {
			It("All Flavor template records are returned", func() {
				router.Handle("/flavor-templates", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/flavor-templates", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var ft *[]hvs.FlavorTemplate
				err = json.Unmarshal(w.Body.Bytes(), &ft)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("When id parameter is added in search API", func() {
			It("Flavor template with the given uuid must be returned", func() {
				router.Handle("/flavor-templates/", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/flavor-templates/?id=426912bd-39b0-4daa-ad21-0c6933230b50", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var ft []hvs.FlavorTemplate
				err = json.Unmarshal(w.Body.Bytes(), &ft)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("When label parameter is added in search API", func() {
			It("Flavor template with the given label must be returned", func() {
				router.Handle("/flavor-templates/", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/flavor-templates/?label=default-uefi", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var ft []hvs.FlavorTemplate
				err = json.Unmarshal(w.Body.Bytes(), &ft)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("When flavorPartContains parameter is added in search API", func() {
			It("Flavor template with the given flavor part must be returned", func() {
				router.Handle("/flavor-templates/", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/flavor-templates/?flavorPartContains=OS", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var ft []hvs.FlavorTemplate
				err = json.Unmarshal(w.Body.Bytes(), &ft)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("When conditionContains parameter is added in search API", func() {
			It("Flavor template with the given condition must be returned", func() {
				router.Handle("/flavor-templates/", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/flavor-templates/?conditionContains=//host_info/uefi_enabled='true'", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var ft []hvs.FlavorTemplate
				err = json.Unmarshal(w.Body.Bytes(), &ft)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("Delete a template which is not in the database", func() {
			It("Appropriate error response should be returned", func() {
				router.Handle("/flavor-templates/{id}", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Delete))).Methods("DELETE")
				req, err := http.NewRequest("DELETE", "/flavor-templates/426912bd-39b0-4daa-ad21-0c6933230b51", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
		Context("Delete a template which is available in the database", func() {
			It("The template with the given uuid must be deleted", func() {
				router.Handle("/flavor-templates/{id}", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Delete))).Methods("DELETE")
				req, err := http.NewRequest("DELETE", "/flavor-templates/426912bd-39b0-4daa-ad21-0c6933230b50", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})
		Context("When includeDeleted parameter is added in search API", func() {
			It("All Flavor template records are returned", func() {
				router.Handle("/flavor-templates/", hvsRoutes.ErrorHandler(hvsRoutes.JsonResponseHandler(flavorTemplateController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/flavor-templates/?includeDeleted=true", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var ft []hvs.FlavorTemplate
				err = json.Unmarshal(w.Body.Bytes(), &ft)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(ft)).To(Equal(1))
			})
		})
	})
})