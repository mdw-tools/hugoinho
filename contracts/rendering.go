package contracts

import "time"

//////////////////////////////////////////////

type Renderer interface {
	Render(interface{}) (string, error)
}

//////////////////////////////////////////////

type RenderedHomePage struct {
	Title string
	Pages []RenderedHomePageEntry
}

type RenderedHomePageEntry struct {
	Path        string
	Title       string
	Date        time.Time
	Description string
}

//////////////////////////////////////////////

type RenderedArticle struct {
	Title       string
	Description string // TODO: rename to Intro
	Date        time.Time
	Tags        []string
	Content     string
}

//////////////////////////////////////////////

type RenderedTag struct {
	Title string
	Name  string
	Pages []RenderedTagEntry
}

type RenderedTagEntry struct {
	Path  string
	Title string
	Date  string
}

//////////////////////////////////////////////

type RenderedAllTagsListing struct {
	Tags []RenderedAllTagsEntry
}

type RenderedAllTagsEntry struct {
	Name  string
	Path  string
	Count int
}

//////////////////////////////////////////////
