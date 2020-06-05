// Copyright 2020 Liquidata, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package enginetest

import "github.com/liquidata-inc/go-mysql-server/sql"

// Harness provides a way for database integrators to validate their implementation against the standard set of queries
// used to develop and test the engine itself. See memory_engine_test.go for an example.
type Harness interface {
	// Parallelism returns how many parallel go routines to use when constructing an engine for test.
	Parallelism() int
	// NewDatabase returns a new sql.Database to use for a test.
	NewDatabase(name string) sql.Database
	// NewTable takes a database previously created by NewDatabase and returns a table created with the given schema.
	NewTable(db sql.Database, name string, schema sql.Schema) (sql.Table, error)
	// NewContext allows a harness to specify any sessions or context variables necessary for the proper functioning of
	// their engine implementation. Every harnessed engine test uses the context created by this method, with some
	// additional information (e.g. current DB) set uniformly. To replicated the behavior of tests during setup,
	// harnesses should generally dispatch to enginetest.NewContext(harness), rather than calling this method themselves.
	NewContext() *sql.Context
}

// SkippingHarness provides a way for integrators to skip tests that are known to be broken. E.g., integrators that
// can't handle every possible SQL type.
type SkippingHarness interface {
	// SkipQueryTest returns whether to skip a test of the provided query string.
	SkipQueryTest(query string) bool
}

// IndexDriverHarness is an extension to Harness that lets an integrator test their implementation alongside an index
// driver they provide.
type IndexDriverHarness interface {
	Harness
	// IndexDriver returns an index driver for the databases given, which will have been created by calls to
	// NewDatabase().
	IndexDriver(dbs []sql.Database) sql.IndexDriver
}

// IndexHarness is an extension to Harness that lets an integrator test their implementation with native
// (table-supplied) indexes. Integrator tables must implement sql.IndexAlterableTable.
type IndexHarness interface {
	Harness
	// SupportsNativeIndexCreation returns whether this harness should accept CREATE INDEX statements as part of test
	// setup.
	SupportsNativeIndexCreation() bool
}

// VersionedDBHarness is an extension to Harness that lets an integrator test their implementation of versioned (AS OF)
// queries. Integrators must implement sql.VersionedDatabase. For each table version being created, there will be a
// call to NewTableAsOf, some number of Delete and Insert operations, and then a call to SnapshotTable.
type VersionedDBHarness interface {
	Harness
	// NewTableAsOf creates a new table with the given name and schema, optionally handling snapshotting with the asOf
	// identifier. NewTableAsOf must ignore tables that already exist in the database. Tables returned by this method do
	// not need to have any previously created data in them, but they can. This behavior is implementation specific, and
	// the harness works either way.
	NewTableAsOf(db sql.VersionedDatabase, name string, schema sql.Schema, asOf interface{}) sql.Table
	// SnapshotTable creates a snapshot of the table named with the given asOf label. Depending on the implementation,
	// NewTableAsOf might do all the necessary work to create such snapshots, so this could be a no-op.
	SnapshotTable(db sql.VersionedDatabase, name string, asOf interface{}) error
}
