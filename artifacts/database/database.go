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

	. "git.tizen.org/tools/weles"

	"github.com/go-gorp/gorp"
	// sqlite3 is imported for side-effects and will be used
	// with the standard library sql interface.
	_ "github.com/mattn/go-sqlite3"
)

type artifactInfoRecord struct {
	ID int64 `db:",primarykey, autoincrement"`
	ArtifactInfo
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
func (aDB *ArtifactDB) InsertArtifactInfo(ai *ArtifactInfo) error {
	ar := artifactInfoRecord{
		ArtifactInfo: *ai,
	}
	return aDB.dbmap.Insert(&ar)
}

// SelectPath selects artifact from database based on its path.
func (aDB *ArtifactDB) SelectPath(path ArtifactPath) (ArtifactInfo, error) {
	ar := artifactInfoRecord{}
	err := aDB.dbmap.SelectOne(&ar, "select * from artifacts where Path=?", path)
	if err != nil {
		return ArtifactInfo{}, err
	}
	return ar.ArtifactInfo, nil
}

// Select fetches artifacts from ArtifactDB.
func (aDB *ArtifactDB) Select(arg interface{}) (artifacts []ArtifactInfo, err error) {
	var (
		results []artifactInfoRecord
		query   string
	)
	// TODO prepare efficient way of executing generic select.
	switch arg.(type) {
	case JobID:
		query = "select * from artifacts where JobID = ?"
	default:
		return nil, ErrUnsupportedQueryType
	}

	_, err = aDB.dbmap.Select(&results, query, arg)
	if err != nil {
		return nil, err
	}
	artifacts = make([]ArtifactInfo, len(results))
	for i, res := range results {
		artifacts[i] = res.ArtifactInfo
	}
	return artifacts, nil
}
