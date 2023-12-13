package toloka

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type ResponseData struct {
	File     string
	Question uint
	Answer   uint
}

type Sentence struct {
	Sentence uint
	Text     string
	Label    string
}

func NewResponseData() map[ResponseData][]Sentence {
	f, err := os.Open("internal/toloka/tasks_from_pool_06-12-2023.tsv")
	if err != nil {
		log.Printf("error while opening file: %v", err)
	}

	resultData := newResultData()
	var result = make(map[ResponseData][]Sentence)
	r := csv.NewReader(f)
	r.Comma = '\t'
	for {
		var row []string
		row, err = r.Read()
		if err == io.EOF {
			break
		}

		if strings.Contains(row[5], "INPUT") {
			continue
		}

		currentRowFromTask := ResultData{
			file:     row[3],
			question: getIntegerFromFile(row[5]),
			answer:   getIntegerFromFile(row[4]),
			sentence: getIntegerFromFile(row[6]),
		}
		rowTransformForResultMap := ResponseData{
			File:     currentRowFromTask.file,
			Question: currentRowFromTask.question,
			Answer:   currentRowFromTask.answer,
		}

		label, ok := resultData[currentRowFromTask]
		if ok {
			result[rowTransformForResultMap] = append(result[rowTransformForResultMap], Sentence{
				Sentence: currentRowFromTask.sentence,
				Text:     row[1],
				Label:    label,
			})

			continue
		}

		result[rowTransformForResultMap] = append(result[rowTransformForResultMap], Sentence{
			Sentence: currentRowFromTask.sentence,
			Text:     row[1],
			Label:    NoData,
		})
	}

	log.Printf("complete reading task file")
	sortSentence(result)

	return result
}

func sortSentence(result map[ResponseData][]Sentence) {
	for _, sentences := range result {
		sort.Slice(sentences, func(i, j int) bool {
			return sentences[i].Sentence < sentences[j].Sentence
		})
	}
}
