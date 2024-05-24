package server_logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const (
	repoDir                   = "C:/Users/misha/pythonProject/chatgpt-research"
	buildSearchDatabaseScript = "scripts.search_database"

	useSkips = "True"

	treshold  = 0.8
	uniqColor = 5
)

type coloredDataSecond struct {
	HTML string `json:"html"` //html output for input text
}

type coloredDataFirst struct {
	ResultSources  [][]string  `json:"result_sources"` //sources with variants on each token
	Tokens         []string    `json:"tokens"`         //all tokens from input text
	ResultDistance [][]float64 `json:"result_dists"`   //cos dist result for each token with variants
}

type buildDatabase struct {
	Result string `json:"result"` // send the result of bilding search database
}

func getColoredFirst(userInput string) {
	cmd := exec.Command(
		"python",
		"-m",
		firstAlgorithmPython,
		"--userinput", userInput,
	)

	cmd.Dir = repoDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Can't execute python script, err:%v", err)
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

	var coloredDataRes coloredDataFirst
	if err = json.Unmarshal(output.Bytes(), &coloredDataRes); err != nil {
		log.Printf("Can't convert bytes to json struct coloredData %v", err)
	}

	topLinksPerEachToken := getTopOneSource(coloredDataRes)
	arrayWithTopUniqColors := buildDictForColor(topLinksPerEachToken, uniqColor)

	HTML := buildPageTemplate(coloredDataRes.Tokens, topLinksPerEachToken, arrayWithTopUniqColors, userInput)

	f, err := os.Create("internal/templates/coloring.html")
	if err != nil {
		log.Printf("can not create html page: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(HTML)
	if err != nil {
		log.Printf("can not render html page: %v", err)
	}

	log.Println("end iteration of algorithm")

}

func getTopOneSource(coloredData coloredDataFirst) []string {
	topSources := make([]string, len(coloredData.ResultSources))
	for i, sourcesEachToken := range coloredData.ResultSources {
		topSources[i] = sourcesEachToken[0]

		if coloredData.ResultDistance[i][0] < treshold {
			topSources[i] = ""
		}
	}

	return topSources
}

func buildDictForColor(links []string, uniqColor int) map[string]string {
	filteredLinks := make([]string, 0)

	for _, link := range links {
		if link != "" {
			filteredLinks = append(filteredLinks, link)
		}
	}

	dictionaryOfLinks := make(map[string]int)
	for _, link := range filteredLinks {
		dictionaryOfLinks[link]++
	}

	type kv struct {
		Key   string
		Value int
	}
	var sortedDict []kv
	for k, v := range dictionaryOfLinks {
		sortedDict = append(sortedDict, kv{k, v})
	}

	sort.Slice(sortedDict, func(i, j int) bool {
		return sortedDict[i].Value > sortedDict[j].Value
	})

	linksWithUniqColors := make(map[string]string)
	for i := 0; i < len(sortedDict) && i < uniqColor; i++ {
		linksWithUniqColors[sortedDict[i].Key] = ""
	}

	uniqColorDict := map[string]string{
		"Fuchsia":      "color1",
		"MediumPurple": "color2",
		"DarkViolet":   "color3",
		"DarkMagenta":  "color4",
		"Indigo":       "color5",
	}

	i := 0
	for link := range linksWithUniqColors {
		if i >= len(uniqColorDict) {
			break
		}

		linksWithUniqColors[link] = uniqColorDict[sortedDict[i].Key]
		i++
	}

	return linksWithUniqColors
}

func buildLinkTemplate(tokens, sourceLink []string, linksWithUniqColors map[string]string) (string, int) {
	tokens = removeSpecialSymbolsFromToken(tokens)
	output := ""
	sentenceLength := 0
	countColoredTokenInSentence := 0
	linkTemplate := "<a href='%s' class=\"%s\">%s</a>"
	withoutLinkTemplate := "<a class=\"%s\">%s</a>"

	for i, key := range tokens {
		src := sourceLink[i]
		flag := false

		if strings.HasSuffix(key, ".") || key == ".\"" {
			sentenceLength = 0
		}

		sentenceLength++

		if src != "" {
			for linkColor, color := range linksWithUniqColors {
				if src == linkColor {
					output += fmt.Sprintf(linkTemplate, src, color, key)
					flag = true
					countColoredTokenInSentence++
					break
				}
			}
			if !flag {
				if i%2 != 0 {
					countColoredTokenInSentence++
					output += fmt.Sprintf(linkTemplate, src, "color7", key)
				} else {
					countColoredTokenInSentence++
					output += fmt.Sprintf(linkTemplate, src, "color8", key)
				}
			}
		} else {
			output += fmt.Sprintf(withoutLinkTemplate, "color0", key)
		}
	}

	return output, countColoredTokenInSentence
}

func listOfColors(dictWithUniqColors map[string]string) string {
	output := ""
	listOfArticles := "<div class=\"item_paragraphes\">" +
		"<a href='%s' class=\"%s\">%s</a>" +
		"</div>\n"

	for key, value := range dictWithUniqColors {
		output += fmt.Sprintf(listOfArticles, key, value, key)
	}

	return output
}

var pageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Result</title>
 	<meta name="viewport" content="width=device-width,initial-scale=1">
    <meta http-equiv="x-ua-compatible" content="crhome=1">
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
	<link href="https://fonts.googleapis.com/css2?family=Lato:ital,wght@0,100;0,300;0,400;0,700;0,900;1,100;1,300;1,400;1,700;1,900&family=Platypi:ital,wght@0,300..800;1,300..800&display=swap" rel="stylesheet">
    <link rel="stylesheet" type="text/css" href="../static/style_result_updated.css">
</head>
<body>
<div class="topper">
  <a href="/home/" class="back-button">Back</a>
</div>
<div class="container">
	<div class="item">
		<h3>Colored percentage: {{ coloredCount }} %</h3>
	</div>
	<div class="item">
		<h3>Result</h3>
		{{ result }}
	</div>
	<div class="item">
		<h3>Top paragraphs</h3>
		{{ list_of_colors }}
	</div>
	<div class="item">
		<h3>Input text</h3>
		{{ gpt_response }}
	</div>
</div>
</body>
</html>
`

func buildPageTemplate(
	tokens []string,
	topLinksPerEachToken []string,
	linksWithUniqColors map[string]string,
	userInput string,
) string {
	template := strings.ReplaceAll(pageTemplate, "{{ gpt_response }}", userInput)

	resultOfColor, countColoredTokenInSentenceArray := buildLinkTemplate(
		tokens, topLinksPerEachToken, linksWithUniqColors)
	resultOfListOfColors := listOfColors(linksWithUniqColors)

	template = strings.ReplaceAll(template, "{{ list_of_colors }}", resultOfListOfColors)
	template = strings.ReplaceAll(template, "{{ result }}", resultOfColor)
	template = strings.ReplaceAll(template, "{{ coloredCount }}",
		strconv.Itoa(int(math.Round((float64(countColoredTokenInSentenceArray) / float64(len(tokens)) * 100)))))

	return template
}

func removeSpecialSymbolsFromToken(tokens []string) (modifyTokens []string) {
	for _, token := range tokens {
		token = strings.ReplaceAll(token, "Ġ", " ")
		token = strings.ReplaceAll(token, "Ċ", "</br>")
		modifyTokens = append(modifyTokens, token)
	}

	return
}

func getColoredSecond(userInput string) {
	cmd := exec.Command(
		"python",
		"-m",
		secondAlgorithmPython,
		"--userinput", userInput,
		"--withskip", useSkips,
	)

	cmd.Dir = repoDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Can't execute python script, err:%v", err)
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

	var coloredDataRes coloredDataSecond
	if err = json.Unmarshal(output.Bytes(), &coloredDataRes); err != nil {
		log.Printf("Can't convert bytes to json struct coloredData %v", err)
	}

	f, err := os.Create("internal/templates/coloring.html")
	if err != nil {
		log.Printf("can not create html page: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(coloredDataRes.HTML)
	if err != nil {
		log.Printf("can not render html page: %v", err)
	}

	log.Println("end iteration of algorithm")
}

func buildSearchDatabase(userInput string) {
	articles := strings.Fields(userInput)
	var prepareArticles string
	for _, article := range articles {
		prepareArticles = prepareArticles + strings.Split(article, "https://en.wikipedia.org/wiki/")[1] + " "
	}

	cmd := exec.Command(
		"python",
		"-m",
		buildSearchDatabaseScript,
		"--articles", prepareArticles,
	)
	cmd.Dir = repoDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Can't execute python script, err:%v", err)
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
	fmt.Println("out2::", string(output.Bytes()))
	var buildDatabaseFinal buildDatabase
	if err = json.Unmarshal(output.Bytes(), &buildDatabaseFinal); err != nil {
		log.Printf("Can't convert bytes to json struct buildDatabase %v", err)
	}

	if buildDatabaseFinal.Result == "Success" {
		log.Println("end building search database")
	} else {
		log.Println("incorrect building search database")
	}

}
