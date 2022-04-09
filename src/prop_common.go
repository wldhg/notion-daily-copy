package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jomei/notionapi"
)

func props_to_texts(props []notionapi.Property) []string {
	var texts []string
	for _, prop := range props {
		switch itemReal := prop.(type) {
		case *notionapi.TitleProperty:
			var str string = ""
			for _, richText := range itemReal.Title {
				str += richText.Text.Content
			}
			texts = append(texts, str)
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
		case *notionapi.NumberProperty:

		default:
			panic(fmt.Errorf("converting a property of type %T to text is not supported", itemReal))
		}
	}
	return texts
}

func add_date_to_page_title(page *notionapi.Page) *notionapi.Page {
	tStr := get_time().Format(addDateFormat)
	if dateProp, ok := page.Properties[addDateProp]; ok {
		switch itemReal := dateProp.(type) {
		case *notionapi.TitleProperty:
			itemReal.Title = append(itemReal.Title, notionapi.RichText{
				Text: notionapi.Text{
					Content: tStr,
				},
			})
			page.Properties[addDateProp] = itemReal
		default:
			panic(fmt.Errorf("adding a date to a page title with a property of type %T is not supported", itemReal))
		}
	} else {
		panic(fmt.Errorf("the created page does not have a %s property", addDateProp))
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
		Date: notionapi.DateObject{
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
