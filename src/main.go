package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/wldh-g/notionapi"
)

var spinnerType = spinner.CharSets[14]
var spinnerTime = 100 * time.Millisecond

var cli *notionapi.Client

var srcDB *notionapi.Database
var tgtDB *notionapi.Database

func main() {
	step0 := spinner.New(spinnerType, spinnerTime)
	step0.Suffix = " Checking environment variables..."
	step0.FinalMSG = "✔️ Checked environment variables.\n"
	step0.Start()
	read_dotenv()
	verify_env_var()
	step0.Stop()

	step1 := spinner.New(spinnerType, spinnerTime)
	step1.Suffix = " Preparing notion API client..."
	step1.FinalMSG = "✔️ Prepared notion API client.\n"
	step1.Start()
	{
		cli = notionapi.NewClient(notionapi.Token(secret))
		srcDB = get_db(srcID)
		tgtDB = get_db(tgtID)
	}
	step1.Stop()

	step2 := spinner.New(spinnerType, spinnerTime)
	step2.Suffix = " Checking items at source database..."
	step2.FinalMSG = "✔️ Checked items at source database.\n"
	step2.Start()
	var srcPages []notionapi.Page
	var srcPageTitles []string
	{
		check_src_props()
		srcTitleProp := get_db_title_prop_name(srcDB)
		srcPages = get_all_pages(srcDB).Results
		srcPageTitles = get_pages_prop_text_slice(srcPages, srcTitleProp)
	}
	step2.Stop()
	fmt.Printf("  %d items found in source database.\n", len(srcPageTitles))
	for _, item := range srcPageTitles {
		fmt.Printf("   - %s\n", item)
	}

	step3 := spinner.New(spinnerType, spinnerTime)
	step3.Suffix = " Adding items to target database..."
	step3.FinalMSG = "✔️ Added items to target database.\n"
	step3.Start()
	var tgtPages []notionapi.Page
	var tgtPageTitles []string
	{
		check_tgt_props()
		tgtTitleProp := get_db_title_prop_name(tgtDB)
		tgtPages = copy_pages(srcPages)
		add_pages(tgtPages)
		tgtPageTitles = get_pages_prop_text_slice(tgtPages, tgtTitleProp)
	}
	step3.Stop()
	fmt.Printf("  %d items added to target database.\n", len(tgtPageTitles))
	for _, item := range tgtPageTitles {
		fmt.Printf("   - %s\n", item)
	}

	fmt.Println("✅ Finished")
}
