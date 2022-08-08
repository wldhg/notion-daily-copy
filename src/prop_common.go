package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/wldh-g/notionapi"
)

func check_is_prop_true(prop notionapi.Property) bool {
	switch itemReal := prop.(type) {
	case *notionapi.CheckboxProperty:
		return itemReal.Checkbox
	case *notionapi.FormulaProperty:
		if itemReal.Formula.Type != "boolean" {
			return false
		}
		return itemReal.Formula.Boolean
	default:
		panic(fmt.Errorf("checking true/false of a property of type %T is not supported", itemReal))
	}
}

func props_to_texts(props []notionapi.Property) []string {
	var texts []string
	for _, prop := range props {
		if prop.GetType() == "pagenated" {
			prop = prop.(*notionapi.PropertyItem).Results[0]
		}
		switch itemReal := prop.(type) {
		case *notionapi.TitleProperty:
			var str string = ""
			for _, richText := range itemReal.Title {
				str += richText.Text.Content
			}
			texts = append(texts, str)
		case *notionapi.SingleTitleProperty:
			texts = append(texts, itemReal.Title.Text.Content)
		case *notionapi.RichTextProperty:
			var str string = ""
			for _, richText := range itemReal.RichText {
				str += richText.Text.Content
			}
			texts = append(texts, str)
		case *notionapi.TextProperty:
			var str string = ""
			for _, richText := range itemReal.Text {
				str += richText.Text.Content
			}
			texts = append(texts, str)
		default:
			panic(fmt.Errorf("converting a property of type %T to text is not supported", itemReal))
		}
	}
	return texts
}

func add_date_to_page_title(page *notionapi.Page) *notionapi.Page {
	tStr := get_time().Format(addDateToTitleFormat)
	if dateProp, ok := page.Properties[addDateToTitleProp]; ok {
		switch itemReal := dateProp.(type) {
		case *notionapi.TitleProperty:
			itemReal.Title = append(itemReal.Title, notionapi.RichText{
				Text: notionapi.Text{
					Content: tStr,
				},
			})
			page.Properties[addDateProp] = itemReal
		case *notionapi.SingleTitleProperty:
			itemReal.Title.Text.Content = tStr
		default:
			panic(fmt.Errorf("adding a date to a page title with a property of type %T is not supported", itemReal))
		}
	} else {
		panic(fmt.Errorf("the created page does not have a %s property", addDateToTitleProp))
	}
	return page
}

func add_date_to_page(page *notionapi.Page) *notionapi.Page {
	date := notionapi.Date(get_0am_time())
	page.Properties[addDateProp] = &notionapi.DateProperty{
		Date: &notionapi.DateObject{
			Start: &date,
			End:   &date,
		},
	}
	return page
}

func copy_time_to_page(page *notionapi.Page, startTTime string, endTTime string, dateOffset int) *notionapi.Page {
	tStr := get_time().AddDate(0, 0, dateOffset).Format(time.RFC3339)
	startTStr := strings.Replace(tStr, tStr[10:16], "T"+startTTime, 1)
	endTStr := strings.Replace(tStr, tStr[10:16], "T"+endTTime, 1)
	startT, startTStrGenErr := time.Parse(time.RFC3339, startTStr)
	if startTStrGenErr != nil {
		panic(startTStrGenErr)
	}
	startTDate := notionapi.Date(startT)
	endT, endTStrGenErr := time.Parse(time.RFC3339, endTStr)
	if endTStrGenErr != nil {
		panic(endTStrGenErr)
	}
	endTDate := notionapi.Date(endT)
	if page.Properties == nil {
		page.Properties = make(map[string]notionapi.Property)
	}
	page.Properties[copyTimeTargetProp] = &notionapi.DateProperty{
		Date: &notionapi.DateObject{
			Start: &startTDate,
			End:   &endTDate,
		},
	}
	return page
}

func add_tags_to_page(page *notionapi.Page) *notionapi.Page {
	if addTag {
		if page.Properties == nil {
			page.Properties = make(map[string]notionapi.Property)
		}
		page.Properties[addTagProp] = notionapi.MultiSelectProperty{
			MultiSelect: []notionapi.Option{
				{
					Name: addTagName,
				},
			},
		}
	}
	return page
}

func add_status_to_page(page *notionapi.Page) *notionapi.Page {
	if addStatus {
		if page.Properties == nil {
			page.Properties = make(map[string]notionapi.Property)
		}
		page.Properties[addStatusProp] = notionapi.StatusProperty{
			Status: notionapi.Option{
				Name: addStatusValue,
			},
		}
	}
	return page
}
