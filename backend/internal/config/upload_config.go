package config

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

// UploadConfig holds backend configuration for file uploads.
type UploadConfig struct {
	Storage   string         `json:"storage"`
	MaxSizeMB int            `json:"max_size_mb"`
	LocalDir  string         `json:"local_dir"`
	S3        UploadS3Config `json:"s3"`
}

// UploadS3Config holds S3-compatible storage configuration.
type UploadS3Config struct {
	Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	UseSSL    bool   `json:"use_ssl"`
	Prefix    string `json:"prefix"`
}

// LoadUploadConfig reads the upload configuration using the provided config instance.
// If no config instance is provided, it uses the global GoFrame configuration.
func LoadUploadConfig(ctx context.Context, cfgs ...*gcfg.Config) (UploadConfig, error) {
	cfg := g.Cfg()
	if len(cfgs) > 0 && cfgs[0] != nil {
		cfg = cfgs[0]
	}

	var uploadCfg UploadConfig
	value, err := cfg.Get(ctx, "upload")
	if err != nil {
		return uploadCfg, err
	}
	if value == nil {
		return uploadCfg, gerror.NewCode(gcode.CodeNotFound, "upload config not found")
	}
	if err := value.Scan(&uploadCfg); err != nil {
		return uploadCfg, err
	}
	return uploadCfg, nil
}
