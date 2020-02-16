package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
)

func createArticle(path string) (article Article) {
	var articleData ArticleData
	pathArr := strings.Split(path, "/")
	name := pathArr[len(pathArr)-1]
	title := name
	canonical := ""

	filename := path + "/" + name + ".md"
	htmlString, err := readMarkdown(filename)

	if err != nil {
		fmt.Println(err)
		return
	}

	// read metadata form json
	configJSON, err := ioutil.ReadFile(path + "/config.json")
	if err != nil {
		fmt.Println(err)
	} else {
		json.Unmarshal([]byte(configJSON), &articleData)
		title = articleData.Title
		canonical = articleData.Title
	}

	return Article{
		Name:      name,
		Path:      path,
		Template:  "article.html",
		Canonical: canonical,
		Title:     title,
		Body:      template.HTML([]byte(htmlString)),
	}
}

func registerArticleInCategory(article Article, category Category) {
	articleTitle := article.Title
	articleURL := "/" + article.Name
	category.Articles[article.Name] = CategoryArticle{Title: articleTitle, URL: articleURL}
}
