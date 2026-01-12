package tests

import (
	"context"
	"testing"

	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func tableExists(ctx context.Context, db gdb.DB, tableName string) (bool, error) {
	v, err := db.GetValue(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = ?)", tableName)
	if err != nil {
		return false, err
	}
	return v.Bool(), nil
}

func Test_SysMessage_Tables_Exist(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := context.Background()
	db := g.DB()

	// Check sys_message
	exists, err := tableExists(ctx, db, "sys_message")
	if err != nil {
		t.Fatalf("Failed to check table existence: %v", err)
	}
	if !exists {
		t.Error("Table sys_message does not exist")
	}

	// Check sys_message_read
	exists, err = tableExists(ctx, db, "sys_message_read")
	if err != nil {
		t.Fatalf("Failed to check table existence: %v", err)
	}
	if !exists {
		t.Error("Table sys_message_read does not exist")
	}
}
