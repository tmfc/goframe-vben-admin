package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

type fakeS3Client struct {
	input *s3.PutObjectInput
	err   error
}

func (f *fakeS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	f.input = params
	return &s3.PutObjectOutput{}, f.err
}

func TestValidateS3ConfigMissing(t *testing.T) {
	err := validateS3Config(s3Config{})
	if err == nil {
		t.Fatal("expected error for missing s3 config")
	}
	if gerror.Code(err) != gcode.CodeMissingConfiguration {
		t.Fatalf("expected missing configuration code, got %v", gerror.Code(err))
	}
}

func TestNormalizeEndpoint(t *testing.T) {
	cases := map[string]string{
		"https://s3.example.com": "https://s3.example.com",
		"http://localhost:9000":  "http://localhost:9000",
		"localhost:9000":         "https://localhost:9000",
		"s3.example.com":         "https://s3.example.com",
	}
	for input, expected := range cases {
		if got := normalizeEndpoint(input); got != expected {
			t.Fatalf("normalizeEndpoint(%q) = %q, want %q", input, got, expected)
		}
	}
}

func TestUploadToS3Success(t *testing.T) {
	cfg := s3Config{
		Endpoint:  "https://s3.example.com",
		Region:    "us-east-1",
		Bucket:    "uploads",
		AccessKey: "access",
		SecretKey: "secret",
	}
	file := buildUploadFile(t, "remote.txt", []byte("hello"), "text/plain")

	previousFactory := s3ClientFactory
	fakeClient := &fakeS3Client{}
	s3ClientFactory = func(ctx context.Context, cfg s3Config) (s3Client, error) {
		return fakeClient, nil
	}
	t.Cleanup(func() { s3ClientFactory = previousFactory })

	key := "uploads/2026/01/remote.txt"
	if err := uploadToS3(context.Background(), cfg, file, key); err != nil {
		t.Fatalf("uploadToS3 failed: %v", err)
	}
	if fakeClient.input == nil || fakeClient.input.Bucket == nil || fakeClient.input.Key == nil {
		t.Fatal("expected PutObject input to be populated")
	}
	if *fakeClient.input.Bucket != cfg.Bucket {
		t.Fatalf("unexpected bucket: %s", *fakeClient.input.Bucket)
	}
	if *fakeClient.input.Key != key {
		t.Fatalf("unexpected key: %s", *fakeClient.input.Key)
	}
	if fakeClient.input.ContentType == nil || *fakeClient.input.ContentType != "text/plain" {
		t.Fatalf("unexpected content type: %v", fakeClient.input.ContentType)
	}
}

func TestUploadToS3ClientInitError(t *testing.T) {
	cfg := s3Config{
		Endpoint:  "https://s3.example.com",
		Region:    "us-east-1",
		Bucket:    "uploads",
		AccessKey: "access",
		SecretKey: "secret",
	}
	file := buildUploadFile(t, "remote.txt", []byte("hello"), "")
	file.FileHeader.Header = nil

	previousFactory := s3ClientFactory
	s3ClientFactory = func(ctx context.Context, cfg s3Config) (s3Client, error) {
		return nil, gerror.New("init failed")
	}
	t.Cleanup(func() { s3ClientFactory = previousFactory })

	if err := uploadToS3(context.Background(), cfg, file, "uploads/2026/01/remote.txt"); err == nil {
		t.Fatal("expected error from uploadToS3")
	} else if gerror.Code(err) != gcode.CodeOperationFailed {
		t.Fatalf("expected operation failed code, got %v", gerror.Code(err))
	}
}

func TestNewS3ClientSuccess(t *testing.T) {
	cfg := s3Config{
		Endpoint:  "https://s3.example.com",
		Region:    "us-east-1",
		Bucket:    "uploads",
		AccessKey: "access",
		SecretKey: "secret",
	}
	client, err := newS3Client(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newS3Client failed: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil s3 client")
	}
}

func buildUploadFile(t *testing.T, filename string, content []byte, contentType string) *ghttp.UploadFile {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	var part io.Writer
	var err error
	if contentType != "" {
		header := textproto.MIMEHeader{}
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
		header.Set("Content-Type", contentType)
		part, err = writer.CreatePart(header)
	} else {
		part, err = writer.CreateFormFile("file", filename)
	}
	if err != nil {
		t.Fatalf("failed to create multipart part: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("failed to write multipart content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(int64(len(content) + 1024))
	if err != nil {
		t.Fatalf("failed to read multipart form: %v", err)
	}
	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("no file found in multipart form")
	}
	return &ghttp.UploadFile{FileHeader: files[0]}
}
