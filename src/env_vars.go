package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var secret string
var srcID string
var tgtID string
var copyProps map[string]string
var dateOffset int64 = 0

var addDate bool
var addDateProp string
var addDateFormat string

var addTag bool
var addTagProp string
var addTagName string

var copyTime bool
var copyTimeTargetProp string
var copyTimeStartProp string
var copyTimeEndProp string

func read_dotenv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	secret = os.Getenv("INTEGRATION_SECRET")
	srcID = os.Getenv("SOURCE_DATABASE")
	tgtID = os.Getenv("TARGET_DATABASE")

	copyPropsStr := os.Getenv("COPY_PROPERTY")
	copyPropsSplit := strings.Split(copyPropsStr, ",")
	copyProps = make(map[string]string)
	for _, prop := range copyPropsSplit {
		copyPropsListPartial := strings.Split(prop, ">")
		if len(copyPropsListPartial) != 2 {
			panic(fmt.Errorf("invalid copy property %s", prop))
		}
		copyProps[copyPropsListPartial[0]] = copyPropsListPartial[1]
	}

	dateOffsetStr := os.Getenv("DATE_OFFSET")
	if len(dateOffsetStr) > 0 {
		dateOffset, err = strconv.ParseInt(dateOffsetStr, 10, 64)
		if err != nil {
			panic(err)
		}
	}

	addDate = os.Getenv("ADD_DATE") == "true"
	addDateProp = os.Getenv("TITLE_PROP_NAME")
	addDateFormat = os.Getenv("DATE_FORMAT")

	addTag = os.Getenv("ADD_TAG") == "true"
	addTagProp = os.Getenv("TAG_PROP_NAME")
	addTagName = os.Getenv("TAG_VALUE")

	copyTime = os.Getenv("COPY_TIME") == "true"
	copyTimeTargetProp = os.Getenv("TIME_PROP_NAME")
	copyTimeStartProp = os.Getenv("START_TIME_PROP_NAME")
	copyTimeEndProp = os.Getenv("END_TIME_PROP_NAME")
}

func verify_env_var() {
	if len(secret) != 50 || strings.Index(secret, "secret_") != 0 {
		panic(fmt.Errorf("invalid secret"))
	}

	if len(srcID) != 32 {
		panic(fmt.Errorf("invalid source database ID"))
	}

	if len(tgtID) != 32 {
		panic(fmt.Errorf("invalid target database ID"))
	}

	if len(copyProps) == 0 {
		panic(fmt.Errorf("no copy properties"))
	}
	for key, prop := range copyProps {
		if len(key) == 0 || len(prop) == 0 {
			panic(fmt.Errorf("invalid copy property"))
		}
	}

	if addDate {
		if len(addDateProp) == 0 {
			panic(fmt.Errorf("no date property name"))
		}
		if len(addDateFormat) == 0 {
			panic(fmt.Errorf("no date format"))
		}
	}

	if addTag {
		if len(addTagProp) == 0 {
			panic(fmt.Errorf("no tag property name"))
		}
		if len(addTagName) == 0 {
			panic(fmt.Errorf("no tag name"))
		}
	}

	if copyTime {
		if len(copyTimeTargetProp) == 0 {
			panic(fmt.Errorf("no target time property name"))
		}
		if len(copyTimeStartProp) == 0 {
			panic(fmt.Errorf("no start time property name"))
		}
		if len(copyTimeEndProp) == 0 {
			panic(fmt.Errorf("no end time property name"))
		}
	}
}
