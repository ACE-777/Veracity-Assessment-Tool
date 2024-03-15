package wiki

import (
	"log"
	"os/exec"
)

const (
	pythonScript = "test.scrape_wiki_dataset"
	repoDir      = "C:/Users/misha/pythonProject/chatgpt-research"
)

func ScrapeDataFromWikiArticles() {
	log.Printf("Begin scrapping data from Wiki articles")
	cmd := exec.Command(
		"python",
		"-m",
		pythonScript,
	)
	cmd.Dir = repoDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println("Can't execute python script")
		log.Println(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(err)
	}

	log.Printf("End scrapping data from Wiki articles")
	defer stdin.Close()
}
