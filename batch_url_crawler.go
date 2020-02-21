package go_batch_crawler

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"net"
	"net/http"
	"time"
)

func NewCrawler(urls []string, handler IResponseHandler) *Crawler {
	crawler := &Crawler{}

	crawler.C = colly.NewCollector(
		colly.Async(false),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36"))

	myTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   20 * time.Second,
			KeepAlive: 20 * time.Second,
		}).DialContext,
		MaxIdleConns:          25,
		IdleConnTimeout:       20 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     true,
	}
	crawler.C.WithTransport(myTransport)

	crawler.C.Limit(&colly.LimitRule{
		Delay:       1 * time.Second,
		RandomDelay: 2 * time.Second,
	})

	crawler.Queue, _ = queue.New(
		100, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 20000}, // Use default queue storage
	)

	for _, url := range urls {
		crawler.Queue.AddURL(url)
	}

	crawler.C.OnError(func(r *colly.Response, err error) {
		handler.HandleError(r, err)
	})

	crawler.C.OnResponse(func(r *colly.Response) {
		handler.HandleResponse(r)
	})

	return crawler
}

func (crawler *Crawler) Run() {
	crawler.Queue.Run(crawler.C)
}
