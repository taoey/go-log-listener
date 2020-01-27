package listener

import (
	"github.com/sirupsen/logrus"
	"os"
)

// 初始化日志
var LOG = logrus.New()

func init() {
	LOG.Out = os.Stdout
	LOG.SetLevel(logrus.DebugLevel)
}
