package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/golang-commonmark/markdown"
)

// Create a site config from the corresponding json file.
// Only includes page title for now
func getSiteConfig() (config *Config) {
	siteConfigJSON, err := ioutil.ReadFile("./src/config/siteConfig.json")
	if err != nil {
		fmt.Print(err)
	}

	json.Unmarshal([]byte(siteConfigJSON), &config)
	return
}

// Create an Index of all contents
// Run this at startup and refer to it on requests
func createIndex() (index *Index) {
	searchDir := "src/markdown"
	categoryList := make(map[string]Category)
	articleList := make(map[string]Article)
	pageList := make(map[string]Page)

	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if path == searchDir {
			return nil
		}

		fileInfo, err := os.Stat(path)
		pathArr := strings.Split(path[len(searchDir+"/"):], "/") // e.g. [src markdown articles cat1 art1 art1.md]

		if fileInfo.IsDir() == true {
			switch {
			// if category
			case pathArr[0] == "articles" && len(pathArr) == 2:
				categoryList[fileInfo.Name()] = createCategory(path)

			// if article
			case pathArr[0] == "articles" && len(pathArr) == 3:
				article := createArticle(path)
				catName := pathArr[1]
				articleList[fileInfo.Name()] = article
				registerArticleInCategory(article, categoryList[catName])

			// if page
			case pathArr[0] == "pages" && len(pathArr) == 2:
				pageList[fileInfo.Name()] = createPage(path)
			}
		}

		return nil
	})

	return &Index{
		Categories: categoryList,
		Pages:      pageList,
		Articles:   articleList,
	}
}

func getPartialList() (allPartials []string) {
	partialDir := "./src/templates/partials"
	partials, _ := ioutil.ReadDir(partialDir)

	for _, partial := range partials {
		fileName := partial.Name()
		if strings.HasSuffix(fileName, ".html") {
			partialPath := filepath.Join(partialDir, fileName)
			allPartials = append(allPartials, partialPath)
		}
	}

	return
}

func requestHandler(index *Index, config *Config, templateList []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var pageType string
		var viewData interface{}
		var templateFile string
		route := r.URL.Path[len("/"):]

		if r.URL.Path == "/" {
			pageType = config.FrontPageType
			route = config.FrontPageName
			// needs error handling
		}
		if _, ok := index.Categories[route]; ok {
			pageType = "category"
		} else if _, ok := index.Pages[route]; ok {
			pageType = "page"
		} else if _, ok := index.Articles[route]; ok {
			pageType = "article"
		}

		switch pageType {
		case "category":
			templateFile = index.Categories[route].Template
			viewData = index.Categories[route]
		case "page":
			templateFile = index.Pages[route].Template
			viewData = index.Pages[route]
		case "article":
			templateFile = index.Articles[route].Template
			viewData = index.Articles[route]
		default:
			templateFile = "error.html"
			viewData = Page{
				Name:     "Error",
				Path:     "",
				Template: "error.html",
				Title:    "Error",
				Body:     template.HTML([]byte("Oh no :(")),
			}
		}

		if err == nil {
			templateList = append(getPartialList(), "./src/templates/"+templateFile)
			templateName := strings.Split(templateFile, ".html")[0] // strip ".html"
			renderTemplate(w, templateList, templateName, viewData)
		} else {
			fmt.Println("an error occured")
		}
	})
}

func renderTemplate(w http.ResponseWriter, templateList []string, templateName string, p interface{}) {
	t, _ := template.New("").ParseFiles(templateList...)
	t.ExecuteTemplate(w, templateName, p)
}

func readMarkdown(mdPath string) (mdString string, err error) {
	artMD, err := ioutil.ReadFile(mdPath)
	if err != nil {
		fmt.Print(err)
		return "", err
	}

	md := markdown.New(markdown.Linkify(false))
	mdString = md.RenderToString([]byte(artMD))
	return mdString, nil
}

func main() {
	siteConfig := getSiteConfig()
	siteIndex := createIndex()
	partialList := getPartialList()

	http.Handle("/", requestHandler(siteIndex, siteConfig, partialList))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
todo:
- create pagination for cats
- add template type config for pages and articles
- add error handling
- create navigation from cats and pages
- use siteConfig
- add rss feed
- add xml sitemap
- add updateIndex method
*/
