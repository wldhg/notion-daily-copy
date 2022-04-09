package main

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
)

func check_prop_existence(db *notionapi.Database, prop string) bool {
	if _, ok := db.Properties[prop]; ok {
		return true
	}

	return false
}

func get_all_pages(db *notionapi.Database) *notionapi.DatabaseQueryResponse {
	items, err := cli.Database.Query(context.Background(), notionapi.DatabaseID(db.ID), &notionapi.DatabaseQueryRequest{
		PropertyFilter: nil,
		CompoundFilter: nil,
	})
	if err != nil {
		panic(err)
	}

	return items
}

func get_pages_prop_text_slice(items []notionapi.Page, prop string) []string {
	var properties []notionapi.Property = make([]notionapi.Property, 0, len(items))
	for _, item := range items {
		properties = append(properties, item.Properties[prop])
	}
	return props_to_texts(properties)
}

func add_pages(pages []notionapi.Page) {
	for _, page := range pages {
		_, err := cli.Page.Create(context.Background(), &notionapi.PageCreateRequest{
			Parent:     page.Parent,
			Properties: page.Properties,
		})
		if err != nil {
			panic(err)
		}
	}
}

func copy_pages(oPages []notionapi.Page) []notionapi.Page {
	var pages []notionapi.Page = make([]notionapi.Page, 0, len(oPages))

	for _, srcPage := range oPages {
		page := &notionapi.Page{}
		page.Parent.DatabaseID = notionapi.DatabaseID(tgtDB.ID)
		page.Properties = make(map[string]notionapi.Property)
		for srcPropName, propObj := range srcPage.Properties {
			if tgtPropName, ok := copyProps[srcPropName]; ok {
				switch propObjReal := propObj.(type) {
				case *notionapi.TitleProperty:
					page.Properties[tgtPropName] = &notionapi.TitleProperty{
						Title: propObjReal.Title,
					}
				case *notionapi.RelationProperty:
					page.Properties[tgtPropName] = &notionapi.RelationProperty{
						Relation: propObjReal.Relation,
					}
				default:
					panic(fmt.Errorf("copying a property of type %T is not supported", propObjReal))
				}
			}
		}

		page = add_date_to_page_title(page)
		page = add_tags_to_page(page)

		if copyTimeStartPropObj, ok := srcPage.Properties[copyTimeStartProp]; ok {
			if copyTimeEndPropObj, ok := srcPage.Properties[copyTimeEndProp]; ok {
				times := props_to_texts([]notionapi.Property{copyTimeStartPropObj, copyTimeEndPropObj})
				if len(times) == 2 {
					if times[0] != "" && times[1] != "" {
						page = copy_time_to_page(page, times[0], times[1])
					}
				} else {
					panic(fmt.Errorf("the page does not have a start and end time"))
				}
			}
		}
		pages = append(pages, *page)
	}

	return pages
}
