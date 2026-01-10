package upload

import (
	"context"
	"path"
	"strings"
	"time"

	"backend/internal/config"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
)

var (
	errCodeUploadTooLarge = gcode.New(413, "Payload Too Large", nil)
	disallowedExtensions  = map[string]struct{}{".exe": {}, ".bat": {}, ".cmd": {}, ".com": {}, ".scr": {}, ".msi": {}, ".ps1": {}, ".sh": {}}
)

type s3Config = config.UploadS3Config

var s3Upload = uploadToS3

func handleUpload(ctx context.Context, file *ghttp.UploadFile) (string, error) {
	if file == nil {
		return "", gerror.NewCode(gcode.CodeMissingParameter, "upload file is required")
	}

	cfg, err := config.LoadUploadConfig(ctx)
	if err != nil {
		return "", err
	}

	switch strings.ToLower(strings.TrimSpace(cfg.Storage)) {
	case "local", "":
		if err := validateUpload(file, cfg.MaxSizeMB); err != nil {
			return "", err
		}
		return saveLocalUpload(file, cfg.LocalDir)
	case "s3":
		if err := validateUpload(file, cfg.MaxSizeMB); err != nil {
			return "", err
		}
		filename := gfile.Basename(file.Filename)
		key := buildObjectKey(cfg.S3.Prefix, filename)
		if err := s3Upload(ctx, cfg.S3, file, key); err != nil {
			return "", err
		}
		return key, nil
	default:
		return "", gerror.NewCode(gcode.CodeNotSupported, "upload storage not supported")
	}
}

func validateUpload(file *ghttp.UploadFile, maxSizeMB int) error {
	if maxSizeMB > 0 {
		limit := int64(maxSizeMB) * 1024 * 1024
		if file.Size > limit {
			return gerror.NewCode(errCodeUploadTooLarge, "payload too large")
		}
	}

	ext := strings.ToLower(gfile.Ext(file.Filename))
	if ext != "" {
		if _, blocked := disallowedExtensions[ext]; blocked {
			return gerror.NewCode(gcode.CodeSecurityReason, "disallowed file type")
		}
	}
	return nil
}

func saveLocalUpload(file *ghttp.UploadFile, baseDir string) (string, error) {
	if strings.TrimSpace(baseDir) == "" {
		return "", gerror.NewCode(gcode.CodeInvalidConfiguration, "upload local_dir is required")
	}

	datePath := uploadDatePath()
	targetDir := gfile.Join(baseDir, datePath)
	savedName, err := file.Save(targetDir)
	if err != nil {
		return "", err
	}
	return "/" + path.Join("uploads", datePath, savedName), nil
}

func buildObjectKey(prefix, filename string) string {
	cleanPrefix := strings.Trim(prefix, "/")
	if cleanPrefix == "" {
		cleanPrefix = "uploads"
	}
	datePath := uploadDatePath()
	return path.Join(cleanPrefix, datePath, filename)
}

func uploadDatePath() string {
	return time.Now().Format("2006/01")
}
