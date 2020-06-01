/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package postgres

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/intel-secl/intel-secl/v3/pkg/model/hvs"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type FlavorGroupStore struct {
	Store *DataStore
}

func NewFlavorGroupStore(store *DataStore) *FlavorGroupStore {
	return &FlavorGroupStore{store}
}

func (f *FlavorGroupStore) Create(flavorGroup *hvs.FlavorGroup) (*hvs.FlavorGroup, error) {
	defaultLog.Trace("postgres/flavorgroup_store:Create() Entering")
	defer defaultLog.Trace("postgres/flavorgroup_store:Create() Leaving")

	dbFlavorGroup, err := toDbFlavorGroup(flavorGroup)
	if err != nil {
		return nil, errors.Wrap(err, "postgres/flavorgroup_store:Create() failed to marshal to dbFlavorgroup")
	}

	dbFlavorGroup.ID, err = uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "postgres/flavorgroup_store:Create() Error generating new UUID")
	}

	if err := f.Store.Db.Create(&dbFlavorGroup).Error; err != nil {
		return nil, errors.Wrap(err, "postgres/flavorgroup_store:Create() failed to create Flavorgroup")
	}
	flavorGroup, err = fromDbFlavorGroup(dbFlavorGroup)
	if err != nil {
		return nil, errors.Wrap(err, "postgres/flavorgroup_store:Create() failed to unmarshal to Flavorgroup")
	}
	return flavorGroup, nil
}

func (f *FlavorGroupStore) Retrieve(flavorGroupId *uuid.UUID) (*hvs.FlavorGroup, error) {
	defaultLog.Trace("postgres/flavorgroup_store:Retrieve() Entering")
	defer defaultLog.Trace("postgres/flavorgroup_store:Retrieve() Leaving")

	dbFlavorGroup := flavorGroup{
		ID: *flavorGroupId,
	}
	err := f.Store.Db.Where(&dbFlavorGroup).First(&dbFlavorGroup).Error
	if err != nil {
		return nil, errors.Wrap(err, "postgres/flavorgroup_store:Retrieve() failed to retrieve Flavorgroup")
	}
	return fromDbFlavorGroup(&dbFlavorGroup)
}

func (f *FlavorGroupStore) Search(fgFilter *hvs.FlavorGroupFilterCriteria) (*hvs.FlavorgroupCollection, error) {
	defaultLog.Trace("postgres/flavorgroup_store:Search() Entering")
	defer defaultLog.Trace("postgres/flavorgroup_store:Search() Leaving")

	tx := buildFlavorGroupSearchQuery(f.Store.Db, fgFilter)

	if tx == nil {
		return nil, errors.New("postgres/flavorgroup_store:Search() Unexpected Error. Could not build" +
			" a gorm query object in FlavorGroups Search function.")
	}

	var dbFlavorgroups []flavorGroup
	if err := tx.Find(&dbFlavorgroups).Error; err != nil {
		return nil, errors.Wrap(err, "postgres/flavorgroup_store:Search() failed to search all "+
			"Flavorgroups")
	}

	return fromDbFlavorGroups(dbFlavorgroups)
}

func (f *FlavorGroupStore) Delete(flavorGroupId *uuid.UUID) error {
	defaultLog.Trace("postgres/flavorgroup_store:Delete() Entering")
	defer defaultLog.Trace("postgres/flavorgroup_store:Delete() Leaving")

	dbFlavorGroup := flavorGroup{
		ID: *flavorGroupId,
	}
	if err := f.Store.Db.Delete(&dbFlavorGroup).Error; err != nil {
		return errors.Wrap(err, "postgres/flavorgroup_store:Delete() failed to delete Flavorgroup")
	}
	return nil
}

