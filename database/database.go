// In-memory storage for application data
package database

import "github.com/hashicorp/go-memdb"

var DB *memdb.MemDB

func init() {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"presentations": {
				Name: "presentations",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
				},
			},
		},
	}
	memdb, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	DB = memdb
}
