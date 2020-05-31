package crawler

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type crawlerStruct struct {
	collector *colly.Collector
}

type crawlerInterface interface {
	GetPlanFee(string) (string, error)
}

// GetPlanFee returns the fee for the given plan or error. The DOM element is mapped
// as a matrix in order to increase page layout change resilience
// If the plan index could not be found, errPlanIndexNotFound will be returned
// If the fee could not be found, errPlanIndexNotFound will be returned
func (c *crawlerStruct) GetPlanFee(plan string) (string, error) {
	fee := ""
	planIndex := -1

	c.collector.OnHTML(`div[id="tarifas-2"]`, func(e *colly.HTMLElement) {
		rows := e.DOM.Children()
		title := rows.First()

		title.Children().EachWithBreak(func(i int, child *goquery.Selection) bool {
			if strings.Contains(child.Text(), plan) {
				planIndex = i
				return false
			}
			return true
		})

		if planIndex >= 0 {
			rows.Children().EachWithBreak(func(i int, child *goquery.Selection) bool {
				if strings.Contains(child.Text(), "Transferência") {

					// Since the plan index is child based, subtracts one
					// to correspond to sibling index
					fee = child.Siblings().Eq(planIndex - 1).Text()
					return false
				}
				return true
			})
		}
	})

	c.collector.Visit("http://smartmei.com.br/")

	if planIndex == -1 {
		return "", errPlanIndexNotFound
	}

	if fee == "" {
		return "", errPlanFeeNotFound
	}

	return fee, nil
}

// Returns a new crawler according to crawlerInterface
func NewCrawler() crawlerInterface {
	return &crawlerStruct{
		collector: colly.NewCollector(
			colly.AllowedDomains("smartmei.com.br", "www.smartmei.com.br"),
			colly.MaxDepth(1),
		),
	}
}
