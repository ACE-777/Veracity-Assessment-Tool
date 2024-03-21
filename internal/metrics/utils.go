package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	_, err = io.Copy(ColoredResultInDir, strings.NewReader(coloredData.HTML))
	if err != nil {
		log.Printf("can not write to result colored data file %v", err)
	}
}