// helper function to build the query object for a FlavorGroup search.
func buildFlavorGroupSearchQuery(tx *gorm.DB, fgFilter *hvs.FlavorGroupFilterCriteria) *gorm.DB {
	defaultLog.Trace("postgres/flavorgroup_store:buildFlavorGroupSearchQuery() Entering")
	defer defaultLog.Trace("postgres/flavorgroup_store:buildFlavorGroupSearchQuery() Leaving")

	if tx == nil {
		return nil
	}

	if fgFilter == nil {
		return tx.Where(&flavorGroup{})
	}
	if fgFilter.Id != "" {
		tx = tx.Where("id = ?", fgFilter.Id)
	} else if fgFilter.NameEqualTo != "" {
		tx = tx.Where("name = ?", fgFilter.NameEqualTo)
	} else if fgFilter.NameContains != "" {
		tx = tx.Where("name like ? ", "%"+fgFilter.NameContains+"%")
	}
	//TODO: Add search for hostId
	return tx
}

func toDbFlavorGroup(fg *hvs.FlavorGroup) (*flavorGroup, error) {
	defaultLog.Trace("postgres/flavorgroup_store:toDbFlavorGroup() Entering")
	defer defaultLog.Trace("postgres/flavorgroup_store:toDbFlavorGroup() Leaving")
	if fg == nil {
		return nil, nil
	}

	flavorMatchPolicyCollection, err := json.Marshal(fg.FlavorMatchPolicyCollection)
	if err != nil {
		return nil, errors.Wrap(err, "postgres/flavorgroup_store:toDbFlavorGroup() failed to" +
			" marshal FlavorMatchPolicyCollection to JSON")
	}
	dbFlavorGroup := flavorGroup{
		ID:                    fg.ID,
		Name:                  fg.Name,
		FlavorTypeMatchPolicy: &postgres.Jsonb{RawMessage: flavorMatchPolicyCollection},
	}
	return &dbFlavorGroup, nil
}

func fromDbFlavorGroups(dbFlavorgroups []flavorGroup) (*hvs.FlavorgroupCollection, error) {
	defaultLog.Trace("postgres/flavorgroup_store:fromDbFlavorGroups() Entering")
	defer defaultLog.Trace("postgres/flavorgroup_store:fromDbFlavorGroups() Leaving")

	var flavorgroupCollection hvs.FlavorgroupCollection
	if dbFlavorgroups == nil || len(dbFlavorgroups) == 0 {
		flavorgroupCollection.Flavorgroups = []*hvs.FlavorGroup{}
		return &flavorgroupCollection, nil
	}

	for _, dbFlavorGroup := range dbFlavorgroups {
		flavorgroup, err := fromDbFlavorGroup(&dbFlavorGroup)
		if err != nil {
			return &flavorgroupCollection, errors.Wrap(err, "postgres/flavorgroup_store:fromDbFlavorGroups() " +
				"failed to unmarshal dbFlavorGroup")
		}
		flavorgroupCollection.Flavorgroups = append(flavorgroupCollection.Flavorgroups, flavorgroup)
	}

	return &flavorgroupCollection, nil
}

func fromDbFlavorGroup(fg *flavorGroup) (*hvs.FlavorGroup, error) {
	log.Trace("postgres/flavorgroup_store:fromDbFlavorGroup() Entering")
	defer log.Trace("postgres/flavorgroup_store:fromDbFlavorGroup() Leaving")

	if fg == nil {
		return nil, nil
	}

	var matchPolicyCollection hvs.FlavorMatchPolicyCollection
	err := json.Unmarshal(fg.FlavorTypeMatchPolicy.RawMessage, &matchPolicyCollection)
	if err != nil {
		return nil, errors.Wrap(err, "postgres/flavorgroup_store:fromDbFlavorGroup()" +
			" Error in unmarshalling the FlavorTypeMatchPolicy")
	}

	flavorGroup := hvs.FlavorGroup{
		ID:   fg.ID,
		Name: fg.Name,
	}

	if &matchPolicyCollection != nil && len(matchPolicyCollection.FlavorMatchPolicies) > 0 {
		flavorGroup.FlavorMatchPolicyCollection = &matchPolicyCollection
	}
	return &flavorGroup, nil
}