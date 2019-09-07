package pulli

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/x-cray/logrus-prefixed-formatter"
	"testing"
)

func TestLogging_InitLogging_setsFormatter(t *testing.T) {
	InitLogging()
	assert.Exactly(t, &prefixed.TextFormatter{
		ForceFormatting: true,
		ForceColors:     true,
		SpacePadding:    12,
	}, logrus.StandardLogger().Formatter)
}

func TestLogging_InitLogging_setsDefaultLevel(t *testing.T) {
	InitLogging()

	assert.Exactly(t, logrus.DebugLevel, logrus.GetLevel())
}

func TestLogging_newEntry_returnsEntryWithFieldRepositoryName(t *testing.T) {
	repoName := "testRepoName"
	log := &logging{currentRepoDisplayName: repoName}
	entry := log.newEntry()

	assert.Exactly(t, repoName, entry.Data["repository"].(string))
}
