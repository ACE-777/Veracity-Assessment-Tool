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

	"toloka-metrics/internal/toloka"
)

type ColoredData struct {
	File     string `json:"file"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Length   []int  `json:"length"`
	Colored  []int  `json:"colored"`
	Labels   []string
}

func GetColored(res map[toloka.ResponseData][]toloka.Sentence) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 4)

	var coloredDataResult []ColoredData

	for k, v := range res {
		k := k
		v := v

		wg.Add(1)
		go func() {

			sem <- struct{}{}
			defer func() {
				wg.Done()
				<-sem
			}()
			//fmt.Println(k, v)
			labels := make([]string, len(v))
			userinput := ""
			for sentenceNumber, sentence := range v {
				userinput += sentence.Text
				userinput += ". "
				labels[sentenceNumber] = sentence.Label
			}

			cmd := exec.Command(
				"python",
				"-m",
				"test.color_build_data",
				"--userinput", userinput,
				"--file", k.File,
				"--question", strconv.Itoa(int(k.Question)),
				"--answer", strconv.Itoa(int(k.Answer)),
			)

			cmd.Dir = "C:/Users/misha/chatgpt-research"

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
			coloredDataResult = append(coloredDataResult, coloredData)

		}()
		break
	}

	wg.Wait()
	fmt.Println(coloredDataResult)

}
