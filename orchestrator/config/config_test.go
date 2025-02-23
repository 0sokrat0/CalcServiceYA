package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	
	tmpDir, err := ioutil.TempDir("", "configtest")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)


	configContent := `
server:
  host: localhost
  port: 8080
  
app:
  name: "Orchestrator"
  time_addition_ms: 1000
  time_subtraction_ms: 1200
  time_multiplication_ms: 2000
  time_division_ms: 2500
`
	
	configFile := filepath.Join(tmpDir, "config.yaml")
	err = ioutil.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	
	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "Orchestrator", cfg.App.Name)
	assert.Equal(t, int64(1000), cfg.App.TIME_ADDITION_MS)
	assert.Equal(t, int64(1200), cfg.App.TIME_SUBTRACTION_MS)
	assert.Equal(t, int64(2000), cfg.App.TIME_MULTIPLICATION_MS)
	assert.Equal(t, int64(2500), cfg.App.TIME_DIVISION_MS)
}
