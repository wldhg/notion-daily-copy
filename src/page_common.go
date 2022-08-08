package main

import (
	"context"
	"fmt"

	"github.com/wldh-g/notionapi"
)

func check_prop_existence(db *notionapi.Database, prop string) bool {
	if _, ok := db.Properties[prop]; ok {
		return true
	}

	return false
}

func check_prop_type(db *notionapi.Database, prop string, t []notionapi.PropertyConfigType) bool {
	if p, ok := db.Properties[prop]; ok {
		for _, tt := range t {
			if p.GetType() == tt {
				return true
			}
		}
	}

	return false
}

func get_all_pages(db *notionapi.Database) *notionapi.DatabaseQueryResponse {
	items, err := cli.Database.Query(context.Background(), notionapi.DatabaseID(db.ID), &notionapi.DatabaseQueryRequest{
		PropertyFilter: nil,
		CompoundFilter: nil,
		PageSize:       100,
	})
	if err != nil {
		panic(err)
	}

	return items
}

func get_prop_list(db *notionapi.Database) []string {
	var props []string = make([]string, 0, len(db.Properties))
	for prop := range db.Properties {
		props = append(props, prop)
	}
	return props
}

func get_pages_prop_text_slice(items []notionapi.Page, prop string) []string {
	var properties []notionapi.Property = make([]notionapi.Property, 0, len(items))
	for _, item := range items {
		receivedProp, err := cli.Page.GetProperty(context.Background(), &item, prop)

		if err != nil {
			panic(err)
		}
		properties = append(properties, *receivedProp)
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
		if filterEnabled {
			receivedFilterProp, err := cli.Page.GetProperty(context.Background(), &srcPage, filterBoolProp)
			if err != nil {
				panic(err)
			}
			if !check_is_prop_true(*receivedFilterProp) {
				continue
			}
		}

		page := &notionapi.Page{}
		page.Parent.DatabaseID = notionapi.DatabaseID(tgtDB.ID)
		page.Properties = make(map[string]notionapi.Property)
		for srcPropName := range srcPage.Properties {
			propObj, err := cli.Page.GetProperty(context.Background(), &srcPage, srcPropName)
			if err != nil {
				panic(err)
			}
			if tgtPropName, ok := copyProps[srcPropName]; ok {
				propObjT := *propObj
				if propObjT.GetType() == "pagenated" {
					propObjT = propObjT.(*notionapi.PropertyItem).Results[0]
				}
				switch propObjReal := propObjT.(type) {
				case *notionapi.TitleProperty:
					page.Properties[tgtPropName] = &notionapi.TitleProperty{
						Title: propObjReal.Title,
					}
				case *notionapi.SingleTitleProperty:
					page.Properties[tgtPropName] = &notionapi.TitleProperty{
						Title: []notionapi.RichText{propObjReal.Title},
					}
				case *notionapi.RelationProperty:
					page.Properties[tgtPropName] = &notionapi.RelationProperty{
						Relation: propObjReal.Relation,
					}
				case *notionapi.SingleRelationProperty:
					page.Properties[tgtPropName] = &notionapi.RelationProperty{
						Relation: []notionapi.Relation{
							propObjReal.Relation,
						},
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
				case *notionapi.FormulaProperty:
					switch propObjReal.Formula.Type {
					case "string":
						richText := notionapi.RichText{
							Text: notionapi.Text{
								Content: propObjReal.Formula.String,
							},
						}
						page.Properties[tgtPropName] = &notionapi.RichTextProperty{
							RichText: []notionapi.RichText{richText},
						}
					default:
						panic(fmt.Sprintf("unknown formula type: %s", propObjReal.Formula.Type))
					}
				default:
					panic(fmt.Errorf("copying a property of type %T is not supported", propObjReal))
				}
			}
		}

		if addDateToTitle {
			page = add_date_to_page_title(page)
		}

		if addDate {
			page = add_date_to_page(page)
		}

		if addTag {
			page = add_tags_to_page(page)
		}

		if addStatus {
			page = add_status_to_page(page)
		}

		if copyTime {
			var copyTimeStartPropObj, copyTimeEndPropObj, copyTimeDateOffsetPropObjP *notionapi.Property
			var copyTimeDateOffsetPropObj notionapi.Property
			var copyTimePropRecvErr error
			if copyTimeStartPropObj, copyTimePropRecvErr = cli.Page.GetProperty(context.Background(), &srcPage, copyTimeStartProp); copyTimePropRecvErr != nil {
				panic(copyTimePropRecvErr)
			}
			if copyTimeEndPropObj, copyTimePropRecvErr = cli.Page.GetProperty(context.Background(), &srcPage, copyTimeEndProp); copyTimePropRecvErr != nil {
				panic(copyTimePropRecvErr)
			}
			copyTimeDateOffsetPropObjP, copyTimePropRecvErr = cli.Page.GetProperty(context.Background(), &srcPage, copyTimeDateOffsetProp)
			if copyTimePropRecvErr != nil {
				copyTimeDateOffsetPropObj = notionapi.NumberProperty{
					Number: 0,
				}
			} else {
				copyTimeDateOffsetPropObj = *copyTimeDateOffsetPropObjP
			}
			copyTimeDateOffsetNumberPropObj := copyTimeDateOffsetPropObj.(*notionapi.NumberProperty)
			copyTimeDateOffset := int(copyTimeDateOffsetNumberPropObj.Number)

			timeStrings := props_to_texts([]notionapi.Property{*copyTimeStartPropObj, *copyTimeEndPropObj})
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
