package upload

import (
	"context"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

type s3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

var s3ClientFactory = newS3Client

func uploadToS3(ctx context.Context, cfg s3Config, file *ghttp.UploadFile, key string) error {
	if err := validateS3Config(cfg); err != nil {
		return err
	}

	client, err := s3ClientFactory(ctx, cfg)
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

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(cfg.Bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return gerror.NewCode(gcode.CodeOperationFailed, "s3 upload failed")
	}
	return nil
}

func newS3Client(ctx context.Context, cfg s3Config) (s3Client, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		awsConfig.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
				if strings.TrimSpace(cfg.Endpoint) == "" {
					return aws.Endpoint{}, &aws.EndpointNotFoundError{}
				}
				return aws.Endpoint{
					URL:           normalizeEndpoint(cfg.Endpoint),
					SigningRegion: cfg.Region,
				}, nil
			}),
		),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	}), nil
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
	if err == nil && parsed.Scheme == "" {
		hostOrPath := parsed.Host
		if hostOrPath == "" {
			hostOrPath = parsed.Path
		}
		return "https://" + hostOrPath
	}
	if err == nil {
		if parsed.Host == "" && !strings.Contains(endpoint, "://") {
			return "https://" + endpoint
		}
		return parsed.String()
	}
	return endpoint
}
