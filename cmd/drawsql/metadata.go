package main

import (
	"encoding/json"
	"os"
)

type TableRelation map[string]string

type Relationships map[string]TableRelation

func (r *Relationships) RefersTo(table string, column string) (string, bool) {
	tr, ok := (*r)[table]
	if ok {
		if ref, ok := tr[column]; ok {
			if ref == table {
				return "", false
			}

			return ref, true
		}
	}

	tr, ok = (*r)["*"]
	if !ok {
		return "", false
	}
	ref, ok := tr[column]
	if ref == table {
		return "", false
	}
	return ref, ok
}

func loadMedataSchema(schemaFile string) (*Relationships, error) {
	if schemaFile == "" {
		return new(Relationships), nil
	}

	data, err := os.ReadFile(schemaFile)
	if err != nil {
		return nil, err
	}

	r := new(Relationships)
	err = json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
