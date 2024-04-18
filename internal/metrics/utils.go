package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"toloka-metrics/internal/toloka"
)

func saveMetaDataPerEachIteration(currentResultDir string, iterator int, coloredData ColoredData) {
	resultMetaInfoNameFile := fmt.Sprintf("%v\\%v",
		currentResultDir,
		fmt.Sprintf("meta_data_%v.%v", iterator, "json"),
	)
	MetaDataInDir, err := os.Create(resultMetaInfoNameFile)
	if err != nil {
		log.Printf("can not create dir for result data:%v", err)
	}
	metaInfo, err := json.MarshalIndent(coloredData, "", "  ")
	if err != nil {
		log.Printf("can not unmarshall meta info at iteration %v, err: %v", iterator, err)
	}

	_, err = io.Copy(MetaDataInDir, strings.NewReader(string(metaInfo)))
	if err != nil {
		log.Printf("can not write to result colored data file %v", err)
	}
}

func saveColoredDataPerEachIteration(currentResultDir string, iterator int, coloredData ColoredData) {
	resultColoredNameFile := fmt.Sprintf("%v\\%v",
		currentResultDir,
		fmt.Sprintf("colored_data_%v.%v", iterator, "html"),
	)
	ColoredResultInDir, err := os.Create(resultColoredNameFile)
	if err != nil {
		log.Printf("can not create dir for result data:%v", err)
	}

	coloredData.HTML = strings.ReplaceAll(coloredData.HTML, "coloredCount", coloredData.PercentageColored)

	_, err = io.Copy(ColoredResultInDir, strings.NewReader(coloredData.HTML))
	if err != nil {
		log.Printf("can not write to result colored data file %v", err)
	}
}

func buildSourcesForOutput(sentences []toloka.Sentence) (sources string) {
	for i, sentence := range sentences {
		sources += sentence.Sources
		if i != len(sentences)-1 {
			sources += ","
		}
	}

	return
}

func getSourceForTestyTasks(sentences []toloka.Sentence, coloredData ColoredData) {
	var source string
	for _, sentence := range sentences {
		source = sentence.Sources
		fmt.Println("source:::", source)

		positionVariants := make(map[int]int, len(coloredData.ResultSources[0]))

		for _, resultSources := range coloredData.ResultSources {
			for j, resultSource := range resultSources {
				if strings.Contains(resultSource, source) {
					positionVariants[j] += 1
				}
			}

		}

		for key, value := range positionVariants {
			fmt.Println("key:", key, "value:", value, "/", len(coloredData.Tokens), "res:", float64(value)/float64(len(coloredData.Tokens)))
		}

		fmt.Println("\n")
	}

}
