package upload

import (
	"context"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func uploadToS3(ctx context.Context, cfg s3Config, file *ghttp.UploadFile, key string) error {
	if err := validateS3Config(cfg); err != nil {
		return err
	}

	endpoint := normalizeEndpoint(cfg.Endpoint)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Region: cfg.Region,
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return gerror.NewCode(gcode.CodeOperationFailed, "failed to initialize s3 client")
	}

	reader, err := file.Open()
	if err != nil {
		return gerror.Wrap(err, "failed to open upload file")
	}
	defer reader.Close()

	contentType := ""
	if file.Header != nil {
		contentType = file.Header.Get("Content-Type")
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = client.PutObject(ctx, cfg.Bucket, key, reader, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return gerror.NewCode(gcode.CodeOperationFailed, "s3 upload failed")
	}
	return nil
}

func validateS3Config(cfg s3Config) error {
	if strings.TrimSpace(cfg.Endpoint) == "" ||
		strings.TrimSpace(cfg.Bucket) == "" ||
		strings.TrimSpace(cfg.AccessKey) == "" ||
		strings.TrimSpace(cfg.SecretKey) == "" {
		return gerror.NewCode(gcode.CodeMissingConfiguration, "s3 configuration is incomplete")
	}
	return nil
}

func normalizeEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return endpoint
	}
	parsed, err := url.Parse(endpoint)
	if err == nil && parsed.Host != "" {
		return parsed.Host
	}
	return endpoint
}
