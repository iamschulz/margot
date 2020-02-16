# Margot is a **Mar**kdown **Go** **T**hing
It produces a website from some Markdown files.

## How to use
### Structure and adding content
A Margot item is usually one of the following three types:
- Page

    A Page is a standalone site without any categorization. Good examples would be an "About Me"-Page, an imprint or some featured content.
    Pages will not show up in the RSS-feed.
    
    To create a page, add your content to `.markdown/pages/YourPage/YourPage.md`

- Category

    A Category is a set of articles. Categories are automatically created from the folder structure in `./markdown/articles`

- Article

    Articles are categorized content. They will show up in their respective categories and in the RSS-feed.

    To create an article, add your content to `.markdown/articley/YourArticle/YourArticle.md`

### Configuring content

It's recommended to place a `config.json` next to your `.md`-files and in category directories. Within those config files, you can specify a title some additional information depending on the item type.

- Page

- Category

- Article

- site config