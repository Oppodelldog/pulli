package pulli

import (
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var log = &logging{}

type logging struct {
	currentRepoDisplayName string
}

func (l *logging) newEntry() *logrus.Entry {
	return logrus.WithField("repository", l.currentRepoDisplayName)
}

func InitLogging() {
	logrus.SetFormatter(&prefixed.TextFormatter{
		ForceFormatting: true,
		ForceColors:     true,
		SpacePadding:    12,
	})
	logrus.SetLevel(logrus.DebugLevel)
}
