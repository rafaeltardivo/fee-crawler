package crawler

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// Fee crawler interface
type feeCrawlerInterface interface {
	findPlanIndex(string, *goquery.Selection) (int, error)
	findFee(string, int, *goquery.Selection) (string, error)
}

// SmartMEI implementation of fee crawler interface
type SmartMEIcrawlerStruct struct{}

// Crawls header looking for plan element
// The container is mapped as a table, so the plan index will be used in the fee search
func (s *SmartMEIcrawlerStruct) findPlanIndex(plan string, rows *goquery.Selection) (int, error) {
	planIndex := -1

	header := rows.First()
	header.Children().EachWithBreak(func(i int, child *goquery.Selection) bool {
		if strings.Contains(child.Text(), plan) {
			planIndex = i
			return false
		}
		return true
	})

	if planIndex < 0 {
		return planIndex, crawlError(fmt.Sprintf("plan not found: %s", plan))
	}
	return planIndex, nil
}

// Crawls rows looking for fee element
// The container is mapped as a table, so the plan index is the search pivot
func (s *SmartMEIcrawlerStruct) findFee(plan string, index int, rows *goquery.Selection) (string, error) {
	var fee = ""

	rows.Children().EachWithBreak(func(i int, row *goquery.Selection) bool {
		if strings.Contains(row.Text(), "Transferência") {
			// Since the plan index is child based, subtracts one
			// to correspond to sibling index
			fee = row.Siblings().Eq(index - 1).Text()
			return false
		}
		return true
	})

	if fee == "" {
		return fee, crawlError(fmt.Sprintf("fee not found"))
	}
	return fee, nil
}

// Crawling proccess enntrypoint
// Returns raw fee
func CrawlFee(plan string) (string, error) {
	var container *colly.HTMLElement
	crawler := NewSmartMEICrawler()
	collector := NewSmartMEICollector()

	collector.OnHTML(`div[id="tarifas-2"]`, func(e *colly.HTMLElement) {
		container = e
	})
	collector.Visit("http://smartmei.com.br/")

	if container == nil {
		return "", crawlError(fmt.Sprintf("container not found"))
	}

	return crawl(plan, container.DOM.Children(), crawler)
}

// Manages crawling proccess and returns fee
func crawl(plan string, rows *goquery.Selection, crawler feeCrawlerInterface) (string, error) {

	index, err := crawler.findPlanIndex(plan, rows)
	if err != nil {
		return "", err
	}

	fee, err := crawler.findFee(plan, index, rows)
	if err != nil {
		return "", err
	}

	return fee, nil
}

// Returns a new SmartMEICollector
func NewSmartMEICollector() *colly.Collector {
	agents := [4]string{
		"Chrome 70.0.3538.77",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:77.0) Gecko/20190101 Firefox/77.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.75.14 (KHTML, like Gecko) Version/7.0.3 Safari/7046A194A",
		"Opera/9.80 (X11; Linux i686; Ubuntu/14.10) Presto/2.12.388 Version/12.16.2",
	}
	// Returns a random user agent
	// The goal is to reduce action patterns
	userAgent := agents[rand.Intn(len(agents)-1)+1]

	return colly.NewCollector(
		colly.UserAgent(userAgent),
		colly.AllowedDomains("smartmei.com.br", "www.smartmei.com.br"),
		colly.MaxDepth(2),
	)
}

// Returns a new SmartMEI Crawler
func NewSmartMEICrawler() feeCrawlerInterface {
	return &SmartMEIcrawlerStruct{}
}
