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
				case *notionapi.SelectProperty:
					page.Properties[tgtPropName] = &notionapi.SelectProperty{
						Select: notionapi.Option{
							Name: propObjReal.Select.Name,
						},
					}
				case *notionapi.MultiSelectProperty:
					multiSelection := make([]notionapi.Option, 0, len(propObjReal.MultiSelect))
					for _, option := range propObjReal.MultiSelect {
						multiSelection = append(multiSelection, notionapi.Option{
							Name: option.Name,
						})
					}
					page.Properties[tgtPropName] = &notionapi.MultiSelectProperty{
						MultiSelect: multiSelection,
					}
				default:
					panic(fmt.Errorf("copying a property of type %T is not supported", propObjReal))
				}
			}
		}

		if addDate {
			page = add_date_to_page_title(page)
		}

		if addTag {
			page = add_tags_to_page(page)
		}

		if copyTime {
			var copyTimeStartPropObj, copyTimeEndPropObj, copyTimeDateOffsetPropObj notionapi.Property
			var copyTimePropChecked bool
			if copyTimeStartPropObj, copyTimePropChecked = srcPage.Properties[copyTimeStartProp]; !copyTimePropChecked {
				panic(fmt.Errorf("the page does not have a property named %s", copyTimeStartProp))
			}
			if copyTimeEndPropObj, copyTimePropChecked = srcPage.Properties[copyTimeEndProp]; !copyTimePropChecked {
				panic(fmt.Errorf("the page does not have a property named %s", copyTimeEndProp))
			}
			copyTimeDateOffsetPropObj, copyTimePropChecked = srcPage.Properties[copyTimeDateOffsetProp]
			if !copyTimePropChecked {
				copyTimeDateOffsetPropObj = notionapi.NumberProperty{
					Number: 0,
				}
			}
			copyTimeDateOffsetNumberPropObj := copyTimeDateOffsetPropObj.(*notionapi.NumberProperty)
			copyTimeDateOffset := int(copyTimeDateOffsetNumberPropObj.Number)

			timeStrings := props_to_texts([]notionapi.Property{copyTimeStartPropObj, copyTimeEndPropObj})
			if len(timeStrings) == 2 {
				if timeStrings[0] != "" && timeStrings[1] != "" {
					page = copy_time_to_page(page, timeStrings[0], timeStrings[1], copyTimeDateOffset)
				}
			} else {
				panic(fmt.Errorf("the page does not have a start and end time"))
			}
		}

		pages = append(pages, *page)
	}

	return pages
}
