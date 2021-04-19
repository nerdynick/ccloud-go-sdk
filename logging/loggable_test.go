package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevelChange(t *testing.T) {
	assert := assert.New(t)
	logger := New("TestLogger")

	assert.True(logger.Log.Core().Enabled(InfoLevel))
	assert.False(logger.Log.Core().Enabled(DebugLevel))
	logger.Log.Info("Info1")
	logger.Log.Debug("Debug1")
	logger.Log.Warn("Warn1")
	logger.Log.Error("Error1")

	logger.SetLogLevel(DebugLevel)
	assert.True(logger.Log.Core().Enabled(InfoLevel))
	assert.True(logger.Log.Core().Enabled(DebugLevel))
	logger.Log.Info("Info2")
	logger.Log.Debug("Debug2")
	logger.Log.Sync()
}
