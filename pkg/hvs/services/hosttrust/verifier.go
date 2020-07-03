/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package hosttrust

import (
	"github.com/google/uuid"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/domain"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/domain/models"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/host-connector/types"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/saml"
	flavorVerifier "github.com/intel-secl/intel-secl/v3/pkg/lib/verifier"
	"github.com/intel-secl/intel-secl/v3/pkg/model/hvs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrInvalidHostManiFest = errors.New("invalid host data")
var ErrManifestMissingHwUUID = errors.New("host data missing hardware uuid")
var ErrMissingHostId = errors.New("host id ")

type verifier struct {
	flavorStore      domain.FlavorStore
	flavorGroupStore domain.FlavorGroupStore
	hostStore        domain.HostStore
	reportStore      domain.ReportStore
	flavorVerifier   flavorVerifier.Verifier
	certsStore       models.CertificatesStore
	samlIssuer       saml.IssuerConfiguration
}

func NewVerifier(cfg domain.HostTrustVerifierConfig) domain.HostTrustVerifier {
	return &verifier{
		flavorStore:      cfg.FlavorStore,
		flavorGroupStore: cfg.FlavorGroupStore,
		hostStore:        cfg.HostStore,
		reportStore:      cfg.ReportStore,
		flavorVerifier:   cfg.FlavorVerifier,
		certsStore:       cfg.CertsStore,
		samlIssuer:       cfg.SamlIssuerConfig,
	}

}

func (v *verifier) Verify(hostId uuid.UUID, hostData *types.HostManifest, newData bool) error {
	defaultLog.Trace("hosttrust/verifier:Verify() Entering")
	defer defaultLog.Trace("hosttrust/verifier:Verify() Leaving")

	if hostData == nil {
		return ErrInvalidHostManiFest
	}
	//TODO: Fix HardwareUUID has to be uuid
	hwUuid, err := uuid.Parse(hostData.HostInfo.HardwareUUID)
	if err != nil || hwUuid == uuid.Nil {
		return ErrManifestMissingHwUUID
	}

	// TODO : remove this when we remove the intermediate collection
	var flvGroups []*hvs.FlavorGroup
	if flvGroupColl, err := v.flavorGroupStore.Search(&models.FlavorGroupFilterCriteria{HostId: hostId.String()}); err != nil {
		return errors.New("hosttrust/verifier:Verify() Store access error")
	} else {
		flvGroups = (*flvGroupColl).Flavorgroups
	}

	// start with the presumption that final trust report would be true. It as some point, we get an invalid report,
	// the Overall trust status would be negative
	var finalReportValid = true // This is the final trust report - initialize
	// create an empty trust report with the host manifest
	finalTrustReport := hvs.TrustReport{HostManifest: *hostData}

	for _, fg := range flvGroups {
		//TODO - handle errors in case of DB transaction
		fgTrustReqs, _ := NewFlvGrpHostTrustReqs(hostId, hwUuid, *fg, v.flavorStore)
		fgCachedFlavors, _ := v.getCachedFlavors(hostId, (*fg).ID)
		if len(fgCachedFlavors) > 0 {
			fgTrustCache, _ := v.validateCachedFlavors(hostId, hostData, fgCachedFlavors)
			fgTrustReport := fgTrustCache.trustReport
			if !fgTrustReqs.MeetsFlavorGroupReqs(fgTrustCache) {
				finalReportValid = false
				fgTrustReport, _ = v.createFlavorGroupReport(hostId, *fgTrustReqs, hostData, fgTrustCache)
			}
			log.Debug("hosttrust/verifier:Verify() Trust status for host id", hostId, "for flavorgroup ", fg.ID, "is", fgTrustReport.Trusted)
			// append the results
			finalTrustReport.Results = append(finalTrustReport.Results, fgTrustReport.Results...)
		}
	}
	// create a new report if we actually have any results and either the Final Report is untrusted or
	// we have new Data from the host and therefore need to update based on the new report.
	if len(finalTrustReport.Results) > 0 && !finalReportValid || newData {
		log.Debugf("hosttrust/verifier:Verify() Generating new SAML for host: %s", hostId)
		samlReportGen := NewSamlReportGenerator(&v.samlIssuer)
		samlReport := samlReportGen.GenerateSamlReport(&finalTrustReport)

		log.Debugf("hosttrust/verifier:Verify() Saving new report for host: %s", hostId)
		v.storeTrustReport( hostId, &finalTrustReport, &samlReport)
	}
	return nil
}

func (v *verifier) getCachedFlavors(hostId uuid.UUID, flavGrpId uuid.UUID) ([]hvs.SignedFlavor, error) {
	defaultLog.Trace("hosttrust/verifier:getCachedFlavors() Entering")
	defer defaultLog.Trace("hosttrust/verifier:getCachedFlavors() Leaving")
	// retrieve the IDs of the trusted flavors from the host store
	if flIds, err := v.hostStore.RetrieveTrustCacheFlavors(hostId, flavGrpId); err != nil {
		return nil, errors.Wrap(err, "hosttrust/verifier:Verify() Error while retrieving TrustCacheFlavors")
	} else {
		result := make([]hvs.SignedFlavor, 0, len(flIds))
		for _, flvId := range flIds {
			if flv, err := v.flavorStore.Retrieve(flvId); err == nil {
				result = append(result, *flv)
			}
		}
		return result, nil
	}
}

func (v *verifier) validateCachedFlavors(hostId uuid.UUID,
	hostData *types.HostManifest,
	cachedFlavors []hvs.SignedFlavor) (hostTrustCache, error) {
	defaultLog.Trace("hosttrust/verifier:validateCachedFlavors() Entering")
	defer defaultLog.Trace("hosttrust/verifier:validateCachedFlavors() Leaving")

	htc := hostTrustCache{
		hostID: hostId,
	}
	var collectiveReport hvs.TrustReport
	var trustCachesToDelete []uuid.UUID
	for _, cachedFlavor := range cachedFlavors {
		//TODO: change the signature verification depending on decision on signed flavors
		report, err := v.flavorVerifier.Verify(hostData, &cachedFlavor, true)
		if err != nil {
			return hostTrustCache{}, errors.Wrap(err, "hosttrust/verifier:validateCachedFlavors() Error from flavor verifier")
		}
		if report.Trusted {
			htc.trustedFlavors = append(htc.trustedFlavors, cachedFlavor.Flavor)
			collectiveReport.Results = append(collectiveReport.Results, report.Results...)
		} else {
			trustCachesToDelete = append(trustCachesToDelete, cachedFlavor.Flavor.Meta.ID)
		}
	}
	// remove cache entries for flavors that could not be verified
	_ = v.hostStore.RemoveTrustCacheFlavors(hostId, trustCachesToDelete)
	htc.trustReport = collectiveReport
	return htc, nil
}

func (v *verifier) storeTrustReport(hostID uuid.UUID, trustReport *hvs.TrustReport, samlReport *saml.SamlAssertion) {
	defaultLog.Trace("hosttrust/verifier:storeTrustReport() Entering")
	defer defaultLog.Trace("hosttrust/verifier:storeTrustReport() Leaving")

	log.Debugf("hosttrust/verifier:storeTrustReport() flavorverify host: %s SAML Report: %s", hostID, samlReport.Assertion)
	hvsReport := models.HVSReport{
		HostID:      hostID,
		TrustReport: *trustReport,
		CreatedAt:   samlReport.CreatedTime,
		Expiration:  samlReport.ExpiryTime,
		Saml:        samlReport.Assertion,
	}
	_, err := v.reportStore.Create(&hvsReport)
	if err != nil {
		log.WithError(err).Errorf("hosttrust/verifier:storeTrustReport() Failed to store Report")
	}
}