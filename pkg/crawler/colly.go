package crawler

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	"log"
	"pool/pkg/pubsub"
)

var Colly *colly.Collector

func Run(pw Processor, publisher *pubsub.Publisher) {
	initColly(pw, publisher)

	q, _ := queue.New(10, &queue.InMemoryQueueStorage{MaxSize: 10000})
	initUrls(pw, q)

	err := q.Run(Colly)
	if err != nil {
		log.Println(err)
	}

	log.Println("Finished!")
	publisher.Close()
}

func initColly(pw Processor, publisher *pubsub.Publisher) {
	Colly = colly.NewCollector(
		colly.AllowedDomains(pw.getDomain()),
		colly.CacheDir(pw.getCacheDir()),
		// TODO: disable on prod
		colly.Debugger(&debug.LogDebugger{}),
		colly.AllowURLRevisit(),
	)

	// TODO: add proxy for each request

	extensions.RandomUserAgent(Colly)

	//rp, e := proxy.RoundRobinProxySwitcher("socks5://host.docker.internal:4781")
	//if e != nil {
	//	log.Fatal(e)
	//}
	//Colly.SetProxyFunc(rp)

	Colly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	err := Colly.Limit(pw.getLimit())
	if err != nil {
		print("Limit: ", err)
	}

	Colly.OnHTML(pw.IpParser(publisher))
}

func initUrls(cp Processor, q *queue.Queue) {
	urlsCollector := Colly.Clone()
	urlsCollector.OnHTML(cp.UrlParser(q))

	err := urlsCollector.Visit(cp.getBaseURL())
	if err != nil {
		log.Println(err)
	}
}
