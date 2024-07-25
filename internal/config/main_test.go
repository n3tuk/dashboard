//nolint:paralleltest // these tests cannot operate in parallel
package config_test

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/n3tuk/dashboard/internal/config"
)

const (
	sendConfigName    = "send.yaml"
	serveConfigName   = "serve.yaml"
	validConfigName   = "valid.yaml"
	invalidConfigName = "invalid.yaml"
	missingConfigName = "missing.yaml"
)

var noConfigFile string

// TestDefaultLoadSend tests that the configuration is loaded when searching for
// the send.yaml file in the default configuration folder (which is relative to
// this file).
func TestDefaultLoadSend(t *testing.T) {
	viper.Reset()

	config.Paths = []string{"testdata"}

	err := config.Load(sendConfigName, noConfigFile)
	require.NoError(t, err)

	assert.NotEmpty(t, viper.GetViper().ConfigFileUsed())
	assert.False(t, viper.GetBool("logging.json"))
	assert.Equal(t, "info", viper.GetString("logging.level"))
}

// TestDefaultLoadServe tests that the configuration is loaded when searching
// for the serve.yaml file in the default configuration folder (which is
// relative to this file).
func TestDefaultLoadServe(t *testing.T) {
	viper.Reset()

	config.Paths = []string{"testdata"}

	err := config.Load(serveConfigName, noConfigFile)
	require.NoError(t, err)

	assert.NotEmpty(t, viper.GetViper().ConfigFileUsed())
	assert.True(t, viper.GetBool("logging.json"))
	assert.Equal(t, "debug", viper.GetString("logging.level"))
}

// TestDefaultLoadMissing tests the case when the expected configuration file is
// missing in the environment, and specifically in this case does not error.
func TestDefaultLoadMissing(t *testing.T) {
	viper.Reset()

	config.Paths = []string{"testdata"}

	err := config.Load(missingConfigName, noConfigFile)
	assert.NoError(t, err)
}

// TestExplicitLoadMissing tests the case when the provided configuration file
// is missing, and specifically in this case does error.
func TestExplicitLoadMissing(t *testing.T) {
	viper.Reset()

	config.Paths = []string{"."}
	file := filepath.Join("testdata", missingConfigName)

	err := config.Load(serveConfigName, file)
	assert.Error(t, err)
}

// TestExplicitLoadValid tests that loading a valid configuration file when
// explicitly requested, loads successfully.
func TestExplicitLoadValid(t *testing.T) {
	viper.Reset()

	config.Paths = []string{"."}
	file := filepath.Join("testdata", validConfigName)

	err := config.Load(serveConfigName, file)
	assert.NoError(t, err)
}

// TestExplicitLoadValid tests that loading an invalid configuration file when
// explicitly requested, results in an error.
func TestExplicitLoadInvalid(t *testing.T) {
	viper.Reset()

	config.Paths = []string{"."}
	file := filepath.Join("testdata", invalidConfigName)

	err := config.Load(serveConfigName, file)
	assert.Error(t, err)
}
