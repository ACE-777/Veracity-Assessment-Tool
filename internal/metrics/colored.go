package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	toloka "toloka-metrics/internal/toloka"
)

const (
	pythonScript = "test.color_build_data"
	repoDir      = "C:/Users/misha/pythonProject/chatgpt-research"

	TP = "TP"
	FP = "FP"
	TN = "TN"
	FN = "FN"

	flag = false

	uniqColor = 5
)

type ColoredData struct {
	File           string      `json:"file"`           //file name
	Question       string      `json:"question"`       //question number
	Answer         string      `json:"answer"`         //answer number
	ResultSources  [][]string  `json:"result_sources"` //sources with variants on each token
	Tokens         []string    `json:"tokens"`         //all tokens from input text
	ResultDistance [][]float64 `json:"result_dists"`   //cos dist result for each token with variants
	Labels         []string    //labels from algoritms for each sentence
	Length         []int       //count of all tokens for each sentence
	Colored        []int       //count of colored tokens for each sentence
	HTML           string      //html output for input
}

var attitudeMetric []float64
var labelsFromSentences []string

func GetColored(res map[toloka.ResponseData][]toloka.Sentence) ([]float64, []string) {
	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		iterator  = 0
		timestamp = time.Now().UTC().Add(time.Hour * 3).Format("2006-01-02T15-04-05")
	)

	sem := make(chan struct{}, 2)

	fileResultName := fmt.Sprintf("result_data_%v.%v", timestamp, "txt")
	fileResultData, err := os.Create(fileResultName)
	if err != nil {
		log.Printf("can not create file for result data:%v", err)
	}

	currentResultDir := fmt.Sprintf("out/colored_%v", timestamp)
	if err = os.Mkdir(currentResultDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	for k, v := range res {
		k := k
		v := v

		wg.Add(1)
		go func() {

			sem <- struct{}{}
			defer func() {
				wg.Done()
				<-sem
				fmt.Println("iterator::", iterator, "/", len(res))
				iterator++

			}()

			labels := make([]string, len(v))
			userInput := ""
			for sentenceNumber, sentence := range v {
				userInput += sentence.Text
				userInput += ". "
				labels[sentenceNumber] = sentence.Label
			}

			//fmt.Println("otput:", userInput)
			cmd := exec.Command(
				"python",
				"-m",
				pythonScript,
				"--userinput", userInput,
				"--file", k.File,
				"--question", strconv.Itoa(int(k.Question)),
				"--answer", strconv.Itoa(int(k.Answer)),
			)
			cmd.Dir = repoDir

			stdin, err := cmd.StdinPipe()
			if err != nil {
				log.Println("Can't execute python script")
				log.Println(err)
			}

			defer stdin.Close()

			var output bytes.Buffer
			cmd.Stdout = &output
			cmd.Stderr = os.Stderr
			if err = cmd.Start(); err != nil {
				log.Printf("error in starting python commnad: %v", err)
			}

			err = cmd.Wait()
			if err != nil {
				log.Println(err)
			}

			var coloredData ColoredData
			if err = json.Unmarshal(output.Bytes(), &coloredData); err != nil {
				log.Printf("Can't convert bytes to json struct coloredData %v", err)
			}

			coloredData.Labels = labels
			topLinksPerEachToken := getTopOneSource(coloredData)
			arrayWithTopUniqColors := buildDictForColor(topLinksPerEachToken, uniqColor)

			sentenceLenght, ColoredCount, HTML := buildPageTemplate(coloredData.Tokens, topLinksPerEachToken, arrayWithTopUniqColors)
			coloredData.Length = sentenceLenght
			coloredData.Colored = ColoredCount
			coloredData.HTML = HTML

			addAttitudeMetricAndWriteToSnapshotFileSafely(&mu, coloredData, fileResultData)

			saveColoredDataPerEachIteration(currentResultDir, iterator, coloredData)
			saveMetaDataPerEachIteration(currentResultDir, iterator, coloredData)
		}()

		break
	}

	wg.Wait()

	log.Println("end collecting coloring data for auc metric")
	return attitudeMetric, labelsFromSentences
}

type ROC struct {
	TPR      float64
	FPR      float64
	Treshold float64
}

func GetAUC(attitudeMetric []float64, labelsFromSentences []string) {
	var roc = make(map[int]ROC, 100)

	tresholds := 0
	for threshold := 0.0; threshold <= 1.0; threshold += 0.01 {
		auc := map[string]int{
			TP: 0,
			FP: 0,
			TN: 0,
			FN: 0,
		}

		for numberLabel, attitude := range attitudeMetric {
			if attitude > threshold && labelsFromSentences[numberLabel] == toloka.True {
				auc[TP]++
			}

			if attitude > threshold && labelsFromSentences[numberLabel] == toloka.False {
				auc[FP]++
			}

			if attitude <= threshold && labelsFromSentences[numberLabel] == toloka.False {
				auc[TN]++
			}

			if attitude <= threshold && labelsFromSentences[numberLabel] == toloka.True {
				auc[FN]++
			}
		}

		currentTPR := calculateTPR(auc)
		currentFPR := calculateFPR(auc)

		currentROC := ROC{
			TPR:      currentTPR,
			FPR:      currentFPR,
			Treshold: threshold,
		}

		roc[tresholds] = currentROC
		tresholds++
	}

	for x := 0; x < 101; x++ {
		fmt.Println(roc[x])
	}

	log.Println("end building data for auc metric")
}

func calculateFPR(auc map[string]int) float64 {
	if auc[FP]+auc[TN] == 0 {
		return 0.0
	}

	return float64(auc[FP]) / float64(auc[FP]+auc[TN])
}

func calculateTPR(auc map[string]int) float64 {
	if auc[TP]+auc[FN] == 0 {
		return 0.0
	}

	return float64(auc[TP]) / float64(auc[TP]+auc[FN])
}

func addAttitudeMetricAndWriteToSnapshotFileSafely(
	mu *sync.Mutex, coloredData ColoredData, fileResultData *os.File,
) {
	mu.Lock()
	defer mu.Unlock()

	for i, coloredCount := range coloredData.Colored {
		attitudeMetric = append(attitudeMetric, float64(coloredCount)/float64(coloredData.Length[i]))
	}

	labelsFromSentences = append(labelsFromSentences, coloredData.Labels...)
	_, err := io.Copy(fileResultData, strings.NewReader(fmt.Sprintf("colored: %v \nlabels: %v \n\n",
		attitudeMetric, labelsFromSentences)))
	if err != nil {
		log.Printf("can not write to result data file %v", err)
	}
}
