package upload

import (
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"mime/multipart"
)

func TestBuildObjectKeyDefaultPrefix(t *testing.T) {
	key := buildObjectKey("", "file.txt")
	if key == "" || key[:7] != "uploads" {
		t.Fatalf("unexpected object key: %s", key)
	}
}

func TestValidateUploadTooLarge(t *testing.T) {
	file := &ghttp.UploadFile{
		FileHeader: &multipart.FileHeader{
			Filename: "file.txt",
			Size:     2 * 1024 * 1024,
		},
	}
	if err := validateUpload(file, 1); err == nil {
		t.Fatal("expected size validation error")
	} else if gerror.Code(err).Code() != 413 {
		t.Fatalf("unexpected error code: %v", gerror.Code(err))
	}
}

func TestValidateUploadDisallowedExtension(t *testing.T) {
	file := &ghttp.UploadFile{
		FileHeader: &multipart.FileHeader{
			Filename: "malware.exe",
			Size:     10,
		},
	}
	if err := validateUpload(file, 1); err == nil {
		t.Fatal("expected disallowed extension error")
	} else if gerror.Code(err) != gcode.CodeSecurityReason {
		t.Fatalf("unexpected error code: %v", gerror.Code(err))
	}
}
