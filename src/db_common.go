package main

import (
	"context"
	"fmt"

	"github.com/wldh-g/notionapi"
)

func get_db(id string) *notionapi.Database {
	db, err := cli.Database.Get(context.Background(), notionapi.DatabaseID(id))
	if err != nil {
		panic(err)
	}
	return db
}

func get_db_title_prop_name(db *notionapi.Database) string {
	// Get title value
	title := ""
	for name, obj := range db.Properties {
		if obj.GetType() == "title" {
			title = name
			break
		}
	}
	if title == "" {
		panic(fmt.Errorf("the database does not have a title property"))
	}
	return title
}
