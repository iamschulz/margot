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
func createIndex(siteConfig *Config) (index *Index) {
	searchDir := "src/markdown"
	categoryList := make(map[string]Category)
	articleList := make(map[string]Article)
	customPageList := make(map[string]CustomPage)

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
				categoryList[fileInfo.Name()] = createCategory(path, siteConfig)

			// if article
			case pathArr[0] == "articles" && len(pathArr) == 3:
				article := createArticle(path)
				catName := pathArr[1]
				articleList[fileInfo.Name()] = article
				registerArticleInCategory(article, categoryList[catName])

			// if customPage
			case pathArr[0] == "customPages" && len(pathArr) == 2:
				customPageList[fileInfo.Name()] = createCustomPage(path)
			}
		}

		return nil
	})

	return &Index{
		Categories:  categoryList,
		CustomPages: customPageList,
		Articles:    articleList,
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

func requestHandler(index *Index, siteconfig *Config, templateList []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var pageType string
		var viewData interface{}
		var templateName string
		route := r.URL.Path[len("/"):]
		path := strings.Split(route, "/")

		if r.URL.Path == "/" {
			pageType = siteconfig.FrontPageType
			route = siteconfig.FrontPageName
			// needs error handling
		}
		if _, ok := index.Categories[path[0]]; ok {
			pageType = "category"
		} else if _, ok := index.CustomPages[route]; ok {
			pageType = "customPage"
		} else if _, ok := index.Articles[route]; ok {
			pageType = "article"
		}

		switch pageType {
		case "category":
			templateName = index.Categories[path[0]].Template
			viewData = index.Categories[path[0]]
			fmt.Println(index.Categories[path[0]]) // todo: get paged category
		case "customPage":
			templateName = index.CustomPages[route].Template
			viewData = index.CustomPages[route]
		case "article":
			templateName = index.Articles[route].Template
			viewData = index.Articles[route]
		default:
			templateName = "error"
			viewData = CustomPage{
				Name:     "Error",
				Path:     "",
				Template: "error",
				Title:    "Error",
				Body:     template.HTML([]byte("")),
			}
		}

		if err == nil {
			templateList = append(getPartialList(), "./src/templates/pageTypes/"+templateName+".html")
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
	siteIndex := createIndex(siteConfig)
	partialList := getPartialList()

	http.Handle("/", requestHandler(siteIndex, siteConfig, partialList))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
todo:
- create pagination for cats
- add template type config for customPages and articles
- add error handling
- create navigation from cats and customPages
- use siteConfig
- add rss feed
- add xml sitemap
- add updateIndex method
- add tests
*/
