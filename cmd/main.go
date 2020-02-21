package main

import (
	"encoding/json"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	crawler "github.com/pytorchtw/go-batch-crawler"
	"io/ioutil"
	"log"
	"strings"
)

type Job struct {
	Title  string
	URL    string
	Status int
}

type Result struct {
	Jobs   []*Job
	Errors []*error
}

type ResponseHandler struct {
	Result Result
}

func ParseTitle(body string) (string, error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	title := htmlquery.FindOne(doc, `//title`)
	if title != nil {
		return htmlquery.InnerText(title), nil
	}
	return "", nil
}

func (h *ResponseHandler) HandleResponse(r *colly.Response) {
	log.Println("OK:", r.Request.URL.String())

	title, err := ParseTitle(string(r.Body))
	if err != nil {
		panic("error parsing title from response")
	}

	job := &Job{
		Title:  title,
		URL:    r.Request.URL.String(),
		Status: r.StatusCode,
	}
	h.Result.Jobs = append(h.Result.Jobs, job)
}

func (h *ResponseHandler) HandleError(r *colly.Response, err error) {
	log.Println("ERROR:", r.Request.URL.String(), r.StatusCode)
	h.Result.Errors = append(h.Result.Errors, &err)
}

func main() {
	newUrls := []string{
		"http://www.bbc.com/",
		"http://edition.cnn.com/",
		"https://www.theguardian.com/international",
		"http://www.breitbart.com/",
		"https://www.infowars.com/",
		"http://www.foxnews.com/",
		"http://www.nbcnews.com/",
		"http://www.theonion.com/",
	}

	handler := &ResponseHandler{}
	c := crawler.NewCrawler(newUrls, handler)
	c.Run()

	file, _ := json.MarshalIndent(handler.Result, "", " ")
	_ = ioutil.WriteFile("result.json", file, 0644)
}
