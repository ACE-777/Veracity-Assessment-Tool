package main

import (
	"flag"
	metrics "toloka-metrics/internal/metrics"
	toloka "toloka-metrics/internal/toloka"
	wiki "toloka-metrics/internal/wiki"
)

func main() {
	scrapeWiki := flag.Bool("scrape_wiki", false, "scrape data from wiki articles, that located in config file")
	flag.Parse()

	// parse labeled data from toloka or other source. Then aggregate it and make unique data for each sentence.
	res := toloka.NewResponseData()

	if *scrapeWiki {
		//scrape wiki articles that are in config file
		wiki.ScrapeDataFromWikiArticles()
	}

	// calculate output of python algorithm with snapshotting each step of data
	percentageColoredInEachSentence, mainLabelInEachSentence := metrics.GetColored(res)

	// calculate and prepare AUC data for future ROC curve
	metrics.GetAUC(percentageColoredInEachSentence, mainLabelInEachSentence)

}
