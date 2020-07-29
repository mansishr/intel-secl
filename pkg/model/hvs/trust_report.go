/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package hvs

//
// TrustReport model returned by verifier.Verify()
//

import (
	"github.com/google/uuid"
	constants "github.com/intel-secl/intel-secl/v3/pkg/hvs/constants/verifier-rules-and-faults"
	asset_tag "github.com/intel-secl/intel-secl/v3/pkg/lib/asset-tag"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/flavor/common"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/host-connector/types"
	ta "github.com/intel-secl/intel-secl/v3/pkg/model/ta"
	"reflect"
)

type TrustReport struct {
	PolicyName   string             `json:"policy_name"`
	Results      []RuleResult       `json:"results"`
	Trusted      bool               `json:"trusted"`
	HostManifest types.HostManifest `json:"host_manifest"`
}

type RuleResult struct {
	Rule     RuleInfo  `json:"rule"`
	FlavorId uuid.UUID `json:"flavor_id,omitempty"`
	Faults   []Fault   `json:"faults,omitempty"`
	Trusted  bool      `json:"trusted"`
}

type RuleInfo struct {
	Name                  string                     `json:"rule_name,omitempty"`
	Markers               []common.FlavorPart        `json:"markers,omitempty"`
	ExpectedPcr           *types.Pcr                 `json:"expected_pcr,omitempty"`
	FlavorID              *uuid.UUID                 `json:"flavor_id,omitempty"`
	FlavorName            *string                    `json:"flavor_name,omitempty"`
	ExpectedValue         *string                    `json:"expected_value,omitempty"`
	ExpectedMeasurements  []ta.FlavorMeasurement     `json:"expected_measurements,omitempty"`
	ExpectedEventLogs     []types.EventLog           `json:"expected,omitempty"`
	ExpectedEventLogEntry *types.EventLogEntry       `json:"expected,omitempty"`
	ExpectedTag           []byte                     `json:"expected_tag,omitempty"`
	Tags                  []asset_tag.TagKvAttribute `json:"tags,omitempty"`
}

type Fault struct {
	Name                   string                 `json:"fault_name"`
	Description            string                 `json:"description"`
	PcrIndex               *types.PcrIndex        `json:"pcr_index,omitempty"`
	ExpectedPcrValue       *string                `json:"expected_value,omitempty"`
	ActualPcrValue         *string                `json:"actual_value,omitempty"`
	MissingEntries         []types.EventLog       `json:"missing_entries,omitempty"`
	UnexpectedEntries      []types.EventLog       `json:"unexpected_entries,omitempty"`
	FlavorId               *uuid.UUID             `json:"flavor_id,omitempty"`
	UnexpectedMeasurements []ta.FlavorMeasurement `json:"unexpected_entries,omitempty"`
	MissingMeasurements    []ta.FlavorMeasurement `json:"missing_entries,omitempty"`
	MismatchMeasurements   []ta.FlavorMeasurement `json:"unexpected_entries,omitempty"`
	ExpectedValue          *string                `json:"expected_value,omitempty"`
	ActualValue            *string                `json:"actual_value,omitempty"`
	MeasurementId          *string                `json:"measurement_id,omitempty"`
	FlavorDigestAlg        *string                `json:"flavor_digest_alg,omitempty"`
	MeasurementDigestAlg   *string                `json:"measurement_digest_alg,omitempty"`
}

func NewTrustReport(report TrustReport) *TrustReport {
	return &TrustReport{PolicyName: report.PolicyName, Results: report.Results, Trusted: report.Trusted}
}

func (t *TrustReport) IsTrusted() bool {
	return t.isTrustedForResults(t.Results)
}

func (t *TrustReport) IsTrustedForMarker(marker string) bool {
	return t.isTrustedForResults(t.GetResultsForMarker(marker))
}
func (t *TrustReport) isTrustedForResults(ruleResults []RuleResult) bool {
	if len(ruleResults) == 0 {
		return false // empty policy is not trusted;  like RequireAllEmptySet fault.
	}

	trusted := true

	for _, result := range ruleResults {
		trusted = trusted && result.IsTrusted()
	}

	return trusted
}

func (t *TrustReport) GetResultsForMarker(marker string) []RuleResult {
	var ruleResults []RuleResult
	for _, result := range t.Results {
		markers := result.Rule.Markers
		if markers != nil {
			if find(markers, marker) {
				ruleResults = append(ruleResults, result)
			}
		}
	}
	return ruleResults
}

// TrustReport.java 106
func (t *TrustReport) CheckResultExists(targetRuleResult RuleResult) bool {
	//In TrustReport.java 107 marker := targetRuleResult.Rule.Markers[0] why dont we iterate all over the markers?
	marker := targetRuleResult.Rule.Markers[0]
	combinedRuleResults := t.GetResultsForMarker(marker.String())
	for _, ruleResult := range combinedRuleResults {
		if targetRuleResult.equals(ruleResult) {
			switch ruleResult.Rule.Name {
			case constants.RulePcrEventLogEqualsExcluding,
				constants.RulePcrEventLogIntegrity,
				constants.RulePcrMatchesConstant:
					// Compare pcrs for only these rules, all other rules will have pcr expected entry = nil
				if targetRuleResult.Rule.ExpectedPcr == nil {
					return false
				} else if reflect.DeepEqual(targetRuleResult.Rule.ExpectedPcr, ruleResult.Rule.ExpectedPcr){
					return true
				} else {
					continue
				}
			default:
				if len(targetRuleResult.Faults) > 0 {
					return false
				}
				return true
			}
		}
	}
	return false
}

func (t *TrustReport) AddResult(ruleResult RuleResult) {
	if !t.CheckResultExists(ruleResult) {
		t.Results = append(t.Results, ruleResult)
	}
}

func find(slice []common.FlavorPart, val string) bool {
	for _, item := range slice {
		if item.String() == val {
			return true
		}
	}
	return false
}

func (r *RuleResult) equals(target RuleResult) bool {
	if target.Rule.Name == r.Rule.Name && r.FlavorId == target.FlavorId {
		return true
	}
	return false
}

func (r *RuleResult) IsTrusted() bool {
	if r.Faults == nil || len(r.Faults) == 0 {
		return true
	}
	return false
}
