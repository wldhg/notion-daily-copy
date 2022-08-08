package main

import (
	"fmt"
	"strings"

	"github.com/wldh-g/notionapi"
)

func check_src_props() {
	for prop := range copyProps {
		if !check_prop_existence(srcDB, prop) {
			availableProps := strings.Join(get_prop_list(srcDB)[:], ",")
			panic(fmt.Errorf("source property %s does not exist (not in: %s)", prop, availableProps))
		}
	}
	if filterEnabled {
		if !check_prop_existence(srcDB, filterBoolProp) {
			availableProps := strings.Join(get_prop_list(srcDB)[:], ",")
			panic(fmt.Errorf("source property %s does not exist (not in: %s)", filterBoolProp, availableProps))
		}
		if !check_prop_type(srcDB, filterBoolProp, []notionapi.PropertyConfigType{"checkbox", "formula"}) {
			panic(fmt.Errorf("source property %s is not a boolean type", filterBoolProp))
		}
	}
	if copyTime && (!check_prop_existence(srcDB, copyTimeStartProp) || !check_prop_existence(srcDB, copyTimeEndProp)) {
		availableProps := strings.Join(get_prop_list(srcDB)[:], ",")
		panic(fmt.Errorf("source properties %s and %s do not exist (not in: %s)", copyTimeStartProp, copyTimeEndProp, availableProps))
	}
}

func check_tgt_props() {
	for _, prop := range copyProps {
		if !check_prop_existence(tgtDB, prop) {
			availableProps := strings.Join(get_prop_list(tgtDB)[:], ",")
			panic(fmt.Errorf("target property %s does not exist (not in: %s)", prop, availableProps))
		}
	}
	if addDateToTitle && !check_prop_existence(tgtDB, addDateToTitleProp) {
		availableProps := strings.Join(get_prop_list(tgtDB)[:], ",")
		panic(fmt.Errorf("target property %s does not exist (not in: %s)", addDateToTitleProp, availableProps))
	}
	if addDate && !check_prop_existence(tgtDB, addDateProp) {
		availableProps := strings.Join(get_prop_list(tgtDB)[:], ",")
		panic(fmt.Errorf("target property %s does not exist (not in: %s)", addDateProp, availableProps))
	}
	if addTag && !check_prop_existence(tgtDB, addTagProp) {
		availableProps := strings.Join(get_prop_list(tgtDB)[:], ",")
		panic(fmt.Errorf("target property %s does not exist (not in: %s)", addTagProp, availableProps))
	}
	if addStatus && !check_prop_existence(tgtDB, addStatusProp) {
		availableProps := strings.Join(get_prop_list(tgtDB)[:], ",")
		panic(fmt.Errorf("target property %s does not exist (not in: %s)", addStatusProp, availableProps))
	}
	if copyTime && !check_prop_existence(tgtDB, copyTimeTargetProp) {
		availableProps := strings.Join(get_prop_list(tgtDB)[:], ",")
		panic(fmt.Errorf("target property %s does not exist (not in: %s)", copyTimeTargetProp, availableProps))
	}
}
