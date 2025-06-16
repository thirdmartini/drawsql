package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"drawsql/pkg/db"
	"drawsql/pkg/renderer"
)

var tableRelations = Relationships{
	"*": {
		"customer_id": "lotuscustomers",
		"device_id":   "fis_inventory",
		"policy_id":   "machinepolicies",
		"user_id":     "lotususers",
	},
	"fleetmachines": {
		"device_id": "fis_inventory",
	},
}

func main() {
	databaseFlag := flag.String("db", "", "DB URN")
	metadataFlag := flag.String("metadata", "schema.json", "Additional Schema Metadata")
	flag.Parse()

	if *databaseFlag == "" {
		flag.Usage()
		return
	}

	dbs, err := db.NewCockroachImpl(*databaseFlag)
	if err != nil {
		panic(err)
	}
	defer dbs.Close()

	tableRelations, err := loadMedataSchema(*metadataFlag)
	if err != nil {
		log.Printf("Error loading medata schema fomr %s: %v", *metadataFlag, err)
		return
	}

	tableNames, err := dbs.ListTables()
	if err != nil {
		panic(err)
	}

	tables := make([]db.Table, 0)

	for _, tableName := range tableNames {
		table, err := dbs.DescribeTable(tableName)
		if err != nil {
			log.Printf("DescribeTable(%v): %v", tableName, err)
			continue
		}

		for idx := range table.Columns {
			if refTable, ok := tableRelations.RefersTo(tableName, table.Columns[idx].Name); ok {
				table.Columns[idx].IsFK = true
				table.Columns[idx].ReferTo = append(table.Columns[idx].ReferTo, fmt.Sprintf("%s.%s", refTable, table.Columns[idx].Name))
			}
		}

		tables = append(tables, table)
	}

	group := db.Group{
		Name:   "phonehome",
		Label:  "phonehome",
		Tables: tables,
	}

	diagram, _, err := renderer.GenerateD2([]db.Group{group}, "down")
	if err != nil {
		panic(err)
	}

	renderResult, err := renderer.RenderSvg(diagram)
	if err != nil {
		panic(err)
	}

	os.WriteFile("graph.svg", renderResult, os.ModePerm)
}
