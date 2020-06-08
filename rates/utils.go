package rates

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
}
