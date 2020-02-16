package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
)

func createPage(path string) (page Page) {
	var pageData PageData
	pathArr := strings.Split(path, "/")
	name := pathArr[len(pathArr)-1]
	title := name

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
		json.Unmarshal([]byte(configJSON), &pageData)
		title = pageData.Title
	}

	return Page{
		Name:     name,
		Path:     path,
		Template: "page.html",
		Title:    title,
		Body:     template.HTML([]byte(htmlString)),
	}
}
