package config

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestLoadUploadConfigFromFile(t *testing.T) {
	t.Parallel()

	cfgAdapter, err := gcfg.NewAdapterFile("testdata/upload-config.toml")
	if err != nil {
		t.Fatalf("failed to create config adapter: %v", err)
	}
	cfg := gcfg.NewWithAdapter(cfgAdapter)

	gtest.C(t, func(t *gtest.T) {
		uploadCfg, err := LoadUploadConfig(context.Background(), cfg)
		t.AssertNil(err)
		t.Assert(uploadCfg.Storage, "local")
		t.Assert(uploadCfg.MaxSizeMB, 25)
		t.Assert(uploadCfg.LocalDir, "resource/uploads")
		t.Assert(uploadCfg.S3.Endpoint, "https://s3.example.com")
		t.Assert(uploadCfg.S3.Bucket, "uploads")
		t.Assert(uploadCfg.S3.AccessKey, "access")
		t.Assert(uploadCfg.S3.SecretKey, "secret")
		t.Assert(uploadCfg.S3.UseSSL, true)
		t.Assert(uploadCfg.S3.Prefix, "uploads")
	})
}

func TestLoadUploadConfigMissing(t *testing.T) {
	t.Parallel()

	cfgAdapter, err := gcfg.NewAdapterFile("testdata/empty-config.toml")
	if err != nil {
		t.Fatalf("failed to create config adapter: %v", err)
	}
	cfg := gcfg.NewWithAdapter(cfgAdapter)

	gtest.C(t, func(t *gtest.T) {
		_, err := LoadUploadConfig(context.Background(), cfg)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)
	})
}
