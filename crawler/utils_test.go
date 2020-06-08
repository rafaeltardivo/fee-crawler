package crawler

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestSanitizeFeeStringNoValuesError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	amount, description, err := sanitizeFeeString("some string")

	g.Expect(amount).To(gomega.Equal(""), "Amount should be empty")
	g.Expect(description).To(gomega.Equal(""), "Description should be empty")
	g.Expect(err).To(gomega.MatchError(crawlError("could not split amount and description")), "Error should be a crawlError")
}

func TestSanitizeFeeStringNoFeeAmountError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	amount, description, err := sanitizeFeeString("some \n string")

	g.Expect(amount).To(gomega.Equal(""), "Amount should be empty")
	g.Expect(description).To(gomega.Equal(""), "Description should be empty")
	g.Expect(err).To(gomega.MatchError(crawlError("fee amount not found")), "Error should be a crawlError")
}

func TestSanitizeFeeStringNoDescription(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	amount, description, err := sanitizeFeeString("\n                    R$ 7,00                ")

	g.Expect(amount).To(gomega.Equal("7.00"), "Amount should be equal to 7.00")
	g.Expect(description).To(gomega.Equal(""), "Description should be empty")
	g.Expect(err).ToNot(gomega.HaveOccurred(), "An error should not have occurred")
}

func TestSanitizeFeeString(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	amount, description, err := sanitizeFeeString("\n                    R$ 7,00               \n Some description ")

	g.Expect(amount).To(gomega.Equal("7.00"), "Amount should be equal to 7.00")
	g.Expect(description).To(gomega.Equal("Some description"), "Description should be equal to Some description")
	g.Expect(err).ToNot(gomega.HaveOccurred(), "An error should not have occurred")
}

func TestSanitizeAmountFeeAmountNotFound(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	amount, err := sanitizeAmount("no value")

	g.Expect(amount).To(gomega.Equal(""), "Amount should be empty")
	g.Expect(err).To(gomega.HaveOccurred(), "An error should not have occurred")
	g.Expect(err).To(gomega.MatchError(crawlError("fee amount not found")), "Error should be a crawlError")
}

func TestSanitizeAmountFeeAmount(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	amount, err := sanitizeAmount("R$ 7,00")

	g.Expect(amount).To(gomega.Equal("7,00"), "Amount should be equal to 7,00")
	g.Expect(err).ToNot(gomega.HaveOccurred(), "An error should not have occurred")
}

func TestSanitizeDescription(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	description := sanitizeDescription("    description    ")

	g.Expect(description).To(gomega.Equal("description"), "Description should not have trailing spaces")
}

func TestNormalizeAmountToBRL(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	description := normalizeAmountToBRL("7,00")

	g.Expect(description).To(gomega.Equal("7.00"), "Amount should be equal to 7.00")
}
