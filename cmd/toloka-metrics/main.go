package main

import (
	metrics "toloka-metrics/internal/metrics"
	toloka "toloka-metrics/internal/toloka"
)

func main() {
	// parse labeled data from toloka or other source. Then aggregate it and make unique data for each sentence.
	res := toloka.NewResponseData()

	// calculate output of python algorithm with snapshotting each step of data
	percentageColoredInEachSentence, mainLabelInEachSentence := metrics.GetColored(res)

	// calculate and prepare AUC data for future ROC curve
	metrics.GetAUC(percentageColoredInEachSentence, mainLabelInEachSentence)

}
