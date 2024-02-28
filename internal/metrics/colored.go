package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"

	toloka "toloka-metrics/internal/toloka"
)

const (
	pythonScript = "test.color_build_data"
	repoDir      = "C:/Users/misha/pythonProject/chatgpt-research"

	TP = "TP"
	FP = "FP"
	TN = "TN"
	FN = "FN"
)

type ColoredData struct {
	File     string `json:"file"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Length   []int  `json:"length"`
	Colored  []int  `json:"colored"`
	Labels   []string
}

func GetColored(res map[toloka.ResponseData][]toloka.Sentence) ([]float64, []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var mu2 sync.Mutex
	sem := make(chan struct{}, 5)

	var attitudeMetric []float64
	var labelsFromSentences []string
	var iterator = 0

	for k, v := range res {
		k := k
		v := v

		wg.Add(1)
		go func() {

			sem <- struct{}{}
			defer func() {
				wg.Done()
				<-sem
				fmt.Println("iterator::", iterator)
				iterator++
			}()

			labels := make([]string, len(v))
			userInput := ""
			mu.Lock()
			for sentenceNumber, sentence := range v {
				userInput += sentence.Text
				userInput += ". "
				labelsFromSentences = append(labelsFromSentences, sentence.Label)
				labels[sentenceNumber] = sentence.Label
			}

			defer mu.Unlock()

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
			mu2.Lock()
			for i, coloredCount := range coloredData.Colored {
				attitudeMetric = append(attitudeMetric, float64(coloredCount)/float64(coloredData.Length[i]))
			}

			defer mu2.Unlock()
		}()
		//break
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
