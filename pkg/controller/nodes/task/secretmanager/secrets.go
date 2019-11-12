package secretmanager

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flytestdlib/logger"
)

type FileEnvSecretManager struct {
	secretPath string
	envPrefix  string
}

func (f FileEnvSecretManager) Get(ctx context.Context, key string) (string, error) {
	envVar := fmt.Sprintf("%s%s", f.envPrefix, key)
	v, ok := os.LookupEnv(envVar)
	if ok {
		return v, nil
	}
	secretFile := filepath.Join(f.secretPath, key)
	if _, err := os.Stat(secretFile); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("secrets not found")
		}
		return "", err
	}
	logger.Debugf(ctx, "reading secrets from filePath [%s]", secretFile)
	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func NewFileEnvSecretManager(cfg *Config) core.SecretManager {
	return FileEnvSecretManager{
		secretPath: cfg.SecretFilePrefix,
		envPrefix:  cfg.EnvironmentPrefix,
	}
}
