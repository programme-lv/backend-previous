//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var TaskOrigins = newTaskOriginsTable("public", "task_origins", "")

type taskOriginsTable struct {
	postgres.Table

	// Columns
	Abbreviation postgres.ColumnString
	FullName     postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type TaskOriginsTable struct {
	taskOriginsTable

	EXCLUDED taskOriginsTable
}

// AS creates new TaskOriginsTable with assigned alias
func (a TaskOriginsTable) AS(alias string) *TaskOriginsTable {
	return newTaskOriginsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new TaskOriginsTable with assigned schema name
func (a TaskOriginsTable) FromSchema(schemaName string) *TaskOriginsTable {
	return newTaskOriginsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new TaskOriginsTable with assigned table prefix
func (a TaskOriginsTable) WithPrefix(prefix string) *TaskOriginsTable {
	return newTaskOriginsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new TaskOriginsTable with assigned table suffix
func (a TaskOriginsTable) WithSuffix(suffix string) *TaskOriginsTable {
	return newTaskOriginsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newTaskOriginsTable(schemaName, tableName, alias string) *TaskOriginsTable {
	return &TaskOriginsTable{
		taskOriginsTable: newTaskOriginsTableImpl(schemaName, tableName, alias),
		EXCLUDED:         newTaskOriginsTableImpl("", "excluded", ""),
	}
}

func newTaskOriginsTableImpl(schemaName, tableName, alias string) taskOriginsTable {
	var (
		AbbreviationColumn = postgres.StringColumn("abbreviation")
		FullNameColumn     = postgres.StringColumn("full_name")
		allColumns         = postgres.ColumnList{AbbreviationColumn, FullNameColumn}
		mutableColumns     = postgres.ColumnList{FullNameColumn}
	)

	return taskOriginsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Abbreviation: AbbreviationColumn,
		FullName:     FullNameColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}