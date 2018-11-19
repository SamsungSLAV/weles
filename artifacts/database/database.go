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
	"log"
	"strings"

	"github.com/SamsungSLAV/weles"

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

const (
	sqlite3BusyTimeout = "?_busy_timeout=5000"
	sqlite3MaxOpenConn = 1
)

// Open opens database connection.
func (aDB *ArtifactDB) Open(dbPath string) error {
	var err error
	aDB.handler, err = sql.Open("sqlite3", dbPath+sqlite3BusyTimeout)
	if err != nil {
		return errors.New(dbOpenFail + err.Error())
	}
	aDB.handler.SetMaxOpenConns(sqlite3MaxOpenConn)

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
func (aDB *ArtifactDB) InsertArtifactInfo(ai *weles.ArtifactInfo) (err error) {
	err = aDB.dbmap.Insert(ai)
	if err != nil {
		log.Println("Failed to insert ArtifactInfo: ", err)
	}
	return
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
func prepareQuery(filter weles.ArtifactFilter, sorter weles.ArtifactSorter,
	paginator weles.ArtifactPagination, totalRecords, remainingRecords bool, offset int64,
) (query string, args []interface{}) {

	if !totalRecords && !remainingRecords {
		query = "select * from artifacts "
	} else {
		query = "select count(*) from artifacts "
	}

	var conditions []string
	conditions, args = prepareQueryFilter(filter)

	if !totalRecords && paginator.ID != 0 {
		if (paginator.Forward && sorter.Order == weles.SortOrderDescending) ||
			(!paginator.Forward && sorter.Order == weles.SortOrderAscending) {
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

	query += prepareQuerySorter(sorter)

	if paginator.Limit != 0 {
		if offset == 0 {
			query += " LIMIT ? "
			args = append(args, paginator.Limit)
		} else {
			query += " LIMIT ? OFFSET ?"
			args = append(args, paginator.Limit, offset)
		}
	}
	return
}

func prepareQuerySorter(sorter weles.ArtifactSorter) string {
	//TODO: make timestamp also db key, add to where clause and order by as described in:
	// https://www.sqlite.org/rowvalue.html#scrolling_window_queries
	if sorter.Order == weles.SortOrderDescending {
		return " ORDER BY ID DESC "
	}
	return " ORDER BY ID ASC "
}

func prepareQueryFilter(filter weles.ArtifactFilter) (conditions []string, args []interface{}) {
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

	return
}

// Filter fetches elements matching ArtifactFilter from database.
func (aDB *ArtifactDB) Filter(filter weles.ArtifactFilter, sorter weles.ArtifactSorter,
	paginator weles.ArtifactPagination) ([]weles.ArtifactInfo, weles.ListInfo, error) {

	results := []weles.ArtifactInfo{}
	var tr, rr int64
	// TODO gorp doesn't support passing list of arguments to where in(...) clause yet.
	// Thats why it's done with the use prepareQuery.
	trans, err := aDB.dbmap.Begin()
	if err != nil {
		return nil, weles.ListInfo{}, errors.New(whileFilter + dbTransOpenFail + err.Error())
	}
	defer func() {
		if err != nil {
			// err should be logged when it occurs.
			if err2 := trans.Rollback(); err2 != nil {
				log.Printf("%v occurred when filtering, trying to rollback transaction failed: %v",
					err, err2)
			}
		}
	}()
	queryForTotal, argsForTotal := prepareQuery(filter, sorter, paginator, true, false, 0)
	queryForRemaining, argsForRemaining := prepareQuery(filter, sorter, paginator, false, true, 0)
	var offset int64

	rr, err = trans.SelectInt(queryForRemaining, argsForRemaining...)
	if err != nil {
		return nil, weles.ListInfo{}, errors.New(whileFilter + dbRemainingFail + err.Error())
	}

	tr, err = trans.SelectInt(queryForTotal, argsForTotal...)
	if err != nil {
		return nil, weles.ListInfo{}, errors.New(whileFilter + dbTotalFail + err.Error())
	}

	if tr == 0 {
		// err needs to be updated for deferred 'if err!=nil' to catch it and roll back the
		// not committed transaction
		err = weles.ErrArtifactNotFound
		return []weles.ArtifactInfo{}, weles.ListInfo{}, err
	}

	if !paginator.Forward {
		offset = rr - int64(paginator.Limit)
	}

	queryForData, argsForData := prepareQuery(filter, sorter, paginator, false, false, offset)
	_, err = trans.Select(&results, queryForData, argsForData...)
	if err != nil {
		return nil, weles.ListInfo{}, errors.New(whileFilter + dbArtifactInfoFail + err.Error())
	}
	if err := trans.Commit(); err != nil {
		return nil, weles.ListInfo{}, errors.New(whileFilter + dbTransCommitFail + err.Error())

	}
	return results,
		weles.ListInfo{
			TotalRecords:     uint64(tr),
			RemainingRecords: uint64(rr - int64(len(results))),
		},
		nil
}

// SetStatus changes artifact's status in ArtifactDB.
func (aDB *ArtifactDB) SetStatus(change weles.ArtifactStatusChange) error {
	ai, err := aDB.SelectPath(change.Path)
	if err != nil {
		log.Println("failed to retrieve artifact based on its path: " + err.Error())
		return err //TODO: aalexanderr - log  error and continue
	}

	ai.Status = change.NewStatus
	if _, err = aDB.dbmap.Update(&ai); err != nil {
		log.Println("failed to update database" + err.Error())
		// TODO: aalexanderr - log critical, stop weles gracefully
	}
	return err
}
