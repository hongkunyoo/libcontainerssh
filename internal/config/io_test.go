package config_test

import (
	"bytes"
	"context"
	"testing"

	configuration "github.com/containerssh/libcontainerssh/config"
	"github.com/containerssh/libcontainerssh/internal/config"
	"github.com/containerssh/libcontainerssh/internal/structutils"
	"github.com/containerssh/libcontainerssh/log"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestSaveLoadYAML(t *testing.T) {
	testSaveLoad(t, config.FormatYAML)
}

func TestSaveLoadJSON(t *testing.T) {
	testSaveLoad(t, config.FormatJSON)
}

func testSaveLoad(t *testing.T, format config.Format) {
	// region Setup
	logger := log.NewTestLogger(t)

	cfg := &configuration.AppConfig{}
	newCfg := &configuration.AppConfig{}
	structutils.Defaults(cfg)

	cfg.Auth.Webhook.URL = "http://localhost:8080"

	buf := &bytes.Buffer{}
	// endregion

	// region Save
	saver, err := config.NewWriterSaver(
		buf,
		logger,
		format,
	)
	assert.NoError(t, err)
	err = saver.Save(cfg)
	assert.Nil(t, err, "failed to load config (%v)", err)
	// endregion

	// region Load
	loader, err := config.NewReaderLoader(buf, logger, format)
	assert.Nil(t, err, "failed to create reader (%v)", err)
	err = loader.Load(context.Background(), newCfg)
	assert.Nil(t, err, "failed to load config (%v)", err)
	// endregion

	// region Assert
	diff := cmp.Diff(
		cfg,
		newCfg,
		cmp.AllowUnexported(configuration.HTTPServerConfiguration{}),
		cmp.AllowUnexported(configuration.HTTPClientConfiguration{}),
		cmp.AllowUnexported(configuration.ClientConfig{}),
		cmp.AllowUnexported(configuration.KubernetesPodConfig{}),
		cmp.AllowUnexported(configuration.KubernetesConnectionConfig{}),
		cmp.AllowUnexported(configuration.DockerExecutionConfig{}),
		cmp.AllowUnexported(configuration.SyslogConfig{}),
		cmpopts.EquateEmpty(),
	)
	assert.Empty(t, diff)
	// endregion
}
