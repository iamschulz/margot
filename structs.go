package main

import "html/template"

// Config stuct contains page meta data
type Config struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Language       string `json:"language"`
	FrontPageType  string `json:"frontPageType"`
	FrontPageName  string `json:"frontPageName"`
	PaginationSize int    `json:"paginationSize"`
}

// Index struct contains all category, custompage and article names
type Index struct {
	Categories  map[string]Category
	CustomPages map[string]CustomPage
	Articles    map[string]Article
}

// Category struct contains category information
type Category struct {
	Name       string
	Path       string
	Template   string
	Title      string
	Pagination int
	Articles   map[string]CategoryArticle
	Body       template.HTML
}

type PagedCategory struct {
	Name     string
	Path     string
	Template string
	Title    string
	Articles map[string]CategoryArticle
	Body     template.HTML
}

// CategoryArticle struct contains information about an Article in a Category
type CategoryArticle struct {
	Title string
	URL   string
}

// CustomPage struct contains page information
type CustomPage struct {
	Name     string
	Path     string
	Template string
	Title    string
	Body     template.HTML
}

// Article struct contains article information
type Article struct {
	Name      string
	Path      string
	Template  string
	Canonical string
	Title     string
	Body      template.HTML
}

// CategoryData contains meta information about a category and matches the categories config.json files
type CategoryData struct {
	Title string `json:"title"`
}

// CustomPageData contains meta information about a customPage and matches the customPage config.json file
type CustomPageData struct {
	Title string `json:"title"`
}

// ArticleData contains meta information about an article and matches the articles config.json file
type ArticleData struct {
	Title     string `json:"title"`
	Canonical string `json:"canonical"`
}
