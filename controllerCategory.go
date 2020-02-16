package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
)

func createCategory(path string) (category Category) {
	var categoryData CategoryData
	pathArr := strings.Split(path, "/")
	name := pathArr[len(pathArr)-1]
	title := "Untitled Category"
	body := []byte("Listing " + name)

	// read metadata form json
	configJSON, err := ioutil.ReadFile(path + "/config.json")
	if err != nil {
		fmt.Print(err)
	} else {
		json.Unmarshal([]byte(configJSON), &categoryData)
		title = categoryData.Title
	}

	return Category{
		Name:     name,
		Path:     path,
		Template: "list.html",
		Title:    title,
		Articles: make(map[string]CategoryArticle),
		Body:     template.HTML(body),
	}
}
