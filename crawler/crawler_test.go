package crawler

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/onsi/gomega"
)

func TestNewSmartMEICollector(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	expectedDomains := []string{"smartmei.com.br", "www.smartmei.com.br"}

	collector := NewSmartMEICollector()

	g.Expect(collector.AllowedDomains).Should(gomega.Equal(expectedDomains), "Return should be equal SmartMEI allowed domains")
}

func TestNewSmartMEICrawler(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	crawler := NewSmartMEICrawler()
	_, sameType := crawler.(feeCrawlerInterface)

	g.Expect(sameType).To(gomega.BeTrue(), "crawler is-one feeCrawlerInterface")
}

type mockedFindPlanIndexCrawler struct{}

func (c *mockedFindPlanIndexCrawler) findPlanIndex(plan string, header *goquery.Selection) (int, error) {
	return -1, crawlError("Mocked error")
}

func (c *mockedFindPlanIndexCrawler) findFee(plan string, index int, rows *goquery.Selection) (string, error) {
	return "7,00", nil
}

func TestCrawlIndexError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	var container *goquery.Selection
	crawler := &mockedFindPlanIndexCrawler{}

	ret, err := crawl("Plan", container, crawler)

	g.Expect(ret).Should(gomega.Equal(""), "Return should be empty")
	g.Expect(err).Should(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).Should(gomega.MatchError(crawlError("Mocked error")), "Error should be a crawlError")
}

type mockedFindFeeErrorCrawler struct{}

func (c *mockedFindFeeErrorCrawler) findPlanIndex(plan string, header *goquery.Selection) (int, error) {
	return 1, nil
}

func (c *mockedFindFeeErrorCrawler) findFee(plan string, index int, rows *goquery.Selection) (string, error) {
	return "", crawlError("Mocked error")
}

func TestCrawlFeeError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	var container *goquery.Selection
	crawler := &mockedFindFeeErrorCrawler{}

	ret, err := crawl("Plan", container, crawler)

	g.Expect(ret).Should(gomega.Equal(""), "Return should be empty")
	g.Expect(err).Should(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).Should(gomega.MatchError(crawlError("Mocked error")), "Error should be a crawlError")
}
