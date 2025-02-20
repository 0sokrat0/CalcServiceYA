package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Создаем временную директорию для тестовой конфигурации
	tmpDir, err := ioutil.TempDir("", "configtest")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Формируем тестовое содержимое файла конфигурации в формате YAML
	configContent := `
server:
  host: "localhost"
  port: "8080"
app:
  name: "TestApp"
  computing_power: 4
`
	// Путь к файлу конфигурации
	configFile := filepath.Join(tmpDir, "config.yaml")
	err = ioutil.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Загружаем конфигурацию
	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Проверяем, что значения полей соответствуют ожидаемым
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "TestApp", cfg.App.Name)
	assert.Equal(t, 4, cfg.App.COMPUTING_POWER)
}
