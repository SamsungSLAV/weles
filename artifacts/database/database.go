/*
 *  Copyright (c) 2017 Samsung Electronics Co., Ltd All Rights Reserved
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License
 */

// Package database is responsible for Weles system's job artifact storage.
package database

import (
	"database/sql"
	"strings"

	"git.tizen.org/tools/weles"

	"github.com/go-gorp/gorp"
	// sqlite3 is imported for side-effects and will be used
	// with the standard library sql interface.
	_ "github.com/mattn/go-sqlite3"
)

type artifactInfoRecord struct {
	ID int64 `db:",primarykey, autoincrement"`
	weles.ArtifactInfo
}

// ArtifactDB is responsible for database connection and queries.
type ArtifactDB struct {
	handler *sql.DB
	dbmap   *gorp.DbMap
}

// Open opens database connection.
func (aDB *ArtifactDB) Open(dbPath string) error {
	var err error
	aDB.handler, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	aDB.dbmap = &gorp.DbMap{Db: aDB.handler, Dialect: gorp.SqliteDialect{}}
	return aDB.initDB()
}

// initDB initializes tables.
func (aDB *ArtifactDB) initDB() error {
	// Add tables.
	aDB.dbmap.AddTableWithName(artifactInfoRecord{}, "artifacts").SetKeys(true, "ID")

	return aDB.dbmap.CreateTablesIfNotExists()
}

// Close closes the database.
func (aDB *ArtifactDB) Close() error {
	return aDB.handler.Close()
}

// InsertArtifactInfo inserts information about artifact to database.
func (aDB *ArtifactDB) InsertArtifactInfo(ai *weles.ArtifactInfo) error {
	ar := artifactInfoRecord{
		ArtifactInfo: *ai,
	}
	return aDB.dbmap.Insert(&ar)
}

// SelectPath selects artifact from database based on its path.
func (aDB *ArtifactDB) SelectPath(path weles.ArtifactPath) (weles.ArtifactInfo, error) {
	ar := artifactInfoRecord{}
	err := aDB.dbmap.SelectOne(&ar, "select * from artifacts where Path=?", path)
	if err != nil {
		return weles.ArtifactInfo{}, err
	}
	return ar.ArtifactInfo, nil
}

// prepareQuery prepares query based on given filter.
// TODO code duplication
func prepareQuery(filter weles.ArtifactFilter) (string, []interface{}) {
	var (
		conditions []string
		query      = "select * from artifacts "
		args       []interface{}
	)
	if len(filter.JobID) > 0 {
		q := make([]string, len(filter.JobID))
		for i, job := range filter.JobID {
			q[i] = "?"
			args = append(args, job)
		}
		conditions = append(conditions, " JobID in ("+strings.Join(q, ",")+")")
	}
	if len(filter.Type) > 0 {
		q := make([]string, len(filter.Type))
		for i, typ := range filter.Type {
			q[i] = "?"
			args = append(args, typ)
		}
		conditions = append(conditions, " Type in ("+strings.Join(q, ",")+")")
	}
	if len(filter.Status) > 0 {
		q := make([]string, len(filter.Status))
		for i, status := range filter.Status {
			q[i] = "?"
			args = append(args, status)
		}
		conditions = append(conditions, " Status in ("+strings.Join(q, ",")+")")
	}
	if len(filter.Alias) > 0 {
		q := make([]string, len(filter.Alias))
		for i, alias := range filter.Alias {
			q[i] = "?"
			args = append(args, alias)
		}
		conditions = append(conditions, " Alias in ("+strings.Join(q, ",")+")")
	}
	if len(conditions) > 0 {
		query += " where " + strings.Join(conditions, " AND ")
	}
	return query, args
}

// Filter fetches elements matching ArtifactFilter from database.
func (aDB *ArtifactDB) Filter(filter weles.ArtifactFilter) ([]weles.ArtifactInfo, error) {
	results := []artifactInfoRecord{}

	query, args := prepareQuery(filter)

	// TODO gorp doesn't support passing list of arguments to where in(...) clause yet.
	// Thats why it's done with the use prepareQuery.
	_, err := aDB.dbmap.Select(&results, query, args...)
	if err != nil {
		return nil, err
	}
	artifacts := make([]weles.ArtifactInfo, len(results))
	for i, res := range results {
		artifacts[i] = res.ArtifactInfo
	}
	return artifacts, nil

}

// Select fetches artifacts from ArtifactDB.
func (aDB *ArtifactDB) Select(arg interface{}) (artifacts []weles.ArtifactInfo, err error) {
	var (
		results []artifactInfoRecord
		query   string
	)
	// TODO prepare efficient way of executing generic select.
	switch arg.(type) {
	case weles.JobID:
		query = "select * from artifacts where JobID = ?"
	case weles.ArtifactType:
		query = "select * from artifacts where Type = ?"
	case weles.ArtifactAlias:
		query = "select * from artifacts where Alias = ?"
	case weles.ArtifactStatus:
		query = "select * from artifacts where Status = ?"
	default:
		return nil, ErrUnsupportedQueryType
	}

	_, err = aDB.dbmap.Select(&results, query, arg)
	if err != nil {
		return nil, err
	}
	artifacts = make([]weles.ArtifactInfo, len(results))
	for i, res := range results {
		artifacts[i] = res.ArtifactInfo
	}
	return artifacts, nil
}

// getID fetches ID of an artifact with provided path.
func (aDB *ArtifactDB) getID(path weles.ArtifactPath) (int64, error) {
	res, err := aDB.dbmap.SelectInt("select ID from artifacts where Path=?", path)
	if err != nil {
		return 0, err
	}
	return res, nil
}

// SetStatus changes artifact's status in ArtifactDB.
func (aDB *ArtifactDB) SetStatus(change weles.ArtifactStatusChange) error {
	ai, err := aDB.SelectPath(change.Path)
	if err != nil {
		return err
	}
	ar := artifactInfoRecord{
		ArtifactInfo: ai,
	}

	id, err := aDB.getID(ar.Path)
	if err != nil {
		return err
	}
	ar.ID = id

	ar.Status = change.NewStatus
	_, err = aDB.dbmap.Update(&ar)
	return err
}
