package crawler

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	logger.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestNewSmartMEICollector(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	collector := NewSmartMEICollector()

	g.Expect(collector.MaxDepth).To(gomega.Equal(2), "Max depth should be 2")
}

func TestNewSmartMEICrawler(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	crawler := NewSmartMEICrawler()
	_, isOne := crawler.(feeCrawlerInterface)

	g.Expect(isOne).To(gomega.BeTrue(), "Crawler is-one feeCrawlerInterface")
}

type mockedFindPlanIndexCrawler struct{}

func (c *mockedFindPlanIndexCrawler) findPlanIndex(plan string, header *goquery.Selection) (int, error) {
	return -1, crawlError("index not found mocked error")
}

func (c *mockedFindPlanIndexCrawler) findFee(plan string, index int, rows *goquery.Selection) (string, error) {
	return "7,00", nil
}

func TestCrawlIndexNotFoundError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	var container *goquery.Selection
	crawler := &mockedFindPlanIndexCrawler{}

	ret, err := crawl("Plan", container, crawler)

	g.Expect(ret).To(gomega.Equal(""), "Return should be empty")
	g.Expect(err).To(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).To(gomega.MatchError(crawlError("index not found mocked error")), "Error should be a crawlError")
}

type mockedFindFeeErrorCrawler struct{}

func (c *mockedFindFeeErrorCrawler) findPlanIndex(plan string, header *goquery.Selection) (int, error) {
	return 1, nil
}

func (c *mockedFindFeeErrorCrawler) findFee(plan string, index int, rows *goquery.Selection) (string, error) {
	return "", crawlError("fee not found mocked error")
}

func TestCrawlFeeNotFoundError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	var container *goquery.Selection
	crawler := &mockedFindFeeErrorCrawler{}

	ret, err := crawl("Plan", container, crawler)

	g.Expect(ret).To(gomega.Equal(""), "Return should be empty")
	g.Expect(err).To(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).To(gomega.MatchError(crawlError("fee not found mocked error")), "Error should be a crawlError")
}
