package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
)

func createCustomPage(path string) (customPage CustomPage) {
	var customPageData CustomPageData
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
		json.Unmarshal([]byte(configJSON), &customPageData)
		title = customPageData.Title
	}

	return CustomPage{
		Name:     name,
		Path:     path,
		Template: "customPage",
		Title:    title,
		Body:     template.HTML([]byte(htmlString)),
	}
}
