package metrics

import (
	"fmt"
	"sort"
	"strings"
)

const (
	treshold = 0.8
)

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

func getTopOneSource(coloredData ColoredData) []string {
	topSources := make([]string, len(coloredData.ResultSources))
	for i, sourcesEachToken := range coloredData.ResultSources {
		topSources[i] = sourcesEachToken[0]

		if coloredData.ResultDistance[i][0] < treshold {
			topSources[i] = ""
		}
	}

	return topSources
}

func buildLinkTemplate(tokens []string, sourceLink []string, linksWithUniqColors map[string]string) (string, []int, []int) {
	tokens = removeSpecialSymbolsFromToken(tokens)
	output := ""
	sentenceLength := 0
	countColoredTokenInSentence := 0
	sentenceLengthArray := make([]int, 0)
	countColoredTokenInSentenceArray := make([]int, 0)
	linkTemplate := "<a href=\"%s\" class=\"%s\">%s</a>"

	for i, key := range tokens {
		src := sourceLink[i]
		flag := false

		if key == "." {
			sentenceLengthArray = append(sentenceLengthArray, sentenceLength)
			countColoredTokenInSentenceArray = append(countColoredTokenInSentenceArray, countColoredTokenInSentence)
			sentenceLength = 0
			countColoredTokenInSentence = 0
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
			output += fmt.Sprintf(linkTemplate, "", "color0", key)
		}
	}

	return output, sentenceLengthArray, countColoredTokenInSentenceArray
}

func listOfColors(dictWithUniqColors map[string]string) string {
	output := ""
	listOfArticles := "<a href=\"%s\" class=\"%s\">%s</a></br>"

	for key, value := range dictWithUniqColors {
		output += fmt.Sprintf(listOfArticles, key, value, key)
	}

	output += "</br>"
	return output
}

var pageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Result</title>
    <link rel="stylesheet" type="text/css" href="../cmd/static/style_result.css">
</head>
<body>
<h1>Result of research</h1>
<pre><b>Input text:</b></pre>
{{ gpt_response }}
<pre><b>Top paragraphs:</b></pre>
{{ list_of_colors }}
<pre><b>Result:</b></pre>
{{ result }}
</body>
</html>
`

func buildPageTemplate(tokens []string, topLinksPerEachToken []string, linksWithUniqColors map[string]string) ([]int, []int, string) {
	template := strings.ReplaceAll(pageTemplate, "{{ gpt_response }}", " ")
	resultOfColor, sentenceLengthArray, countColoredTokenInSentenceArray := buildLinkTemplate(tokens, topLinksPerEachToken, linksWithUniqColors)
	resultOfListOfColors := listOfColors(linksWithUniqColors)
	template = strings.ReplaceAll(template, "{{ list_of_colors }}", resultOfListOfColors)
	template = strings.ReplaceAll(template, "{{ result }}", resultOfColor)

	return sentenceLengthArray, countColoredTokenInSentenceArray, template
}

func removeSpecialSymbolsFromToken(tokens []string) (modifyTokens []string) {
	for _, token := range tokens {
		token = strings.ReplaceAll(token, "Ġ", " ")
		token = strings.ReplaceAll(token, "Ċ", "</br>")
		modifyTokens = append(modifyTokens, token)
	}

	return
}
