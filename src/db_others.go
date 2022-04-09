package main

import "fmt"

func check_src_props() {
	for prop := range copyProps {
		if !check_prop_existence(srcDB, prop) {
			panic(fmt.Errorf("source property %s does not exist", prop))
		}
	}
	if copyTime && (!check_prop_existence(srcDB, copyTimeStartProp) || !check_prop_existence(srcDB, copyTimeEndProp)) {
		panic(fmt.Errorf("source properties %s and %s do not exist", copyTimeStartProp, copyTimeEndProp))
	}
}

func check_tgt_props() {
	for _, prop := range copyProps {
		if !check_prop_existence(tgtDB, prop) {
			panic(fmt.Errorf("target property %s does not exist", prop))
		}
	}
	if addDate && !check_prop_existence(tgtDB, addDateProp) {
		panic(fmt.Errorf("target property %s does not exist", addDateProp))
	}
	if addTag && !check_prop_existence(tgtDB, addTagProp) {
		panic(fmt.Errorf("target property %s does not exist", addTagProp))
	}
	if copyTime && !check_prop_existence(tgtDB, copyTimeTargetProp) {
		panic(fmt.Errorf("target property %s do not exist", copyTimeTargetProp))
	}
}
