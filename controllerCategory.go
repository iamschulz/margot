package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
)

func createCategory(path string, siteConfig *Config) (category Category) {
	var categoryData CategoryData
	pathArr := strings.Split(path, "/")
	name := pathArr[len(pathArr)-1]
	title := "Untitled Category"
	body := []byte("Listing " + name)
	pagination := 0
	categoryArticles := make(map[string]CategoryArticle)
	/*
		categoryArticles is always 0 at this point
		workaround: build cats as they're needed while looping through arts
	*/

	// read metadata form json
	configJSON, err := ioutil.ReadFile(path + "/config.json")
	if err != nil {
		fmt.Print(err)
	} else {
		json.Unmarshal([]byte(configJSON), &categoryData)
		title = categoryData.Title
	}

	// todo: return map of categories, sorted after pagination
	return Category{
		Name:       name,
		Path:       path,
		Template:   "list",
		Title:      title,
		Pagination: pagination,
		Articles:   categoryArticles,
		Body:       template.HTML(body),
	}
}
