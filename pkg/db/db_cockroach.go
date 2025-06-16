package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type CockroachImpl struct {
	sql *sql.DB
}

type crTableDescription struct {
	SchemaName        string
	TableName         string
	Type              string
	Owner             string
	EstimatedRowCount int64
	Locality          *string
}

type crRecordField struct {
	Name    string  `json:"field,omitempty"`
	Type    string  `json:"type,omitempty"`
	Null    bool    `json:"null,omitempty"`
	Default *string `json:"key,omitempty"`
	Ignore  *string
}

func (db *CockroachImpl) getTableFields(table string) ([]crRecordField, error) {
	list := make([]crRecordField, 0)

	rows, err := db.sql.Query(fmt.Sprintf("SHOW COLUMNS FROM %s;", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//var ignore string
	for rows.Next() {
		f := crRecordField{}
		err = rows.Scan(
			&f.Name,
			&f.Type,
			&f.Null,
			&f.Default,
			&f.Ignore,
			&f.Ignore,
			&f.Ignore,
		)
		if err != nil {
			break
		}
		list = append(list, f)
	}
	return list, err
}

type crTableConstraint struct {
	TableName      string
	ConstraintName string
	ConstraintType string
	Details        string
	Validated      bool
}

func (c *crTableConstraints) IsPrimaryKey(colName string) bool {
	for _, constraint := range *c {
		if constraint.ConstraintType == "PRIMARY KEY" {
			if strings.Contains(constraint.Details, colName) {
				return true
			}
		}
	}
	return false
}

type crTableConstraints []crTableConstraint

func (db *CockroachImpl) getTableConstraints(table string) (crTableConstraints, error) {
	constraints := crTableConstraints{}

	rows, err := db.sql.Query(fmt.Sprintf("SHOW CONSTRAINTS FROM %s;", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		f := crTableConstraint{}
		err = rows.Scan(
			&f.TableName,
			&f.ConstraintName,
			&f.ConstraintType,
			&f.Details,
			&f.Validated,
		)
		if err != nil {
			break
		}
		constraints = append(constraints, f)
	}
	return constraints, err
}

func (db *CockroachImpl) Close() error {
	return db.sql.Close()
}

func (db *CockroachImpl) ListTables() ([]string, error) {
	list := make([]string, 0)

	rows, err := db.sql.Query("SHOW TABLES;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//var ignore string
	for rows.Next() {
		f := crTableDescription{}
		err = rows.Scan(
			&f.SchemaName,
			&f.TableName,
			&f.Type,
			&f.Owner,
			&f.EstimatedRowCount,
			&f.Locality,
		)
		if err != nil {
			break
		}
		list = append(list, f.TableName)
	}
	return list, err
}

func (db *CockroachImpl) DescribeTable(name string) (Table, error) {
	table := Table{
		Name: name,
	}

	columnDefinitions, err := db.getTableFields(name)
	if err != nil {
		return table, err
	}

	constraints, err := db.getTableConstraints(name)
	if err != nil {
		return table, err
	}

	fmt.Printf("constraints: %+v\n", constraints)

	for _, definition := range columnDefinitions {
		//fmt.Printf("%s:%v\n", definition.Name, constraints.IsPrimaryKey(definition.Name))

		column := &Column{
			Name: definition.Name,
			Type: definition.Type,
			IsPK: constraints.IsPrimaryKey(definition.Name),
		}
		table.Columns = append(table.Columns, *column)
	}

	return table, nil
}

func NewCockroachImpl(dbURN string) (*CockroachImpl, error) {
	if dbURN == "" {
		return nil, fmt.Errorf("no database URN provided")
	}

	sqlDB, err := sql.Open("postgres", dbURN)
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return &CockroachImpl{
		sql: sqlDB,
	}, nil
}
