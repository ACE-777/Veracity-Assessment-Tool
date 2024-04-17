package toloka

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	True      = "True"
	False     = "False"
	CannotSay = "Cannot say"
	Nonsense  = "Nonsense"
	NoData    = "NoData"

	resultFile = "internal/toloka/results_from_pool_06-12-2023.tsv"
)

type ResultData struct {
	file     string
	question uint
	answer   uint
	sentence uint
}

type ResultLabelsTable struct {
	true      uint
	false     uint
	cannotSay uint
	nonsense  uint
}

// label win that takes part more than another labels, in bad case noData return for sentence
func (r *ResultLabelsTable) resultLabel() string {
	total := r.true + r.false + r.cannotSay + r.nonsense
	maxCount, maxLabel := r.maxLabel()
	if maxCount > total/2 {
		return maxLabel
	}

	return NoData
}

func (r *ResultLabelsTable) maxLabel() (maxCount uint, maxLabel string) {
	maxLabel = True
	maxCount = r.true

	if r.false > maxCount {
		maxCount = r.false
		maxLabel = False
	}

	if r.cannotSay > maxCount {
		maxCount = r.cannotSay
		maxLabel = NoData
	}

	if r.nonsense > maxCount {
		maxCount = r.nonsense
		maxLabel = NoData
	}

	return maxCount, maxLabel
}

func newResultData() map[ResultData]string {
	f, err := os.Open(resultFile)
	if err != nil {
		log.Printf("error while opening file: %v", err)
	}

	var responseDataFromResults = make(map[ResultData]ResultLabelsTable)
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

		currentRowFromResult := ResultData{
			file:     row[3],
			question: getIntegerFromFile(row[5]),
			answer:   getIntegerFromFile(row[4]),
			sentence: getIntegerFromFile(row[6]),
		}
		incrementResultLabelsTable(responseDataFromResults, currentRowFromResult, row[7])
	}

	log.Printf("complete readind result file")

	return definitionOfMaxLabel(responseDataFromResults)
}

func getIntegerFromFile(data string) uint {
	number, err := strconv.Atoi(data)
	if err != nil {
		log.Printf("error in converting string to int from file: %v", err)
	}

	return uint(number)
}

func incrementResultLabelsTable(
	responseDataFromResults map[ResultData]ResultLabelsTable,
	currentRowFromResult ResultData,
	row string,
) {
	if row == True {
		responseDataFromResults[currentRowFromResult] = ResultLabelsTable{
			true:      responseDataFromResults[currentRowFromResult].true + 1,
			false:     responseDataFromResults[currentRowFromResult].false,
			cannotSay: responseDataFromResults[currentRowFromResult].cannotSay,
			nonsense:  responseDataFromResults[currentRowFromResult].nonsense,
		}
	}

	if row == False {
		responseDataFromResults[currentRowFromResult] = ResultLabelsTable{
			true:      responseDataFromResults[currentRowFromResult].true,
			false:     responseDataFromResults[currentRowFromResult].false + 1,
			cannotSay: responseDataFromResults[currentRowFromResult].cannotSay,
			nonsense:  responseDataFromResults[currentRowFromResult].nonsense,
		}
	}

	if row == CannotSay {
		responseDataFromResults[currentRowFromResult] = ResultLabelsTable{
			true:      responseDataFromResults[currentRowFromResult].true,
			false:     responseDataFromResults[currentRowFromResult].false,
			cannotSay: responseDataFromResults[currentRowFromResult].cannotSay + 1,
			nonsense:  responseDataFromResults[currentRowFromResult].nonsense,
		}
	}

	if row == Nonsense {
		responseDataFromResults[currentRowFromResult] = ResultLabelsTable{
			true:      responseDataFromResults[currentRowFromResult].true,
			false:     responseDataFromResults[currentRowFromResult].false,
			cannotSay: responseDataFromResults[currentRowFromResult].cannotSay,
			nonsense:  responseDataFromResults[currentRowFromResult].nonsense + 1,
		}
	}
}

func definitionOfMaxLabel(responseDataFromResults map[ResultData]ResultLabelsTable) map[ResultData]string {
	var responseResultDataFromResults = make(map[ResultData]string, len(responseDataFromResults))
	for resultData, labels := range responseDataFromResults {
		responseResultDataFromResults[resultData] = labels.resultLabel()
	}

	return responseResultDataFromResults
}
