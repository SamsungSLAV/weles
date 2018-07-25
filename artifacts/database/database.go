/*
 *  Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
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
	"errors"
	"strings"

	"git.tizen.org/tools/weles"

	"github.com/go-gorp/gorp"
	// sqlite3 is imported for side-effects and will be used
	// with the standard library sql interface.
	_ "github.com/mattn/go-sqlite3"
)

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
	aDB.dbmap.AddTableWithName(weles.ArtifactInfo{}, "artifacts").SetKeys(true, "ID")

	return aDB.dbmap.CreateTablesIfNotExists()
}

// Close closes the database.
func (aDB *ArtifactDB) Close() error {
	return aDB.handler.Close()
}

// InsertArtifactInfo inserts information about artifact to database.
func (aDB *ArtifactDB) InsertArtifactInfo(ai *weles.ArtifactInfo) error {
	return aDB.dbmap.Insert(ai)
}

// SelectPath selects artifact from database based on its path.
func (aDB *ArtifactDB) SelectPath(path weles.ArtifactPath) (weles.ArtifactInfo, error) {
	//	ar := artifactInfoRecord{}
	ai := weles.ArtifactInfo{}
	err := aDB.dbmap.SelectOne(&ai, "select * from artifacts where Path=?", path)
	if err != nil {
		return weles.ArtifactInfo{}, err
	}
	return ai, nil
}

// prepareQuery prepares query based on given filter.
// TODO code duplication
func prepareQuery(
	filter weles.ArtifactFilter,
	sorter weles.ArtifactSorter,
	paginator weles.ArtifactPagination,
	getTotal, getRemaining bool, offset int64) (string, []interface{}) {
	var (
		conditions []string
		query      string
		args       []interface{}
	)
	if getTotal == false && getRemaining == false {
		query = "select * from artifacts "
	} else {
		query = "select count(*) from artifacts "
	}

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
	if getTotal == false && paginator.ID != 0 {
		if (paginator.Forward == true && sorter.SortOrder == weles.SortOrderDescending) || (paginator.Forward == false && (sorter.SortOrder == weles.SortOrderAscending || sorter.SortOrder == "")) {
			conditions = append(conditions, " ID < ? ")
			args = append(args, paginator.ID)
		} else {
			conditions = append(conditions, " ID > ? ")
			args = append(args, paginator.ID)
		}
	}

	if len(conditions) > 0 {
		query += " where " + strings.Join(conditions, " AND ")
	}
	//TODO: make timestamp also db key, add to where clause and order by as described in:
	// https://www.sqlite.org/rowvalue.html#scrolling_window_queries
	if sorter.SortOrder == weles.SortOrderDescending {
		query += " ORDER BY ID DESC "
	} else if sorter.SortOrder == weles.SortOrderAscending || sorter.SortOrder == "" {
		query += " ORDER BY ID ASC "
	}
	if paginator.Limit != 0 {
		if offset == 0 {
			query += " LIMIT ? "
			args = append(args, paginator.Limit)
		} else {
			query += " LIMIT ? OFFSET ?"
			args = append(args, paginator.Limit, offset)
		}
	}
	return query, args
}

// Filter fetches elements matching ArtifactFilter from database.
func (aDB *ArtifactDB) Filter(filter weles.ArtifactFilter, sorter weles.ArtifactSorter, paginator weles.ArtifactPagination) ([]weles.ArtifactInfo, weles.ListInfo, error) {
	results := []weles.ArtifactInfo{}
	var tr, rr int64
	// TODO gorp doesn't support passing list of arguments to where in(...) clause yet.
	// Thats why it's done with the use prepareQuery.
	trans, err := aDB.dbmap.Begin()
	if err != nil {
		return nil, weles.ListInfo{}, errors.New("Failed to open transaction while filtering " + err.Error())
	}
	queryForTotal, argsForTotal := prepareQuery(filter, sorter, paginator, true, false, 0)
	queryForRemaining, argsForRemaining := prepareQuery(filter, sorter, paginator, false, true, 0)
	var offset int64

	rr, err = aDB.dbmap.SelectInt(queryForRemaining, argsForRemaining...)
	if err != nil {
		return nil, weles.ListInfo{}, errors.New("Failed to get remaining records " + err.Error())
	}

	tr, err = aDB.dbmap.SelectInt(queryForTotal, argsForTotal...)
	if err != nil {
		return nil, weles.ListInfo{}, errors.New("Failed to get total records " + err.Error())
	}
	// TODO: refactor this file. below is to ignore pagination object when pagination is turned off.
	if paginator.Limit == 0 {
		paginator.Forward = true
		paginator.ID = 0
	}

	if paginator.Forward == false {
		offset = rr - int64(paginator.Limit)
	}

	queryForData, argsForData := prepareQuery(filter, sorter, paginator, false, false, offset)
	_, err = aDB.dbmap.Select(&results, queryForData, argsForData...)
	if err != nil {
		return nil, weles.ListInfo{}, err
	}
	if err := trans.Commit(); err != nil {
		return nil, weles.ListInfo{}, errors.New("Failed to commit transaction while filtering " + err.Error())

	}
	return results, weles.ListInfo{TotalRecords: uint64(tr), RemainingRecords: uint64(rr - int64(len(results)))}, nil
}

// Select fetches artifacts from ArtifactDB.
func (aDB *ArtifactDB) Select(arg interface{}) (artifacts []weles.ArtifactInfo, err error) {
	var (
		results []weles.ArtifactInfo
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
	query += " ORDER BY id"

	_, err = aDB.dbmap.Select(&results, query, arg)
	if err != nil {
		return nil, err
	}
	return results, nil
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

	id, err := aDB.getID(ai.Path)
	if err != nil {
		return err
	}
	ai.ID = id

	ai.Status = change.NewStatus
	_, err = aDB.dbmap.Update(&ai)
	return err
}
