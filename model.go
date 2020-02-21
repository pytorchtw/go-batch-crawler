package go_batch_crawler

import (
	"github.com/gocolly/colly"
)

type IQueue interface {
	AddURL(string) error
	Run(*colly.Collector) error
}

type IResponseHandler interface {
	HandleResponse(r *colly.Response)
	HandleError(r *colly.Response, err error)
}

type Crawler struct {
	C     *colly.Collector
	Queue IQueue
}
