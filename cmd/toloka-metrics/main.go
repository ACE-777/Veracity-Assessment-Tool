package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"toloka-metrics/internal/toloka"
)

func main() {
	res := toloka.NewResponseData()

	for k, v := range res {
		fmt.Println("key:", k, "value", v, "\n")
		userinput := ""
		for _, sentence := range v {
			userinput += sentence.Text
			userinput += ". "
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

		fmt.Println(" Result: ", output.String())
		break
	}

	//cmd := exec.Command("python", "-m", "test.color_build_data", "--userinput", "Elvis")
	//cmd.Dir = "C:/Users/misha/chatgpt-research"
	//
	//stdin, err := cmd.StdinPipe()
	//if err != nil {
	//	log.Println("Can't execute python script")
	//	log.Println(err)
	//}
	//defer stdin.Close()
	//
	//var output bytes.Buffer
	//
	//cmd.Stdout = &output
	//cmd.Stderr = os.Stderr
	//if err = cmd.Start(); err != nil { //Use start, not run
	//	log.Printf("error in starting python commnad: %v", err)
	//}
	//
	////_, err = io.WriteString(stdin, addNewlineIfMissing(faceVectorStr))
	////if err != nil {
	////	log.Println(err)
	////}
	//
	////log.Println("Vector was given:")
	////log.Println(addNewlineIfMissing(faceVectorStr))
	//
	//err = cmd.Wait()
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//fmt.Println(" Result: ", output.String())

}
