package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func downloadPage(link string) string {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	html := string(content)
	return html
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Must pass series name as an argument")
	}
	var seriesName string
	if len(os.Args) == 2 {
		seriesName = os.Args[1]
	} else {
		seriesName = strings.Join(os.Args[1:], " ")
	}

	seriesName = strings.TrimSpace(seriesName)

	if err := os.MkdirAll(seriesName, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	searchPageUrl := fmt.Sprintf("https://libgen.is/search.php?&req=%s&column=series",
		strings.ReplaceAll(seriesName, " ", "+"))
	searchPage := downloadPage(searchPageUrl)

	downloadLinksRegex := regexp.MustCompile(`http://library.lol/main/[A-Za-z0-9]+`)
	downloadPagesUrls := downloadLinksRegex.FindAllString(searchPage, -1)

	for _, downloadPageUrl := range downloadPagesUrls {
		downloadPage_ := downloadPage(downloadPageUrl)

		downloadLinkRegex := regexp.MustCompile(`<h2><a href="(.*?)">GET</a></h2>`)
		downloadLink := downloadLinkRegex.FindStringSubmatch(downloadPage_)[1]

		getFilenameRegex := regexp.MustCompile(`http://[0-9.]+/main/[0-9]+/[0-9a-z]+/(.*)`)
		filename := getFilenameRegex.FindStringSubmatch(downloadLink)[1]
		filename = strings.ReplaceAll(filename, "%20", " ")
		removePercentEscapedCharactersRegex := regexp.MustCompile(`%2[0-9A-Z]`)
		filename = removePercentEscapedCharactersRegex.ReplaceAllString(filename, "")

		out, err := os.Create(filepath.Join(seriesName, filename))
		defer out.Close()
		if err != nil {
			log.Panic(fmt.Scanf("Failed to download %s", downloadLink))

		}

		resp, err := http.Get(downloadLink)
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Panic(err)
		}
	}
}
