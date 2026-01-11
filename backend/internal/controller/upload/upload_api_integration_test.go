package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

type apiEnvelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestUploadLocalSuccess(t *testing.T) {
	uploadDir := gfile.Temp(guid.S())
	if err := gfile.Mkdir(uploadDir); err != nil {
		t.Fatalf("failed to create upload dir: %v", err)
	}
	t.Cleanup(func() { gfile.Remove(uploadDir) })

	restoreConfig := setUploadConfig(t, "local", uploadDir, 1)
	t.Cleanup(restoreConfig)

	s := startUploadAPIServer(t)
	client := g.Client()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

	filePath := filepath.Join(uploadDir, "sample.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to write upload fixture: %v", err)
	}

	gtest.C(t, func(t *gtest.T) {
		content := client.PostContent(context.Background(), "/api/v1/upload", g.Map{
			"file": "@file:" + filePath,
		})
		env := decodeEnvelope(t, content)
		t.Assert(env.Code, gcode.CodeOK.Code())

		var res struct {
			Path string `json:"path"`
		}
		t.AssertNil(json.Unmarshal(env.Data, &res))
		t.Assert(strings.HasPrefix(res.Path, "/uploads/"), true)
		t.Assert(strings.HasSuffix(res.Path, "sample.txt"), true)
	})
}

func TestUploadLocalTooLarge(t *testing.T) {
	t.Skip("TODO: Fix file size validation in test environment")
}

func TestUploadLocalDisallowedType(t *testing.T) {
	uploadDir := gfile.Temp(guid.S())
	if err := gfile.Mkdir(uploadDir); err != nil {
		t.Fatalf("failed to create upload dir: %v", err)
	}
	t.Cleanup(func() { gfile.Remove(uploadDir) })

	restoreConfig := setUploadConfig(t, "local", uploadDir, 1)
	t.Cleanup(restoreConfig)

	s := startUploadAPIServer(t)
	client := g.Client()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

	filePath := filepath.Join(uploadDir, "malware.exe")
	if err := os.WriteFile(filePath, []byte("not really"), 0o644); err != nil {
		t.Fatalf("failed to write disallowed fixture: %v", err)
	}

	gtest.C(t, func(t *gtest.T) {
		content := client.PostContent(context.Background(), "/api/v1/upload", g.Map{
			"file": "@file:" + filePath,
		})
		env := decodeEnvelope(t, content)
		t.Assert(env.Code, gcode.CodeSecurityReason.Code())
	})
}

func TestUploadS3Success(t *testing.T) {
	t.Skip("TODO: Fix S3 upload test - needs complete S3 configuration or better mocking")
}

func TestUploadS3Failure(t *testing.T) {
	t.Skip("TODO: Fix S3 upload test - needs complete S3 configuration or better mocking")
}

func startUploadAPIServer(t *testing.T) *ghttp.Server {
	t.Helper()
	s := g.Server(guid.S())
	s.SetAddr(ghttp.FreePortAddress)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(NewV1())
	})
	s.SetDumpRouterMap(false)
	s.Start()
	time.Sleep(100 * time.Millisecond)
	t.Cleanup(func() { s.Shutdown() })
	return s
}

func setUploadConfig(t *testing.T, storage, localDir string, maxSizeMB int) func() {
	t.Helper()

	adapter, err := gcfg.NewAdapterFile()
	if err != nil {
		t.Fatalf("failed to create config adapter: %v", err)
	}
	adapter.SetContent(fmt.Sprintf(`
[upload]
storage = "%s"
max_size_mb = %d
local_dir = "%s"

[upload.s3]
endpoint = ""
region = ""
bucket = ""
access_key = ""
secret_key = ""
use_ssl = true
prefix = "uploads"
`, storage, maxSizeMB, localDir))

	cfg := g.Cfg()
	previous := cfg.GetAdapter()
	cfg.SetAdapter(adapter)
	return func() {
		cfg.SetAdapter(previous)
	}
}

func decodeEnvelope(t *gtest.T, content string) apiEnvelope {
	t.Helper()
	var env apiEnvelope
	t.AssertNil(json.Unmarshal([]byte(content), &env))
	return env
}
